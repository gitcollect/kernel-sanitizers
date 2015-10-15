// go run qemuer.go -ninst=4 -maxcpu=4 --timeout=30m \
//	-diskimage=/disk/image/with/linux \
//	-id_rsa=/id_rsa/file/for/ssh \
//	-bzimage=/kernel/src/arch/x86/boot/bzImage \
//	-binary=/src/trinity/trinity \
//	-args="--dangerous -q -m -C 32"

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	diskimage = flag.String("diskimage", "", "path to linux disk image (required)")
	bzimage   = flag.String("bzimage", "", "path to kernel image, e.g. arch/x86/boot/bzImage (required)")
	qemu      = flag.String("qemu", "qemu-system-x86_64", "qemu binary")
	port      = flag.Int("port", 23505, "use ports [port, port+ninst) for ssh")
	id_rsa    = flag.String("id_rsa", "", "path to id_rsa file to ssh into the image (required)")
	ninst     = flag.Int("ninst", 1, "number of instances to use (1-100)")
	binary    = flag.String("binary", "", "path to binary to execute (required)")
	args      = flag.String("args", "", "arguments for the binary")
	maxcpu    = flag.Int("maxcpu", 1, "maximum number of cpus (1-128)")
	memmb     = flag.Int("memmb", 2048, "per-instance memory, in mb (128-8192)")
	timeout   = flag.Duration("timeout", 60*time.Minute, "timeout for binary execution")
	v         = flag.Bool("v", false, "dump instance output to stdout")

	images = make(chan string)
)

func main() {
	flag.Parse()
	if *binary == "" || *diskimage == "" || *bzimage == "" || *id_rsa == "" ||
		*ninst <= 0 || *ninst > 100 || *maxcpu <= 0 || *maxcpu > 128 ||
		*memmb < 128 || *memmb > 819200 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *v && *ninst > 1 {
		// Intermixed output from several intances is unuseful.
		*ninst = 1
	}
	log.Printf("fuzzing on %v instance(s)", *ninst)
	go CopyImageLoop()
	for i := 0; i < *ninst; i++ {
		go FuzzLoop(i)
		time.Sleep(time.Second)
	}
	select {}
}

func CopyImageLoop() {
	for i := 0; ; i++ {
		t0 := time.Now()
		oldf, err := os.Open(*diskimage)
		if err != nil {
			log.Fatalf("failed to open original image: %v", err)
		}
		defer oldf.Close()
		image := fmt.Sprintf("asanimage%v", i)
		newf, err := os.Create(image)
		if err != nil {
			log.Fatalf("failed to open new image: %v", err)
		}
		_, err = io.Copy(newf, oldf)
		if err != nil {
			log.Fatalf("failed to copy image: %v", err)
		}
		oldf.Close()
		newf.Close()
		log.Printf("copied image (%v)", time.Now().Sub(t0))
		images <- image
	}
}

func FuzzLoop(i int) {
	ncpu := *maxcpu
	for run := 0; ; run++ {
		logname := fmt.Sprintf("asan%v-%v-%v.log", i, run, time.Now().Unix())
		logf, err := os.Create(logname)
		if err != nil {
			log.Printf("failed to create log file: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}
		log := io.Writer(logf)
		if *v {
			log = io.MultiWriter(log, os.Stdout)
		}
		inst := &Instance{
			name:    fmt.Sprintf("asan%v", i),
			id:      i,
			runid:   run,
			logname: logname,
			log:     log,
			ncpu:    ncpu,
			cmds:    make(map[*Command]bool),
		}
		inst.Run()
		inst.Shutdown()
	}
}

type Instance struct {
	sync.Mutex
	name    string
	image   string
	id      int
	runid   int
	ncpu    int
	logname string
	log     io.Writer
	cmds    map[*Command]bool
	qemu    *Command
}

type Command struct {
	sync.Mutex
	cmd    *exec.Cmd
	done   chan struct{}
	failed bool
	out    []byte
	outpos int
	inw    io.Writer
}

var bootSem = make(chan bool, 1)

func (inst *Instance) Run() bool {
	// Obtain brand new image.
	bootSem <- true
	t0 := time.Now()
	inst.image = <-images

	// Start the instance.
	inst.qemu = inst.CreateCommand(*qemu,
		"-hda", inst.image,
		"-m", strconv.Itoa(*memmb),
		"-smp", strconv.Itoa(inst.ncpu),
		"-net", fmt.Sprintf("user,hostfwd=tcp::%v-:22", *port+inst.id),
		"-net", "nic",
		"-nographic",
		"-kernel", *bzimage,
		"-append", "console=ttyS0 root=/dev/sda debug earlyprintk=serial",
		"-enable-kvm")
	if !inst.qemu.WaitForOut("Starting OpenBSD Secure Shell server: sshd", 5*time.Minute) {
		<-bootSem
		inst.Logf("failed to start qemu")
		return false
	}
	<-bootSem
	inst.Logf("started vm (%v)", time.Now().Sub(t0))

	// Run the binary.
	t0 = time.Now()
	nfailfast := 0
	reported := true // not very useful (and broken because of command prompt wait)
	deadline := time.Now().Add(*timeout)
	for try := 0; ; try++ {
		start := time.Now()
		conoutlen := inst.qemu.OutLen()
		lastout := time.Now()
		// Copy the binary into the instance.
		inst.CreateSCPCommand(*binary, "/tmp/kasaner_binary").Wait(3 * time.Minute)
		cmd := inst.CreateSSHCommand("/tmp/kasaner_binary " + *args)
		for {
			time.Sleep(time.Second)
			// Check for AddressSanitizer reports in output.
			if !reported && inst.qemu.CheckOut("KASAN") {
				reported = true
				inst.Logf("found KASAN report in %v", inst.logname)
			}
			if time.Now().After(deadline) {
				inst.Logf("test %v interrupted", try)
				return true
			}
			if cmd.Exited() {
				t := time.Now().Sub(start)
				inst.Logf("test %v finished after %v", try, t)
				if t < 20*time.Second {
					// If the binary exits too quickly often,
					// assume the vm is broken.
					nfailfast++
					if nfailfast > 3 {
						inst.Logf("vm is broken (binary can't start)")
						return false
					}
				} else {
					nfailfast = 0
				}
				break
			}
			if conoutlen != inst.qemu.OutLen() {
				conoutlen = inst.qemu.OutLen()
				lastout = time.Now()
			} else if lastout.Add(5 * time.Minute).Before(time.Now()) {
				// If no output for more than 5 minutes,
				// assume the vm is broken.
				inst.Logf("vm hang after %v (no output)", time.Now().Sub(start))
				return false
			}
		}
	}
}

func (inst *Instance) Shutdown() {
	for try := 0; try < 10; try++ {
		inst.qemu.cmd.Process.Kill()
		time.Sleep(time.Second)
		inst.Lock()
		n := len(inst.cmds)
		inst.Unlock()
		if n == 0 {
			os.Remove(inst.image)
			return
		}
	}
	inst.Logf("hanged processes after kill")
	inst.Lock()
	for cmd := range inst.cmds {
		cmd.cmd.Process.Kill()
	}
	inst.Unlock()
	time.Sleep(3 * time.Second)
	os.Remove(inst.image)
}

func (inst *Instance) CreateCommand(args ...string) *Command {
	fmt.Fprintf(inst.log, "command %v\n", args)
	cmd := &Command{}
	cmd.done = make(chan struct{})
	outr, outw, err := os.Pipe()
	if err != nil {
		inst.Logf("failed to create pipe: %v", err)
		cmd.failed = true
		close(cmd.done)
		return cmd
	}
	inr, inw, err := os.Pipe()
	if err != nil {
		inst.Logf("failed to create pipe: %v", err)
		cmd.failed = true
		close(cmd.done)
		return cmd
	}
	cmd.cmd = exec.Command(args[0], args[1:]...)
	cmd.cmd.Stdout = io.MultiWriter(inst.log, outw)
	cmd.cmd.Stderr = io.MultiWriter(inst.log, outw)
	cmd.cmd.Stdin = inr
	cmd.inw = inw
	go func() {
		err := cmd.cmd.Start()
		if err == nil {
			inst.Lock()
			inst.cmds[cmd] = true
			inst.Unlock()
			err = cmd.cmd.Wait()
			inst.Lock()
			delete(inst.cmds, cmd)
			inst.Unlock()
		}
		fmt.Fprintf(inst.log, "command %v DONE %v\n", args, err)
		inw.Close()
		inr.Close()
		outw.Close()
		outr.Close()
		cmd.failed = err != nil
		close(cmd.done)
	}()
	go func() {
		var buf [4 << 10]byte
		for {
			n, err := outr.Read(buf[:])
			if err != nil {
				return
			}
			cmd.Lock()
			cmd.out = append(cmd.out, buf[:n]...)
			cmd.Unlock()
		}
	}()
	return cmd
}

func (inst *Instance) CreateSSHCommand(args ...string) *Command {
	args1 := []string{"ssh", "-i", *id_rsa, "-p", strconv.Itoa(*port + inst.id),
		"-o", "ConnectionAttempts=10", "-o", "ConnectTimeout=60",
		"-o", "BatchMode=yes", "-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no", "root@localhost"}
	return inst.CreateCommand(append(args1, args...)...)
}

func (inst *Instance) CreateSCPCommand(from, to string) *Command {
	return inst.CreateCommand("scp", "-i", *id_rsa, "-P", strconv.Itoa(*port+inst.id),
		"-o", "ConnectionAttempts=10", "-o", "ConnectTimeout=60",
		"-o", "BatchMode=yes", "-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		from, "root@localhost:"+to)
}

func (inst *Instance) Logf(str string, args ...interface{}) {
	fmt.Fprintf(inst.log, str+"\n", args...)
	if !*v {
		log.Printf("%v-%v: "+str, append([]interface{}{inst.name, inst.runid}, args...)...)
	}
}

func (cmd *Command) Wait(max time.Duration) bool {
	select {
	case <-cmd.done:
		return !cmd.failed
	case <-time.After(max):
		return false
	}
}

func (cmd *Command) Exited() bool {
	select {
	case <-cmd.done:
		return true
	default:
		return false
	}
}

func (cmd *Command) OutLen() int {
	cmd.Lock()
	defer cmd.Unlock()
	return len(cmd.out)
}

func (cmd *Command) CheckOut(what string) bool {
	cmd.Lock()
	defer cmd.Unlock()
	idx := bytes.Index(cmd.out[cmd.outpos:], []byte(what))
	if idx == -1 {
		newpos := len(cmd.out) - len(what)
		if newpos > cmd.outpos {
			cmd.outpos = newpos
		}
	} else {
		cmd.outpos += idx + len(what)
	}
	return idx != -1
}

func (cmd *Command) WaitForOut(what string, max time.Duration) bool {
	if cmd.CheckOut(what) {
		return true
	}
	ticker := time.NewTicker(time.Second)
	timeout := time.NewTimer(max)
	defer func() {
		ticker.Stop()
		timeout.Stop()
	}()
	for {
		select {
		case <-cmd.done:
			return cmd.CheckOut(what)
		case <-ticker.C:
			if cmd.CheckOut(what) {
				return true
			}
		case <-timeout.C:
			return false
		}
	}
}

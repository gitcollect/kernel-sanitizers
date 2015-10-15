#!/usr/bin/python

import json
import os

import matplotlib.pyplot as plt

def get_files(path):
  for root, dirs, files in os.walk(path):
    return files 

def get_lines(filename):
  rv = []
  with open(filename) as f:
    for line in f:
      rv.append(line.strip())
  return rv

def get_log_paths(dirname):
  files = get_files(dirname)
  rv = []
  for filename in files:
    path = os.path.join(dirname, filename)
    if path.startswith(dirname + '/' + 'asan0'):
      rv.append(path)
  return rv

def log_path_to_time(path):
  return path.split('-')[-1].split('.')[-2]

def get_log_paths_times(paths):
  rv = {}
  for path in paths:
    rv[path] = int(log_path_to_time(path))
  return rv

def normalize_log_paths_times(paths_times):
  m = min(paths_times.items(), key=lambda x: x[1])[1]
  for key in paths_times.keys():
    paths_times[key] -= m
  return paths_times

def log_line_to_time(line):
  rv = None
  try:
    rv = int(line.split(']')[0].lstrip('[ ').split('.')[0])
  except:
    return None
  return rv

def get_race_time(log_line, races):
  for race in races:
    if race in log_line:
      time = log_line_to_time(log_line)
      if time != None:
	return time
  return None

def get_races_times(log_path, races):
  rv = []
  with open(log_path) as f:
    for line in f:
      line = line.strip()
      time = get_race_time(line, races)
      if time != None:
        rv.append(time)
  return rv

def get_all_races_times(dirname, races):
  log_paths = get_log_paths(dirname)
  log_paths_times = normalize_log_paths_times(get_log_paths_times(log_paths))

  all_times = []

  for path in log_paths:
    times = get_races_times(path, races)
    times = [x + log_paths_times[path] for x in times]
    all_times += times

  all_times.sort()
  return all_times

def generate_step_plot_ranges(steps):
  x_max = max(steps) + 1
  x = range(x_max)

  y = []
  v = 0
  for i in xrange(x_max):
    if i in steps:
      v += 1
    y.append(v)

  return x, y

def get_all_races_plot_ranges(dirname, races):
  times = get_all_races_times(dirname, races)
  times = times[:-3]
  x, y = generate_step_plot_ranges(times)
  return x, y

old_dirname = './old'
new_dirname = './new'
old_races = get_lines('./old-races')
new_races = get_lines('./new-races')

x, y = get_all_races_plot_ranges(old_dirname, old_races)
plt.step(x, y, color='g')

x, y = get_all_races_plot_ranges(new_dirname, old_races)
plt.step(x, y, color='r')

plt.show()

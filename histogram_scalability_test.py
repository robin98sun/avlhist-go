#!/usr/bin/env python3

import math
import os
import sys
import json
import time
import subprocess
from random import random
import requests
import multiprocessing
import re

if __name__ == "__main__":
    
    import argparse
    parser = argparse.ArgumentParser(description='generate jobs according to the configuration file')

    parser.add_argument(
        '--logs', type=str, required=True, nargs="+",
        help='the log file path'
    )

    args = parser.parse_args()

    program_start_time = time.time()

    def new_perf_obj(hist_count, window_size, percentile, latency, count):
        return {
            "hist_count": hist_count,
            "window_size": window_size,
            "percentile": percentile,
            "latency": latency,
            "count": count,
        }

    def save_to_row(rows, row_key, col_key, perf_obj):
        if row_key not in rows:
            rows[row_key] = {}
        if col_key not in rows[row_key]:
            rows[row_key][col_key] = perf_obj
        else:
            for prop in ["latency"]:
                rows[row_key][col_key][prop] = (rows[row_key][col_key][prop]*rows[row_key][col_key]["count"] + perf_obj[prop])/(perf_obj["count"]+rows[row_key][col_key]["count"])
            rows[row_key][col_key]["count"] = perf_obj["count"]+rows[row_key][col_key]["count"]


    rows = {}

    for file in args.logs:
        with open(file, 'r') as f:
            lines = f.readlines()
            f.close()


        for line in lines:
            if "BenchmarkTestScheduler_MultiplyHistograms/multiply_" in line:
                line = re.sub(' +', ' ', line)
                line = re.sub('\n', '', line)
                hist_count = int(line.split("_")[2])
                window_size = int(line.split("_")[7])
                percentile = float(line.split("_")[10].split(" ")[0])

                latency = round(float(line.split(" ")[1]) * float(line.split(" ")[2])/1000/1000, 1)
                
                perf_obj = new_perf_obj(hist_count, window_size, percentile, latency, 1)

                row_key = "{}-{}".format(hist_count, percentile)
                col_key = "{}".format(window_size)
                save_to_row(rows, row_key, col_key, perf_obj)
                
    lines = [] 
    for row_key in rows:
        line = ""
        for col_key in rows[row_key]:
            perf_obj = rows[row_key][col_key]
            if line == "":
                line = "{}".format(perf_obj["hist_count"])
            line = "{}, {}".format(line, perf_obj["latency"])
        lines.append(line)

    sorted(lines, key = lambda x: int(x.split(",")[0]))
    for line in lines:
        print(line)



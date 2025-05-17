import os
import sys
import re
from subprocess import PIPE, run

def Out(command):
    result = run(command, stdout=PIPE, stderr=PIPE, universal_newlines=True, shell=True, encoding='utf-8')
    return result.stdout

def LaunchTest(command_line, filename, success):
    print(f"Launching: {command_line}")
    output = Out(command_line).encode('utf-8', errors='ignore').decode(errors='ignore')

    if re.search(success, output):
        print(f"Erorr: found proof for file ${name}. Goéland is unsound.")
        exit(1)

if len(sys.argv) < 3:
    print(f"python3 {sys.argv[0]} problem_folder timeout [additional arguments]")
else:
    folder = sys.argv[1]
    folder_split = folder.split("/")
    folder += "/"

    entries = os.listdir(folder)
    timeout = sys.argv[2]
    additionalArgs = " ".join(sys.argv[3:])
    for index, file in enumerate(entries):
        print(f"Problem {index+1}/{len(entries)} : {folder+file}")
        LaunchTest("timeout "+timeout+" ../src/_build/goeland " + additionalArgs + " " +folder+file, file, "% RES : VALID")

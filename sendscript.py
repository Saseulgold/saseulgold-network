import csv
import subprocess
import time
import sys

csv_file_path = sys.argv[1]

command_template = ["./sg", "wallet", "send", "-t", "{address}", "-a", "{amount}"]

# Reading and iterating CSV files
with open(csv_file_path, newline='', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile)  # Read as dictionary using first row as header
    
    for row in reader:
        print(row)
        command = [arg.format(**row) if "{" in arg else arg for arg in command_template]
        print(command)

        try:
            stdout = subprocess.run(command, check=True, capture_output=True, text=True)
            print(stdout)
        except subprocess.CalledProcessError as e:
            print(f"Error Occurred: {e}")
        


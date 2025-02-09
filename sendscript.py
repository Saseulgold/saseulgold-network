import csv
import subprocess
import time
import sys

csv_file_path = sys.argv[1]

command_template = ["./sg", "wallet", "send", "-t", "{address}", "-a", "{amount}"]

# CSV 파일 읽기 및 반복 실행
with open(csv_file_path, newline='', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile)  # 첫 번째 행을 헤더로 사용하여 딕셔너리로 읽음
    
    for row in reader:
        print(row)
        command = [arg.format(**row) if "{" in arg else arg for arg in command_template]
        print(command)

        try:
            stdout = subprocess.run(command, check=True, capture_output=True, text=True)
            print(stdout)
        except subprocess.CalledProcessError as e:
            print(f"오류 발생: {e}")
        


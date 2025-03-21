# 📌 SaseulGold mining guide

### Minimum (recommended) node operation specifications
**8gb (12gb)** or more memory  
cpu at least **4 cores (8 cores)**  
**96GB (256GB)** or higher disk

### Recommended OS Environment
Cuda Version : **12.4**  
Driver Version : **550.127.05**  
Aws Ec2 AMI: Deep Learning Base OSS Nvidia Driver GPU AMI (**Ubuntu 22.04**)

## 🔹 Prerequisites
SaseulGold mining requires **a GPU environment**. Ensure the following commands are available before proceeding.

✅ `nvcc` (CUDA Compiler)  
✅ `nvidia-smi` (NVIDIA System Management Interface)

**Based on AWS EC2**, **Deep Learning Base OSS Nvidia Driver GPU AMI (Ubuntu 22.04)** images allow you to run immediately without additional settings.

---

## 🛠️ Installing and Running the SaseulGold Client

### 1️⃣ Downloading and Installing the SaseulGold Client
```bash
wget https://github.com/Saseulgold/saseulgold-network/raw/refs/heads/main/sg-main.zip
unzip sg-main.zip -d sg_network  # Uncompressed
cd sg_network  # Move folder
```
### 2️⃣ Running a pre-installation script
```bash
sh init.sh
```
### 3️⃣ Create a Wallet
```bash
./sg wallet create
✅ example:
Private Key: 8d7a0bb37a9044aba0dab18968b8ad6f071790c2429de209855bc041d904833d
Public Key: 0b2da6013dda3bfd5dc1a24e4927f151b33b027dd8b7402516d4c5020c04fd18
Wallet address: eb4e0202345542b5e3405debd9385043f4a852411a1a
```
### 4️⃣ Initial Wallet Settings
```bash
./sg wallet set -k ${private_key}
```
### 5️⃣ Test (Check balance)
```bash
./sg wallet balance
```
### 6️⃣ Start mining
```bash
./sg mining start
```
### 7️⃣ End of mining
```bash
./sg mining stop
```


---

## 🔄 Setting Up an Auto-Restart Script for Mining
To ensure the mining process restarts automatically if it stops, you can set up the following script.
### 1️⃣ Create the Restart Script
```bash
sudo vim /home/SSG/auto_start_ssg.sh
```
Add the following content to the file:
```bash
#!/bin/bash

# Assumes the script is run from the sg_network folder
cd <sg_network folder>

gpu_status="driver not installed"
gpu_model="unknown"

if command -v nvidia-smi > /dev/null 2>&1; then
    output=$(nvidia-smi -L 2>&1) 

    case "$output" in
        *"Unknown Error"* | *"No devices"* | *"NVIDIA_SMI has Failed"*)
            gpu_status="not found"
            ;;
        *"version mismatch"*)
            gpu_status="version mismatch"
            ;;
        *)
            gpu_status="detected"

            if [ "$(nvidia-smi | grep -c 'cmine')" -gt 1 ]; then
                echo "Multiple SSG Detected - Reboot."
                sleep 120
                reboot
            fi
            ;;
    esac
fi

miner_status=$(pgrep cmine > /dev/null && echo "running" || echo "stopped")

if [[ "$gpu_status" == "detected" && "$miner_status" == "stopped" ]]; then
    ./sg mining start
    sleep 1
    miner_status="starting"
fi

timestamp=$(date +"%y-%m-%d %H:%M:%S")

ssg_address=$(grep "Address:" wallet.info | awk '{print $2}')

echo "Miner Status: $miner_status"
echo "GPU Status: $gpu_status"
echo "Checked Time: $timestamp"
echo "SSG Address: $ssg_address"
echo "============================================================="
echo "============================================================="
```

### 2️⃣ Schedule with Crontab (Run Every 5 Minutes)
```bash
sudo crontab -e
```
Add the following line, replacing /path/to/sg_network with the actual path to your sg_network folder:
```bash
*/5 * * * * bash -c /path/to/sg_network/auto_start_ssg.sh >> /path/to/sg_network/auto_start_ssg.log 2>&1
```
>Tip: To find the full path, run pwd inside the sg_network folder and use that in the crontab command.

### 3️⃣ Grant Execution Permissions
```bash
sudo chmod 755 auto_start_ssg.sh
```
or
```bash
sudo chmod -x auto_start_ssg.sh
```
>Recommand : Use sudo chmod 755 to make the script executable. Avoid sudo chmod -x as it removes execution permissions.

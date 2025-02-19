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

## 🛠️ Setup & Installation for BLE Gateway
This section provides a step-by-step guide to set up the BLE Gateway on an Ubuntu machine.

### Prerequisites
- Hardware: Raspberry Pi 3B+ or any device with BLE (Bluetooth Low Energy) support
- Operating System: Ubuntu 20.04 or later
- Software: Go (Golang) programming language, Bluetooth tools, and dependencies
---
### Steps

#### 1. **Clone the Repository**
First, clone this repository to your local machine.
```
git clone https://github.com/yourusername/ble-gateway.git
cd ble-gateway
```

#### 2. Update System Packages
Ensure that your system is up-to-date.
```
sudo apt update
sudo apt upgrade -y
```

#### 3. Install Bluetooth Packages
Install the necessary Bluetooth packages for Ubuntu.
```
sudo apt install -y bluetooth bluez
```
Enable and start the Bluetooth service.
```
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```
Verify that the Bluetooth service is running.
```
sudo systemctl status bluetooth
```

#### 4. Install Go
Install the Go programming language to run the BLE Gateway code.
```
sudo apt install -y golang
```
Verify the Go installation by checking its version.
```
go version
```

#### 5. Set Up Go Workspace
Navigate to the ```ble-gateway``` directory and initialize the Go module. This will create a ```go.mod``` file in the directory.
```
go mod init ble-gateway
```

#### 6. Install Go Bluetooth Library
To use Bluetooth Low Energy (BLE) in Go, install a compatible Bluetooth library. Here’s an example using ```tinygo-org/bluetooth```:
```
go get tinygo.org/x/bluetooth
```
Note: ```tinygo-org/bluetooth``` is primarily designed for TinyGo. Ensure this library meets your requirements, or consider alternatives if compatibility issues arise.

#### 7. Create the main.go File
Create the ```main.go``` file to define the main logic for the BLE Gateway.
```
nano main.go
```
Inside ```main.go```, you can implement the BLE scanning and gRPC communication logic, or refer to the example code provided in the repository.

#### 8. Run the BLE Gateway
Once all dependencies are installed, you can run the BLE Gateway using the following command:
```
go run main.go
```

---
Following these steps will set up the BLE Gateway on an Ubuntu machine, ready to detect BLE signals and communicate with the server.
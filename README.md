# BALogin: The Beacon-based Automatic Login System for Enhancing Platform Accessibility for the Elderly

This repository contains the implementation of a **Beacon-based Automatic Login System** designed to enhance digital platform accessibility, particularly for elderly users. Utilizing Bluetooth Low Energy (BLE) technology, this system enables automatic login upon proximity to specific locations such as public institutions or care facilities, making login processes seamless and convenient for users with limited digital experience.

## Project Structure

The system is divided into main components:

1. **BLE Sensor** - A portable device carried by the user to act as a BLE beacon, periodically transmitting unique UUID signals to nearby gateways. It operates on low energy to optimize battery usage, making it ideal for prolonged usage.
2. **BLE Gateway** - A receiver stationed at designated locations to detect BLE signals from the sensor and send UUIDs to the server for automatic login verification.

---

### 1. BLE Sensor

- **Device:** ESP32-WROOM-32
- **Communication Protocol:** Beacon for communication with the gateway.
- **File Structure:**
  ```
  ble-device/
  └── src/
  │  └── main.cpp
  └── platformio.ini
  ```

### 2. BLE Gateway

- **Device:** Raspberry Pi 3B+ with Ubuntu Server 24.04.1 LTS
- **Database:** SQLite for local storage of registered UUIDs.
- **Communication Protocol:** gRPC for communication with the server.
- **File Structure:** 
  ```
  ble-gateway/
  ├── ble/                    
  │   └── ble.go             
  ├── db/                     
  │   └── db.go              
  ├── handler/                
  │   ├── create.go         
  │   └── status.go           
  ├── proto/
  │   ├── ble.proto
  │   ├── ble.pb.go
  │   └── ble_grpc.pb.go
  ├── main.go
  ├── go.mod
  └── go.sum
  ```
- **Database Structure:**
  ```sql
  CREATE TABLE devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_name TEXT NOT NULL,
    uuid TEXT NOT NULL UNIQUE,
    is_active INTEGER NOT NULL DEFAULT 0
  );
  ```

---

## System Workflow

1. **User Registration:** A unique UUID is assigned to each BLE sensor during user registration. This UUID is stored on both the sensor and the server database.
2. **Automatic Login:** When the user, carrying the BLE sensor, approaches a gateway, the gateway detects the BLE signal, retrieves the UUID, and sends it to the server via gRPC for login.
3. **Automatic Logout:** If the gateway loses the BLE signal for a set period, it triggers an automatic logout by informing the server.
<img width="793" alt="Screenshot 2024-11-10 at 4 13 31 PM" src="https://github.com/user-attachments/assets/2bb4de5d-2123-4db1-892b-7ed4205eaebb">

---

## System Architecture
<img width="730" alt="Screenshot 2024-11-10 at 4 13 45 PM" src="https://github.com/user-attachments/assets/23853f81-5f06-4859-9545-2497572c7792">

---

## Setup & Installation
1. **BLE Sensor**
   [Click Details](ble-device/README.md)
   ```
   Hardware: ESP32-WROOM-32
   IDE: PlatformIO
   ```
2. **BLE Gateway**
   [Click Details](ble-gateway/README.md)
   ```
   Hardware: Raspberry Pi 3B+
   OS: Ubuntu Server 24.04.1 LTS
   Database: SQLite (using `go-sqlite3` v1.14.24)
   Go Version: 1.22.2 (linux/arm64)
   gRPC Library: grpc-go v1.67.1
   Protocol Buffers Library: protobuf v1.35.1
   BLE Library: tinygo-org/bluetooth v0.10.0
   ```

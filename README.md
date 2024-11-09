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

---

## System Architecture

---

## Setup & Installation
1. **BLE Sensor**
   [Click Details]
   ```
   Hardware: ESP32-WROOM-32
   IDE: PlatformIO
   ```
2. **BLE Gateway**
   [Click Details]
   ```
   Hardware: Raspberry Pi 3B+
   OS: Ubuntu Server 24.04.1 LTS
   Database:
   Go Version:
   gRPC Library:
   BLE Library:
   ```

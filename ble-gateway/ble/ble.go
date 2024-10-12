package ble

import (
    "database/sql"
    "fmt"
    "log"
    "sync"
    "time"
    "tinygo.org/x/bluetooth"
    _ "github.com/mattn/go-sqlite3" // SQLite3 드라이버
    pb "ble-gateway/proto"
    "ble-gateway/handler"
)

// BLE 어댑터 초기화
var adapter = bluetooth.DefaultAdapter

var connectedDevices = make(map[string]string) // MAC 주소 -> UUID 매핑
var lastSeen = make(map[string]time.Time)      // MAC 주소 -> 마지막으로 감지된 시간
var mu sync.Mutex

const RSSIThreshold = -90
const timeoutDuration = 30 * time.Second
const scanInterval = 3 * time.Second

// SQLite 데이터베이스 열기
func openDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./ble.db")
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    return db, nil
}

// 특정 UUID의 is_active 값 확인
func isDeviceActive(db *sql.DB, uuid string) (bool, error) {
    var isActive int
    query := `SELECT is_active FROM devices WHERE uuid = ?`
    err := db.QueryRow(query, uuid).Scan(&isActive)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil // UUID가 데이터베이스에 없으면 연결하지 않음
        }
        return false, fmt.Errorf("failed to check device active status: %v", err)
    }
    return isActive == 1, nil // is_active가 1이면 true, 아니면 false
}

// BLE 스캔 재시작 및 최신 장치 상태 반영 함수
func RestartScan(client pb.DeviceServiceClient) {
    db, err := openDB()
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer db.Close()

    for {
        time.Sleep(scanInterval)

        fmt.Println("Restarting BLE scan to refresh device states...")
        checkTimeouts(client)

        err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
            if result.LocalName() == "" {
                return
            }

            macAddress := result.Address.String()
            lastSeen[macAddress] = time.Now()

            device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
            if err != nil {
                handleDisconnect(macAddress, client)
                return
            }
            defer device.Disconnect()

            services, err := device.DiscoverServices(nil)
            if err != nil {
                handleDisconnect(macAddress, client)
                return
            }

            for _, service := range services {
                uuid := service.UUID().String()
                if uuid != "00001801-0000-1000-8000-00805f9b34fb" {
                    if result.RSSI > RSSIThreshold {
                        isActive, err := isDeviceActive(db, uuid)
                        if err != nil {
                            log.Printf("Error checking device active status: %v", err)
                            continue
                        }

                        if isActive {
                            handleConnect(macAddress, uuid, client)
                            fmt.Printf("Device Address: %s, UUID: %s, RSSI: %d\n", macAddress, uuid, result.RSSI)
                        } else {
                            fmt.Printf("Device UUID %s is not active. Skipping connection.\n", uuid)
                            handleDisconnect(macAddress, client)
                        }
                    } else {
                        handleDisconnect(macAddress, client)
                    }
                }
            }
        })

        if err != nil {
            log.Printf("Error restarting BLE scan: %v", err)
        }
    }
}

// BLE 장치 연결/해제 함수
func handleDisconnect(macAddress string, client pb.DeviceServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    if uuid, exists := connectedDevices[macAddress]; exists {
        fmt.Printf("Device %s disconnected.\n", uuid)
        delete(connectedDevices, macAddress)
        delete(lastSeen, macAddress)
        handler.SendDeviceStatus(client, uuid, 0)
    }
}

func handleConnect(macAddress string, uuid string, client pb.DeviceServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    if _, exists := connectedDevices[macAddress]; !exists {
        fmt.Printf("Device %s connected.\n", uuid)
        connectedDevices[macAddress] = uuid
        lastSeen[macAddress] = time.Now()
        handler.SendDeviceStatus(client, uuid, 1)
    } else {
        lastSeen[macAddress] = time.Now()
    }
}

// 장치의 타임아웃을 체크하고, 감지되지 않으면 disconnect 처리
func checkTimeouts(client pb.DeviceServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    currentTime := time.Now()

    for macAddress, lastSeenTime := range lastSeen {
        if currentTime.Sub(lastSeenTime) > timeoutDuration {
            fmt.Printf("Device %s timed out (no signal for %v).\n", connectedDevices[macAddress], timeoutDuration)
            handleDisconnect(macAddress, client)
        }
    }
}

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
    "tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// 연결된 장치의 UUID를 저장하는 맵과 락
var connectedDevices = make(map[string]bool)
var mu sync.Mutex

// RSSI 임계값
const RSSIThreshold = -90

// 연결 및 해제 시 호출할 URL
const connectURL = "https://6e316739-cad9-44a2-8fae-647a4b9f2444.mock.pstmn.io/connected"
const disconnectURL = "https://6e316739-cad9-44a2-8fae-647a4b9f2444.mock.pstmn.io/disconnected"

// POST 요청으로 UUID 전송 함수
func sendPostRequest(url string, uuid string) {
    jsonData := map[string]string{"uuid": uuid}  // UUID 데이터를 JSON 형식으로 준비
    jsonValue, _ := json.Marshal(jsonData)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))  // POST 요청 전송
    if err != nil {
        log.Printf("Failed to send POST request: %v", err)
        return
    }
    defer resp.Body.Close()
    fmt.Printf("Sent POST request to %s with UUID: %s, Status: %s\n", url, uuid, resp.Status)
}

// BLE 장치 연결이 끊겼을 때 처리하는 함수
func handleDisconnect(deviceUUID string) {
    mu.Lock()
    defer mu.Unlock()
    if connectedDevices[deviceUUID] {
        fmt.Printf("Device %s disconnected.\n", deviceUUID)
        connectedDevices[deviceUUID] = false
        sendPostRequest(disconnectURL, deviceUUID)  // UUID 전송
    }
}

// BLE 장치가 연결되었을 때 처리하는 함수
func handleConnect(deviceUUID string) {
    mu.Lock()
    defer mu.Unlock()
    if !connectedDevices[deviceUUID] {
        fmt.Printf("Device %s connected.\n", deviceUUID)
        connectedDevices[deviceUUID] = true
        sendPostRequest(connectURL, deviceUUID)  // UUID 전송
    }
}

// BLE 스캔 재시작 및 최신 장치 상태 반영 함수
func restartScan() {
    for {
        time.Sleep(10 * time.Second)

        fmt.Println("Restarting BLE scan to refresh device states...")

        err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
            if result.LocalName() == "" {
                return
            }

            deviceUUID := result.Address.String()

            // 연결 끊김 처리 (RSSI 확인)
            if result.RSSI <= RSSIThreshold {
                handleDisconnect(deviceUUID)
            } else {
                device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
                if err != nil {
                    // 연결 실패 시 Disconnect 처리
                    handleDisconnect(deviceUUID)
                } else {
                    defer device.Disconnect()

                    // 서비스 탐색을 통해 UUID를 확인
                    services, err := device.DiscoverServices(nil)
                    if err != nil {
                        handleDisconnect(deviceUUID) // 서비스 탐색 실패 시에도 Disconnect 처리
                        return
                    }

                    // 연결 시에도 서비스 UUID 확인 및 전송
                    for _, service := range services {
                        uuid := service.UUID().String()
                        if uuid != "00001801-0000-1000-8000-00805f9b34fb" {
                            if result.RSSI > RSSIThreshold {
                                handleConnect(uuid) // 연결 처리
                                fmt.Printf("Device Address: %s, UUID: %s, RSSI: %d\n", deviceUUID, uuid, result.RSSI)
                            } else {
                                handleDisconnect(uuid) // RSSI가 떨어지면 Disconnect
                            }
                        }
                    }
                }
            }

        })

        if err != nil {
            log.Printf("Error restarting BLE scan: %v", err)
        }
    }
}

func main() {
    fmt.Println("Initializing BLE adapter...")
    must("enable BLE stack", adapter.Enable())

    go restartScan()

    select {}
}

func must(action string, err error) {
    if err != nil {
        log.Fatalf("Failed to %s: %v", action, err)
    }
}

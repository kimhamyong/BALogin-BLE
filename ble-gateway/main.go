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

// 연결된 장치의 MAC 주소와 UUID를 저장하는 맵
var connectedDevices = make(map[string]string) // MAC 주소 -> UUID 매핑
var lastSeen = make(map[string]time.Time)      // MAC 주소 -> 마지막으로 감지된 시간
var mu sync.Mutex

// RSSI 임계값
const RSSIThreshold = -90

// 연결 및 해제 시 호출할 URL
const connectURL = "https://6e316739-cad9-44a2-8fae-647a4b9f2444.mock.pstmn.io/connected"
const disconnectURL = "https://6e316739-cad9-44a2-8fae-647a4b9f2444.mock.pstmn.io/disconnected"

// 타임아웃 설정 (30초 동안 감지되지 않으면 disconnect 처리)
const timeoutDuration = 30 * time.Second

// 신호 감지 주기를 3초로 설정
const scanInterval = 3 * time.Second

// POST 요청으로 UUID 전송 함수
func sendPostRequest(url string, uuid string) {
    jsonData := map[string]string{"uuid": uuid} // UUID 데이터를 JSON 형식으로 준비
    jsonValue, _ := json.Marshal(jsonData)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue)) // POST 요청 전송
    if err != nil {
        log.Printf("Failed to send POST request: %v", err)
        return
    }
    defer resp.Body.Close()
    fmt.Printf("Sent POST request to %s with UUID: %s, Status: %s\n", url, uuid, resp.Status)
}

// BLE 장치 연결이 끊겼을 때 처리하는 함수
func handleDisconnect(macAddress string) {
    mu.Lock()
    defer mu.Unlock()

    if uuid, exists := connectedDevices[macAddress]; exists {
        fmt.Printf("Device %s disconnected.\n", uuid)
        delete(connectedDevices, macAddress) // MAC 주소에서 장치를 삭제
        delete(lastSeen, macAddress)         // 감지된 시간 기록 삭제
        sendPostRequest(disconnectURL, uuid) // UUID 전송
    }
}

// BLE 장치가 연결되었을 때 처리하는 함수
func handleConnect(macAddress string, uuid string) {
    mu.Lock()
    defer mu.Unlock()

    if _, exists := connectedDevices[macAddress]; !exists {
        fmt.Printf("Device %s connected.\n", uuid)
        connectedDevices[macAddress] = uuid
        lastSeen[macAddress] = time.Now()  // 현재 시간을 마지막 감지 시간으로 기록
        sendPostRequest(connectURL, uuid)  // UUID 전송
    } else {
        lastSeen[macAddress] = time.Now()  // 이미 연결된 장치의 경우, 마지막 감지 시간만 업데이트
    }
}

// BLE 스캔 재시작 및 최신 장치 상태 반영 함수
func restartScan() {
    for {
        time.Sleep(scanInterval) // 신호 감지 주기를 3초로 설정

        fmt.Println("Restarting BLE scan to refresh device states...")

        // 타임아웃 체크를 먼저 수행
        checkTimeouts()

        err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
            if result.LocalName() == "" {
                return
            }

            macAddress := result.Address.String() // MAC 주소 사용

            // 장치가 감지되면 타임아웃 처리를 막기 위해 lastSeen을 즉시 업데이트
            lastSeen[macAddress] = time.Now()

            device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
            if err != nil {
                handleDisconnect(macAddress) // 연결 실패 시 Disconnect 처리
                return
            }
            defer device.Disconnect()

            // 서비스 탐색을 통해 UUID를 확인
            services, err := device.DiscoverServices(nil)
            if err != nil {
                handleDisconnect(macAddress) // 서비스 탐색 실패 시에도 Disconnect 처리
                return
            }

            // UUID 추출 및 연결 처리
            for _, service := range services {
                uuid := service.UUID().String()
                if uuid != "00001801-0000-1000-8000-00805f9b34fb" { // 기본 서비스 UUID 필터링
                    if result.RSSI > RSSIThreshold {
                        handleConnect(macAddress, uuid) // 연결 처리
                        fmt.Printf("Device Address: %s, UUID: %s, RSSI: %d\n", macAddress, uuid, result.RSSI)
                    } else {
                        handleDisconnect(macAddress) // RSSI가 떨어지면 Disconnect
                    }
                }
            }

        })

        if err != nil {
            log.Printf("Error restarting BLE scan: %v", err)
        }
    }
}

// 장치의 타임아웃을 체크하고, 일정 시간 동안 감지되지 않은 장치를 disconnect 처리하는 함수
func checkTimeouts() {
    mu.Lock()
    defer mu.Unlock()

    currentTime := time.Now()

    for macAddress, lastSeenTime := range lastSeen {
        if currentTime.Sub(lastSeenTime) > timeoutDuration {
            // 장치가 타임아웃에 도달한 경우 disconnect 처리
            fmt.Printf("Device %s timed out (no signal for %v).\n", connectedDevices[macAddress], timeoutDuration)
            handleDisconnect(macAddress)
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

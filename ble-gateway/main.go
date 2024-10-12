package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    "tinygo.org/x/bluetooth"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

var adapter = bluetooth.DefaultAdapter

// 연결된 장치의 MAC 주소와 UUID를 저장하는 맵
var connectedDevices = make(map[string]string) // MAC 주소 -> UUID 매핑
var lastSeen = make(map[string]time.Time)      // MAC 주소 -> 마지막으로 감지된 시간
var mu sync.Mutex

// RSSI 임계값
const RSSIThreshold = -90

// gRPC 서버 주소
const grpcServerAddress = "localhost:50051"

// 타임아웃 설정 (30초 동안 감지되지 않으면 disconnect 처리)
const timeoutDuration = 30 * time.Second

// 신호 감지 주기를 3초로 설정
const scanInterval = 3 * time.Second

// gRPC 클라이언트 초기화
func newBLEServiceClient() pb.BLEServiceClient {
    conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    return pb.NewBLEServiceClient(conn)
}

// gRPC를 통해 장치 상태 전송 함수
func sendDeviceStatus(client pb.BLEServiceClient, uuid string, status string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // BLE 장치 상태를 서버에 전송
    res, err := client.SendDeviceStatus(ctx, &pb.DeviceStatus{
        Uuid:   uuid,
        Status: status,
    })

    if err != nil {
        log.Printf("Failed to send device status: %v", err)
        return
    }
    fmt.Printf("Response from server: %s\n", res.Message)
}

// BLE 장치 연결이 끊겼을 때 처리하는 함수
func handleDisconnect(macAddress string, client pb.BLEServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    if uuid, exists := connectedDevices[macAddress]; exists {
        fmt.Printf("Device %s disconnected.\n", uuid)
        delete(connectedDevices, macAddress) // MAC 주소에서 장치를 삭제
        delete(lastSeen, macAddress)         // 감지된 시간 기록 삭제
        sendDeviceStatus(client, uuid, "disconnected") // UUID와 상태를 서버에 전송
    }
}

// BLE 장치가 연결되었을 때 처리하는 함수
func handleConnect(macAddress string, uuid string, client pb.BLEServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    if _, exists := connectedDevices[macAddress]; !exists {
        fmt.Printf("Device %s connected.\n", uuid)
        connectedDevices[macAddress] = uuid
        lastSeen[macAddress] = time.Now()  // 현재 시간을 마지막 감지 시간으로 기록
        sendDeviceStatus(client, uuid, "connected") // UUID와 상태를 서버에 전송
    } else {
        lastSeen[macAddress] = time.Now()  // 이미 연결된 장치의 경우, 마지막 감지 시간만 업데이트
    }
}

// BLE 스캔 재시작 및 최신 장치 상태 반영 함수
func restartScan(client pb.BLEServiceClient) {
    for {
        time.Sleep(scanInterval) // 신호 감지 주기를 3초로 설정

        fmt.Println("Restarting BLE scan to refresh device states...")

        // 타임아웃 체크를 먼저 수행
        checkTimeouts(client)

        err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
            if result.LocalName() == "" {
                return
            }

            macAddress := result.Address.String() // MAC 주소 사용

            // 장치가 감지되면 타임아웃 처리를 막기 위해 lastSeen을 즉시 업데이트
            lastSeen[macAddress] = time.Now()

            device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
            if err != nil {
                handleDisconnect(macAddress, client) // 연결 실패 시 Disconnect 처리
                return
            }
            defer device.Disconnect()

            // 서비스 탐색을 통해 UUID를 확인
            services, err := device.DiscoverServices(nil)
            if err != nil {
                handleDisconnect(macAddress, client) // 서비스 탐색 실패 시에도 Disconnect 처리
                return
            }

            // UUID 추출 및 연결 처리
            for _, service := range services {
                uuid := service.UUID().String()
                if uuid != "00001801-0000-1000-8000-00805f9b34fb" { // 기본 서비스 UUID 필터링
                    if result.RSSI > RSSIThreshold {
                        handleConnect(macAddress, uuid, client) // 연결 처리
                        fmt.Printf("Device Address: %s, UUID: %s, RSSI: %d\n", macAddress, uuid, result.RSSI)
                    } else {
                        handleDisconnect(macAddress, client) // RSSI가 떨어지면 Disconnect
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
func checkTimeouts(client pb.BLEServiceClient) {
    mu.Lock()
    defer mu.Unlock()

    currentTime := time.Now()

    for macAddress, lastSeenTime := range lastSeen {
        if currentTime.Sub(lastSeenTime) > timeoutDuration {
            // 장치가 타임아웃에 도달한 경우 disconnect 처리
            fmt.Printf("Device %s timed out (no signal for %v).\n", connectedDevices[macAddress], timeoutDuration)
            handleDisconnect(macAddress, client)
        }
    }
}

func main() {
    fmt.Println("Initializing BLE adapter...")
    must("enable BLE stack", adapter.Enable())

    // gRPC 클라이언트 생성
    client := newBLEServiceClient()

    go restartScan(client)

    // 블로킹 상태로 대기
    select {}
}

func must(action string, err error) {
    if err != nil {
        log.Fatalf("Failed to %s: %v", action, err)
    }
}

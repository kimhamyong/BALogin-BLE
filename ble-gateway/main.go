package main

import (
    "fmt"
    "log"
    "tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

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

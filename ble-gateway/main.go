package main

import (
    "fmt"
    "log"
    "tinygo.org/x/bluetooth"
    "ble-gateway/handler"
    "ble-gateway/ble"
)

var adapter = bluetooth.DefaultAdapter

func main() {
    fmt.Println("Starting program...")

    fmt.Println("Initializing BLE adapter...")
    must("enable BLE stack", adapter.Enable())
    fmt.Println("BLE adapter initialized.")

    // gRPC 클라이언트 생성
    fmt.Println("Creating gRPC client...")
    client := handler.NewBLEServiceClient()
    fmt.Println("gRPC client created.")

    // BLE 스캔 시작
    fmt.Println("Starting BLE scan...")
    go ble.RestartScan(client)

    fmt.Println("Waiting for BLE scan results...")
    select {}
}

func must(action string, err error) {
    if err != nil {
        log.Fatalf("Failed to %s: %v", action, err)
    }
    fmt.Printf("%s succeeded.\n", action)
}

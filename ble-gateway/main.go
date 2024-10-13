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

    must("enable BLE stack", adapter.Enable())
    fmt.Println("BLE adapter initialized.")

    client := handler.ServiceClient()
    fmt.Println("gRPC client created.")

    fmt.Println("Starting BLE scan...")
    go ble.RestartScan(client)

    fmt.Println("Waiting for server request...")
    go handler.ServiceServer()
    select {}
}

func must(action string, err error) {
    if err != nil {
        log.Fatalf("Failed to %s: %v", action, err)
    }
    fmt.Printf("%s succeeded.\n", action)
}

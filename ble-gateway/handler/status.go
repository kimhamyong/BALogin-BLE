package handler

import (
    "context"
    "log"
    "time"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

// gRPC server address
const BaloginServerAddress = "localhost:50051" // for testing

// Function to create a gRPC client
func ServiceClient() pb.DeviceServiceClient {
    conn, err := grpc.Dial(BaloginServerAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Printf("Failed to connect to gRPC server: %v", err)
        return nil // Return nil if the connection fails, allowing the caller to handle the error
    }
    return pb.NewDeviceServiceClient(conn)
}

// Function to send device status via gRPC
func SendDeviceStatus(client pb.DeviceServiceClient, uuid string, status int32) {
    if client == nil {
        log.Printf("Client is not initialized")
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a 5-second timeout
    defer cancel()

    // Send BLE device status to the server
    res, err := client.SendDeviceStatus(ctx, &pb.DeviceStatus{
        Uuid:   uuid,
        Status: status,
    })

    if err != nil {
        log.Printf("Failed to send device status: %v", err)
        return
    }
    log.Printf("Response from server: %s\n", res.Message)
}

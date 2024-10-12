package main

import (
    "context"
    "log"
    "time"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

// gRPC 서버 주소
const grpcServerAddress = "localhost:50051"

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
    log.Printf("Response from server: %s\n", res.Message)
}

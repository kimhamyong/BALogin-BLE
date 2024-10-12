package handler

import (
    "context"
    "log"
    "time"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

// gRPC 서버 주소
const grpcServerAddress = "localhost:50053"

// gRPC 클라이언트 생성 함수
func NewBLEServiceClient() pb.DeviceServiceClient {
    conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Printf("Failed to connect to gRPC server: %v", err)
        return nil // 실패한 경우 nil을 반환하여 호출자가 처리할 수 있도록 함
    }
    return pb.NewDeviceServiceClient(conn)
}

// gRPC를 통해 장치 상태 전송 함수
func SendDeviceStatus(client pb.DeviceServiceClient, uuid string, status int32) {
    if client == nil {
        log.Printf("Client is not initialized")
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 타임아웃 5초로 설정
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

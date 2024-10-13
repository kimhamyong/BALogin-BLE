package handler

import (
    "context"
    "fmt"
    "log"
    "net"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"  // 생성된 프로토콜 버퍼 코드
)

// DeviceServiceServer 구조체 정의
type server struct {
    pb.UnimplementedDeviceServiceServer
}

// RequestUnusedUUID: 서버에서 UUID 요청이 오면 호출되는 함수
func (s *server) RequestUnusedUUID(ctx context.Context, req *pb.UUIDRequest) (*pb.Response, error) {
    fmt.Println("Server request received.")  // 서버 요청을 알리는 출력
    return &pb.Response{Message: "UUID request processed"}, nil
}

// gRPC 서버 시작 함수
func ServiceServer() {
    // gRPC 서버 리스너 설정
    lis, err := net.Listen("tcp", ":50052")  // 50052 포트에서 대기
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterDeviceServiceServer(grpcServer, &server{})  // 서비스 핸들러 등록

    fmt.Println("gRPC server is running on port 50052...")  // 서버가 실행 중임을 알림
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

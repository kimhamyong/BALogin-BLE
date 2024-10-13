package handler

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net"

    _ "github.com/mattn/go-sqlite3"  // SQLite3 드라이버
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

// DeviceServiceServer 구조체 정의
type server struct {
    pb.UnimplementedDeviceServiceServer
}

// SQLite 데이터베이스 열기
func openDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./ble.db")
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    return db, nil
}

// is_active가 0인 UUID 찾기
func findInactiveUUID(db *sql.DB) (string, error) {
    var uuid string
    query := `SELECT uuid FROM devices WHERE is_active = 0 LIMIT 1`
    err := db.QueryRow(query).Scan(&uuid)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("There are not enough devices available.") //사용가능한 ble 기기가 없는 경우
        }
        return "", fmt.Errorf("failed to query inactive UUID: %v", err)
    }
    return uuid, nil
}

// UUID의 is_active 값을 1로 업데이트하는 함수
func updateUUIDStatusToActive(db *sql.DB, uuid string) error {
    query := `UPDATE devices SET is_active = 1 WHERE uuid = ?`
    _, err := db.Exec(query, uuid)
    if err != nil {
        return fmt.Errorf("failed to update UUID status: %v", err)
    }
    return nil
}

// RequestUnusedUUID: 서버에서 UUID 요청이 오면 호출되는 함수
func (s *server) RequestUnusedUUID(ctx context.Context, req *pb.UUIDRequest) (*pb.Response, error) {
    fmt.Println("Server request received.")

    // 데이터베이스 열기
    db, err := openDB()
    if err != nil {
        log.Printf("Failed to open database: %v", err)
        return &pb.Response{Message: "Failed to process request"}, err
    }
    defer db.Close()

    // is_active가 0인 UUID 찾기
    uuid, err := findInactiveUUID(db)
    if err != nil {
        log.Printf("Failed to find inactive UUID: %v", err)
        return &pb.Response{Message: "No inactive UUID found"}, err
    }

    // UUID가 성공적으로 처리되었음을 응답
    responseMessage := fmt.Sprintf("UUID: %s", uuid)
    fmt.Println(responseMessage)

    // UUID의 is_active 값을 1로 업데이트
    err = updateUUIDStatusToActive(db, uuid)
    if err != nil {
        log.Printf("Failed to update UUID status to active: %v", err)
        return &pb.Response{Message: "Failed to update UUID status"}, err
    }

    return &pb.Response{Message: responseMessage}, nil
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

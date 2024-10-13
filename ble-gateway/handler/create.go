package handler

import (
    "context"
    "fmt"
    "log"
    "net"
    "ble-gateway/db"
    "google.golang.org/grpc"
    pb "ble-gateway/proto"
)

// DeviceServiceServer structure definition
type server struct {
    pb.UnimplementedDeviceServiceServer
}

// RequestUnusedUUID: Function called when a UUID request is made to the server
func (s *server) RequestUnusedUUID(ctx context.Context, req *pb.UUIDRequest) (*pb.Response, error) {
    fmt.Println("Request to server.")

    // Call the service to activate and process the UUID
    uuid, err := db.GetAndActivateUUID()
    if err != nil {
        log.Printf("Failed to handle UUID: %v", err)
        return &pb.Response{Message: "Failed to process request"}, err
    }

    responseMessage := fmt.Sprintf("%s", uuid)
    return &pb.Response{Message: responseMessage}, nil
}

// ServiceServer: Function to start the gRPC server
func ServiceServer() {
    // Set up gRPC server listener
    lis, err := net.Listen("tcp", ":50052") // Waiting on port 50052
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterDeviceServiceServer(grpcServer, &server{}) // Register the service handler

    fmt.Println("Running on port 50052...") // Notify that the server is running
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

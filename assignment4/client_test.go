package main

import (
	"context"
	"testing"

	"net"
	pb "path/to/your/generated/protos/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic("Server exited with error")
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// Test client AddUser function.
func TestClientAddUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	user := &pb.User{Name: "Eve", Email: "eve@example.com"}
	resp, err := client.AddUser(ctx, user)
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}
	if resp.Id == 0 {
		t.Errorf("Expected valid user ID, got %d", resp.Id)
	}
}

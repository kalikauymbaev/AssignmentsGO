package main

import (
	"context"
	"log"
	"net"

	pb "path/to/your/generated/protos/user"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users  []pb.User // Simulated database of users.
	nextID int32     // Next user ID to assign.
}

func (s *server) AddUser(ctx context.Context, user *pb.User) (*pb.UserResponse, error) {
	user.Id = s.nextID
	s.users = append(s.users, *user)
	s.nextID++
	return &pb.UserResponse{Id: user.Id}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	for _, user := range s.users {
		if user.Id == req.Id {
			return &user, nil
		}
	}
	return nil, grpc.Errorf(grpc.Code().NotFound, "User not found")
}

func (s *server) ListUsers(ctx context.Context, empty *pb.EmptyRequest) (*pb.UserList, error) {
	return &pb.UserList{Users: s.users}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, &server{})
	log.Println("Server listening at", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

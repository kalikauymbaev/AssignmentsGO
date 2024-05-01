package main

import (
	"context"
	"log"

	pb "path/to/your/generated/protos/user"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	// Test AddUser
	user := &pb.User{Name: "John Doe", Email: "john@example.com"}
	res, err := c.AddUser(context.Background(), user)
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("Added user ID: %d", res.Id)

	// Test GetUser
	getUser, err := c.GetUser(context.Background(), &pb.UserRequest{Id: res.Id})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("Got user: %s", getUser.Name)

	// Test ListUsers
	list, err := c.ListUsers(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not list users: %v", err)
	}
	for _, u := range list.Users {
		log.Printf("User: %s, Email: %s", u.Name, u.Email)
	}
}

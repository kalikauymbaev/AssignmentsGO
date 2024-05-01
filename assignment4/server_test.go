package main

import (
	"context"
	"testing"

	pb "path/to/your/generated/protos/user"
)

// Create a new server instance for testing.
func setupServer() *server {
	return &server{
		users:  []pb.User{},
		nextID: 1,
	}
}

// Test adding a user.
func TestAddUser(t *testing.T) {
	s := setupServer()
	user := &pb.User{Name: "Alice", Email: "alice@example.com"}
	resp, err := s.AddUser(context.Background(), user)
	if err != nil {
		t.Errorf("AddUser returned an error: %v", err)
	}
	if resp.Id != 1 {
		t.Errorf("Expected user ID 1, got %d", resp.Id)
	}
}

// Test retrieving an existing user.
func TestGetUser(t *testing.T) {
	s := setupServer()
	user := &pb.User{Name: "Bob", Email: "bob@example.com"}
	s.AddUser(context.Background(), user)
	req := &pb.UserRequest{Id: 1}
	resp, err := s.GetUser(context.Background(), req)
	if err != nil {
		t.Errorf("GetUser returned an error: %v", err)
	}
	if resp.Name != "Bob" {
		t.Errorf("Expected user name Bob, got %s", resp.Name)
	}
}

// Test listing users.
func TestListUsers(t *testing.T) {
	s := setupServer()
	user1 := &pb.User{Name: "Carol", Email: "carol@example.com"}
	user2 := &pb.User{Name: "Dave", Email: "dave@example.com"}
	s.AddUser(context.Background(), user1)
	s.AddUser(context.Background(), user2)
	resp, err := s.ListUsers(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		t.Errorf("ListUsers returned an error: %v", err)
	}
	if len(resp.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(resp.Users))
	}
}

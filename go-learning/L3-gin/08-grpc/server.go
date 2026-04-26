package main

import (
	"context"
	"fmt"
	"grpc/pb"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	user   map[int64]*pb.User
	mu     sync.RWMutex
	nextID int64
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		user: make(map[int64]*pb.User),
	}
}
func (s *UserServiceServer) initSampleData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.user[1] = &pb.User{
		Id:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Age:      25,
		Active:   true,
		Tags:     []string{"admin", "developer"},
		Metadata: map[string]string{
			"department": "engineering",
			"location":   "Beijing",
		},
	}
	s.user[2] = &pb.User{
		Id:       2,
		Username: "bob",
		Email:    "bob@example.com",
		Age:      30,
		Active:   true,
		Tags:     []string{"user", "tester"},
		Metadata: map[string]string{
			"department": "qa",
			"location":   "Shanghai",
		},
	}
	s.nextID = 3
}
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.user[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}
	return user, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Username == "" || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Username and email are required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	user := &pb.User{
		Id:       s.nextID,
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Age,
		Active:   true,
		Tags:     req.Tags,
		Metadata: req.Metadata,
	}
	s.user[s.nextID] = user
	s.nextID++
	return &pb.CreateUserResponse{
		User:    user,
		Success: true,
		Message: "User created successfully",
	}, nil
}
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	users := make([]*pb.User, 0, len(s.user))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for id, user := range s.user {
		if id >= (int64(req.Page-1)*int64(req.PageSize)) && id < (int64(req.Page)*int64(req.PageSize)) {
			users = append(users, user)
		}
	}
	return &pb.ListUserResponse{
		Users:    users,
		Total:    int32(len(s.user)),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.user[int64(req.Id)]
	if !ok {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age != 0 {
		user.Age = req.Age
	}
	if req.Tags != nil {
		user.Tags = req.Tags
	}
	if req.Metadata != nil {
		user.Metadata = req.Metadata
	}
	return &pb.UpdateUserResponse{
		User:    user,
		Success: true,
		Message: "User updated successfully",
	}, nil
}
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.user[int64(req.Id)]
	if !ok {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	delete(s.user, int64(req.Id))
	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
	s.mu.RLock()
	user := make([]*pb.User, 0, len(s.user))
	for _, v := range s.user {
		user = append(user, v)
	}
	s.mu.RUnlock()
	limit := int(req.Limit)
	if limit <= 0 {
		limit = len(user)
	}
	if limit > len(user) {
		limit = len(user)
	}
	interval := time.Duration(req.IntervalMs) * time.Millisecond
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	for i := 0; i < limit; i++ {
		if err := stream.Send(user[i]); err != nil {
			return err
		}
		if i < limit-1 {
			time.Sleep(interval)
		}
	}
	return nil

}
func (s *UserServiceServer) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
	var createdUsers []*pb.User
	successCount := 0
	failCount := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			break
		}
		resp, err := s.CreateUser(stream.Context(), req)
		if err != nil {
			failCount++
			log.Printf("Failed to create user: %v", err)
			continue
		}
		createdUsers = append(createdUsers, resp.User)
		successCount++
	}
	return stream.SendAndClose(&pb.BatchCreateUsersResponse{
		Users:        createdUsers,
		SuccessCount: int32(successCount),
		FailCount:    int32(failCount),
		Message:      fmt.Sprintf("Batch create completed: %d success, %d fail", successCount, failCount),
	})
}
func startServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	userService := NewUserServiceServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("gRPC server listening on %s", port)
	log.Println("Available methods:")
	log.Println("  - GetUser")
	log.Println("  - CreateUser")
	log.Println("  - ListUsers")
	log.Println("  - UpdateUser")
	log.Println("  - DeleteUser")
	log.Println("  - StreamUsers (server streaming)")
	log.Println("  - BatchCreateUsers (client streaming)")
	log.Println("  - ChatUsers (bidirectional streaming)")

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

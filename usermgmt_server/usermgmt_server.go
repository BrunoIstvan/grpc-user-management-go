package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/BrunoIstvan/grpc-user-management-go/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {

	log.Printf("Receive: %v", in.GetName())

	var user_id int32 = int32(rand.Intn(1000))

	return &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}, nil

}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterUserManagementServer(serv, &UserManagementServer{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := serv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

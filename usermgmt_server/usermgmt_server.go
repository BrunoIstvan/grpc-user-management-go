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
	user_list *pb.UserList
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

func (server *UserManagementServer) Run() error {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterUserManagementServer(serv, server)

	log.Printf("Server listening at %v", lis.Addr())

	return serv.Serve(lis)

}

func (serv *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {

	log.Printf("Receive: %v", in.GetName())

	var user_id int32 = int32(rand.Intn(1000))

	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}

	serv.user_list.Users = append(serv.user_list.Users, created_user)

	return created_user, nil

}

func (serv *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUserParams) (*pb.UserList, error) {

	return serv.user_list, nil

}

func main() {

	var user_mgmt_server *UserManagementServer = NewUserManagementServer()

	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}

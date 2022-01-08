package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/BrunoIstvan/grpc-user-management-go/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port      = ":50051"
	user_file = "user.json"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	// user_list *pb.UserList
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		// user_list: &pb.UserList{},
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

	readBytes, err := ioutil.ReadFile(user_file)

	var users_list *pb.UsersList = &pb.UsersList{}

	var user_id int32 = int32(rand.Intn(1000))

	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}

	if err != nil {

		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")

			users_list.Users = append(users_list.Users, created_user)

			// jsonBytes, err := protojson.Marshal(users_list)

			// if err != nil {
			// 	log.Fatalf("JSON Marshaling failed: %v", err)
			// }

			// if err := ioutil.WriteFile(user_file, jsonBytes, 0664); err != nil {
			// 	log.Fatalf("Failed write to file: %v", err)
			// }

			marshalAndWriteFile(users_list)

			return created_user, nil

		} else {
			log.Fatalln("Error reading file: %v", err)
		}

	}

	if err := protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err)
	}
	users_list.Users = append(users_list.Users, created_user)

	// jsonBytes, err := protojson.Marshal(users_list)

	// if err != nil {
	// 	log.Fatalf("JSON Marshaling failed: %v", err)
	// }

	// if err := ioutil.WriteFile(user_file, jsonBytes, 0664); err != nil {
	// 	log.Fatalf("Failed write to file: %v", err)
	// }

	marshalAndWriteFile(users_list)

	return created_user, nil

}

func marshalAndWriteFile(users_list *pb.UsersList) {
	jsonBytes, err := protojson.Marshal(users_list)

	if err != nil {
		log.Fatalf("JSON Marshaling failed: %v", err)
	}

	if err := ioutil.WriteFile(user_file, jsonBytes, 0664); err != nil {
		log.Fatalf("Failed write to file: %v", err)
	}
}

func (serv *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUserParams) (*pb.UsersList, error) {

	jsonBytes, err := ioutil.ReadFile(user_file)
	if err != nil {
		log.Fatalf("Failed read from file: %v", err)
	}

	var users_list *pb.UsersList = &pb.UsersList{}

	if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err)
	}

	return users_list, nil

}

func main() {

	var user_mgmt_server *UserManagementServer = NewUserManagementServer()

	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}

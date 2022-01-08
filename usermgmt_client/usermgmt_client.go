package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/BrunoIstvan/grpc-user-management-go/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)

	new_users["User 1"] = 38
	new_users["User 2"] = 8
	new_users["User 3"] = 29
	new_users["User 4"] = 66

	for name, age := range new_users {

		resp, err := client.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("Could not create user: %v", err)
		}

		log.Printf(`User Details: [Name: %s, Age: %d, Id: %d]`, resp.GetName(), resp.GetAge(), resp.GetId())

	}

	params := &pb.GetUserParams{}

	resp, err := client.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("Could not retrieve users: %v", err)
	}

	log.Print("\nUser List: \n")
	fmt.Printf("resp.GetUsers(): %v\n", resp.GetUsers())

}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/BrunoIstvan/grpc-user-management-go/usermgmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port      = ":50051"
	user_file = "user.json"
)

type UserManagementServer struct {
	conn *pgx.Conn

	pb.UnimplementedUserManagementServer
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
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

	log.Printf("Receive to create: %v", in.GetName())

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}

	tx, err := serv.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users (name, age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("Failed to create user into database: %v", err)
	}

	tx.Commit(context.Background())

	return created_user, nil

}

func (serv *UserManagementServer) UpdateUser(ctx context.Context, in *pb.User) (*pb.Message, error) {

	log.Printf("Receive to update: %v", in.GetName())

	updated_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: in.Id}

	tx, err := serv.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), "update users set name = $1, age = $2 where id = $3", updated_user.Name, updated_user.Age, updated_user.Id)
	if err != nil {
		log.Fatalf("Failed to update user into database: %v", err)
	}

	tx.Commit(context.Background())

	return &pb.Message{Content: "Updated successfuly"}, nil

}

func (serv *UserManagementServer) DeleteUser(ctx context.Context, in *pb.Number) (*pb.Message, error) {

	log.Printf("Receive to delete: %v", in.GetId())

	tx, err := serv.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), "delete from users where id = $1", in.GetId())
	if err != nil {
		log.Fatalf("Failed to delete user into database: %v", err)
	}

	tx.Commit(context.Background())

	return &pb.Message{Content: "Deleted successfuly"}, nil

}

func (serv *UserManagementServer) GetAllUsers(ctx context.Context, in *pb.GetUserParams) (*pb.UsersList, error) {

	var users_list *pb.UsersList = &pb.UsersList{}

	rows, err := serv.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)

	}

	return users_list, nil

}

func (serv *UserManagementServer) GetUserById(ctx context.Context, in *pb.Number) (*pb.User, error) {

	rows, err := serv.conn.Query(context.Background(), "select * from users where id = $1", in.GetId())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		return &user, nil

	}

	return nil, nil

}

func main() {

	database_url := "postgres://root:root@localhost:5432/root"

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())

	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	user_mgmt_server.conn = conn

	createTableSql := `
	create table if not exists users(
		id SERIAL PRIMARY KEY,
		name text,
		age int
	);`

	_, err = conn.Exec(context.Background(), createTableSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create table: %v\n", err)
		os.Exit(1)
	}

	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}

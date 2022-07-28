package main

import (
	"context"
	pb "example.com/go-usermgmt-grpc/usermgmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Println("received %v", in.GetName())

	createSQL := `CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name text,
	age integer
)`
	_, err := s.conn.Exec(context.Background(), createSQL)
	if err != nil {
		log.Fatalf("failed to create table %v", err)
	}
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("failed to begin: %v", err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO users (name, age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("could not insert %v", err)
	}
	tx.Commit(context.Background())
	return created_user, nil
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Listening on port %s", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) GetUsers(ctx context.Context, userParams *pb.GetUserParams) (*pb.UserList, error) {
	users_list := &pb.UserList{}
	user := &pb.User{}
	rows, err := s.conn.Query(context.Background(), `select * from users`)
	if err != nil {
		log.Fatalf("could not read from table %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Age)
		users_list.Users = append(users_list.Users, user)
	}
	return users_list, nil
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

func main() {
	conn, err := pgx.Connect(context.Background(), `postgres://postgres:mysecretpassword@localhost:5432/postgres`)
	if err != nil {
		log.Fatalf("error connecting to DB: %v", err)
	}
	defer conn.Close(context.Background())
	user_mgt_server := NewUserManagementServer()
	user_mgt_server.conn = conn
	if err := user_mgt_server.Run(); err != nil {
		log.Fatalf("failed to run %v", err)
	}
}

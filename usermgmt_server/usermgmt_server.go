package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "example.com/go-usermgmt-grpc/usermgmt"
)

const (
	port     = ":50051"
	fileName = "users_list.json"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Println("received %v", in.GetName())
	var user_id = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	readBytes := readFromFile()
	users_list := &pb.UserList{}
	if readBytes != nil {
		if err := protojson.Unmarshal(readBytes, users_list); err != nil {
			log.Fatalf("failed unmarshaling bytes %v", err)
		}
	}
	users_list.Users = append(users_list.Users, created_user)
	writeFile(users_list)
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
	readBytes := readFromFile()
	if readBytes != nil {
		if err := protojson.Unmarshal(readBytes, users_list); err != nil {
			log.Fatalf("failed unmarshaling bytes %v", err)
		}
	}
	return users_list, nil
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

func main() {
	user_mgt_server := NewUserManagementServer()
	if err := user_mgt_server.Run(); err != nil {
		log.Fatalf("failed to run %v", err)
	}
}

func readFromFile() []byte {
	readBytes, err := os.ReadFile(fileName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to read file: %v", err)
	}
	return readBytes
}

func writeFile(usersList *pb.UserList) {
	jsonBytes, err := protojson.Marshal(usersList)
	if err != nil {
		log.Fatalf("failed marshalling json: %v", err)
	}

	err = ioutil.WriteFile(fileName, jsonBytes, 0664)
	if err != nil {
		log.Fatalf("failed to write to a file: %v", err)
	}
}

package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "example.com/go-usermgmt-grpc/usermgmt"
)

const (
	address = "localhost:50051"
)

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer connection.Close()

	c := pb.NewUserManagementClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	new_users := make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("error creating a new user %v", err)
		}
		log.Printf(`User details:
	Name %s
	Age %d
	Id %d`, r.GetName(), r.GetAge(), r.GetId())
	}

	params := &pb.GetUserParams{}
	list, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not get a list %v", err)
	}
	log.Printf("The list of users %v", list)
}

package main

import (
	"log"

	rpc "ffRPC/ffRPCUserLogin"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address     = "fps.coola.tv:50051"
	defaultName = "world"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("keys/server.crt", "fps.coola.tv")
	if err != nil {
		log.Fatalf("failed to create client TLS credentials: %v", err)
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := rpc.NewUserLoginServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Login(context.Background(), &rpc.LoginRequest{LoginInfo: "loginInfo"})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Printf("LoginResult: %s", r.LoginResult)
}

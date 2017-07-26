/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"log"
	"net"

	"ffRPC/ffRPCLoginServer"
	"ffRPC/ffRPCUserLogin"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:50051"
)

// userLoginService is used to implement ffRPCServiceServer.UserLoginServiceServer.
type userLoginService struct{}

// Login implements UserLoginServiceServer.Login
func (s *userLoginService) Login(ctx context.Context, in *ffRPCUserLogin.LoginRequest) (*ffRPCUserLogin.LoginReply, error) {
	return &ffRPCUserLogin.LoginReply{LoginResult: "Hello " + in.LoginInfo}, nil
}

// loginServerService is used to implement ffRPCServiceServer.LoginServerServiceServer.
type loginServerService struct{}

// Login implements LoginServerServiceServer.Login
func (s *loginServerService) Login(ctx context.Context, in *ffRPCLoginServer.LoginRequest) (*ffRPCLoginServer.LoginReply, error) {
	return &ffRPCLoginServer.LoginReply{LoginResult: "Hello " + in.LoginInfo}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server_no_passwd.key")
	if err != nil {
		log.Fatalf("failed to create server TLS credentials: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	ffRPCUserLogin.RegisterUserLoginServiceServer(s, &userLoginService{})
	ffRPCLoginServer.RegisterLoginServerServiceServer(s, &loginServerService{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

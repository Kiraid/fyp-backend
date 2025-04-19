package rpcserver

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"search.com/m/rpc_server/pb"
)


func Server(){
	productServer := NewProductServer()
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productServer)
	
	listener, err := net.Listen("tcp", ":8085")
	log.Print("Starting GRPC Server on Port :8085")
	if err != nil {
		log.Fatal("cannot start grpc server: ",err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot serve grpc server: ", err)
	}
}
package grpc_client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"fyp.com/m/grpc_client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(product *pb.Product){
	rpcHost := os.Getenv("RPC_HOST")
	rpcPort := os.Getenv("RPC_PORT")
	if rpcHost == "" || rpcPort == "" {
        log.Printf("gRPC environment variables not set: RPC_HOST=%s, RPC_PORT=%s", rpcHost, rpcPort)
        return
    }	
	address := fmt.Sprintf("%s:%s", rpcHost, rpcPort)
	log.Printf("→ dialing gRPC at %q", address)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Printf("Cannot connect to server: %v", err) 
        return
    }
	defer conn.Close()
	
	ProductClient := pb.NewProductServiceClient(conn)
	SendProduct(ProductClient, product)
}

func SendProduct(productClient pb.ProductServiceClient, product *pb.Product){
	req := &pb.SendProductRequest{Product: product}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err :=	productClient.SendProduct(ctx, req)
	if err != nil {
		log.Printf("Couldnot send the product: %s", err)
		return
	}
	log.Printf("Sent Product with id: %s", res.Id)
	
}
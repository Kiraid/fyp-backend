package rpcserver

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"search.com/m/rpc_server/pb"
	"search.com/m/storing"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	
}

func NewProductServer() *ProductServer {
	return &ProductServer{}
}

func (server *ProductServer) SendProduct(ctx context.Context,req *pb.SendProductRequest) (*pb.SendProductResponse, error){
	product := req.GetProduct()
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Errorf(codes.Canceled, "Request was canceled")
	}
	
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Errorf(codes.DeadlineExceeded, "Request was canceled")
	}
	prod := storing.NewProduct()
	prod.ID = int64(product.GetId())
	prod.Name = product.GetName()
	prod.Description = product.GetDescription()
	prod.ImagePath = product.GetImagePath()
	prod.UserID = int64(product.GetUserId())
	prod.Category_name = product.GetCategoryName()
	prod.Price = product.GetPrice()
	
	err := prod.Save()
	if err != nil {
		log.Printf("failed to save product to EKS: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save product")
	}
	return &pb.SendProductResponse{Id: fmt.Sprintf("%d", prod.ID)}, nil
	
}
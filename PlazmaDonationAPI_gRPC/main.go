package main

import (
	pb "PlazmaDonation/Gen_code"
	services "PlazmaDonation/Server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Println("Failed to listen server")
	}
	ser := grpc.NewServer()
	pb.RegisterUserServiceServer(ser, &services.Server{})
	if err := ser.Serve(lis); err != nil {
		log.Println("failed to serve")
	}
}

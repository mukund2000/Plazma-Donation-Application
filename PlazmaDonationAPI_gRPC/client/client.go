package main

import (
	"context"
	"log"
	"time"

	pb "PlazmaDonation/PlazmaDonationAPI_gRPC/Gen_code"
	"google.golang.org/grpc"
)

const (
	address = "localhost:3000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("did not connect:", err)
	}
	defer conn.Close()

	con := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//r, err := con.CreateUser(ctx, &pb.UserDetails{
	//	Name:       "Alex",
	//	Address:    "jkshdkaba",
	//	PhoneNum:   "14527845",
	//	UserType:   "Donor",
	//	DiseaseDes: "hjdshjfkas"})
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println("User 1:", r)
	//r1, err := con.CreateUser(ctx, &pb.UserDetails{
	//	Name:       "John",
	//	Address:    "hjbsdcuru",
	//	PhoneNum:   "12345678",
	//	UserType:   "Patient",
	//	DiseaseDes: "bhgtfrdesdc"})
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println("User 2:", r1)
	//
	//r2, err := con.Login(ctx, &pb.UserDetails{
	//	SecretCode: "9727887"})
	//if err != nil {
	//	log.Println("Unable to login user.", err)
	//}
	//log.Println("User gets logged in:", r2)

	//r3, err := con.SendRequest(ctx, &pb.UserRequest{
	//	Id:         "7131847",
	//	SecretCode: "9727887"})
	//if err != nil {
	//	log.Println("unable to send request", err)
	//}
	//log.Println(r3)

	//r4, err := con.AcceptRequest(ctx, &pb.UserRequest{
	//	Id:         "8498081",
	//	SecretCode: "9984059"})
	//if err != nil {
	//	log.Println("unable to accept request", err)
	//}
	//log.Println(r4)
	//
	//r5, err := con.CreateUser(ctx, &pb.UserDetails{
	//	Name:       "Krish",
	//	Address:    "bhjbxkskbjhsa",
	//	PhoneNum:   "4875135",
	//	UserType:   "Donor",
	//	DiseaseDes: "bjsrbjbhbs"})
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println("User 3:", r5)
	//
	//r6, err := con.CreateUser(ctx, &pb.UserDetails{
	//	Name:       "Fred",
	//	Address:    "vhvdhuvvdu",
	//	PhoneNum:   "145784214",
	//	UserType:   "Patient",
	//	DiseaseDes: "hsuhbhbsiu"})
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println("User 4:", r6)
	//
	//r7, err := con.SendRequest(ctx, &pb.UserRequest{
	//	Id:         "954425",
	//	SecretCode: "4941318"})
	//if err != nil {
	//	log.Println("unable to send request", err)
	//}
	//log.Println(r7)
	//
	//r8, err := con.CancelRequest(ctx, &pb.UserRequest{
	//	Id:         "1902081",
	//	SecretCode: "6122540"})
	//if err != nil {
	//	log.Println("unable to send request", err)
	//}
	//log.Println(r8)
	//
	//r9, err := con.GetUser(ctx, &pb.UserRequest{
	//	Id:         "3024728",
	//	SecretCode: "9431445"})
	//
	//if err != nil {
	//	log.Println("unable to access user", err)
	//}
	//log.Println(r9)

	r10, err := con.GetAllDonors(ctx, &pb.UserDetails{
		SecretCode: "6122540"})
	if err != nil {
		log.Println("unable to access Donors", err)
	}
	log.Println(r10)
}

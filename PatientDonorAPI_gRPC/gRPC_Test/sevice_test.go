package gRPC_Test

import (
	common "PlazmaDonation/PatientDonorAPI_gRPC/Common"
	pb "PlazmaDonation/PatientDonorAPI_gRPC/Gen_code"
	services "PlazmaDonation/PatientDonorAPI_gRPC/Server"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/http"
	"testing"
)

func userDialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, &services.Server{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
func getIdToken(email string) (string, string) {
	fireApp, err := common.GetFirebaseInstance()
	if err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	fireAuth, err := fireApp.Auth(ctx)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	token, err := fireAuth.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	genToken, err := fireAuth.CustomToken(ctx, token.UID)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	//generating id token from customtoken using token.uid
	values := map[string]string{"token": genToken, "returnSecureToken": "true"}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key="+common.ApiKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	idToken := fmt.Sprintf("%v", res["idToken"])
	return idToken, token.UID
}
func TestCreateUser(t *testing.T) {
	//change to localhost when locally testing
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithInsecure(), grpc.WithContextDialer(userDialer()))
	if err != nil {
		log.Println(err, "server")
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	client := pb.NewUserServiceClient(conn)
	request := &pb.UserDetails{
		Email:      "mukundrastogixyz@gmail.com",
		Name:       "Mukund",
		Address:    "bdhbjksbkj",
		UserType:   "Donor",
		DiseaseDes: "dshviwdhbkk",
		PhoneNum:   "39839872930",
	}
	response, err := client.CreateUser(ctx, request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	t.Log(response)
	t.Log("Owner user created")
}

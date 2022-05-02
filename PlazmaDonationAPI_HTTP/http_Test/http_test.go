package http_Test

import (
	pb "PlazmaDonationHTTP/GeneratedCode"
	service "PlazmaDonationHTTP/Server"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	request := &pb.UserDetails{
		Email:      "test01@gmail.com",
		Name:       "Test",
		Address:    "indore",
		UserType:   0,
		DiseaseDes: "sdbjbijads",
		PhoneNum:   "14585622265",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/createUser", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.CreateUser) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.UserDetails{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestLogin(t *testing.T) {
	request := &pb.LoginRequest{
		Email:    "mukundrastogixyz@gmail.com",
		Password: "Test123",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Login) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestDeleteUser(t *testing.T) {
	request := &pb.DeleteUserRequest{
		UserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/deleteUser", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.DeleteUser) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestGetUser(t *testing.T) {
	request := &pb.UserRequest{
		UserId:          "",
		RequestedUserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/getUser", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetUser) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.UserDetails{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestGetAllDonors(t *testing.T) {
	request := &pb.UserDetails{
		Id: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/getAllDonors", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetAllDonors) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.ListUser{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestGetAllPatients(t *testing.T) {
	request := &pb.UserDetails{
		Id: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/getAllPatients", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetAllPatients) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.ListUser{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestUpdateContactDetails(t *testing.T) {
	request := &pb.UserDetails{
		Id:         "",
		PhoneNum:   "",
		Address:    "",
		DiseaseDes: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/updateContactDetails", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateContactDetails) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestSendRequest(t *testing.T) {
	request := &pb.UserRequest{
		UserId:          "",
		RequestedUserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/sendRequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.SendRequest) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestAcceptRequest(t *testing.T) {
	request := &pb.UserRequest{
		UserId:          "",
		RequestedUserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/acceptRequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.AcceptRequest) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestCancelRequest(t *testing.T) {
	request := &pb.UserRequest{
		UserId:          "",
		RequestedUserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/cancelRequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.CancelRequest) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
func TestCancelConnection(t *testing.T) {
	request := &pb.UserRequest{
		UserId:          "",
		RequestedUserId: "",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	//t.Log(jsonData)
	req, err := http.NewRequest(http.MethodPost, "/cancelConnection", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.CancelConnection) //register function
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	data := json.NewDecoder(rr.Body)
	res := &pb.Success{}
	if err := data.Decode(res); err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}

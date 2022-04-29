package main

import (
	common "PlazmaDonationHTTP/Common"
	pb "PlazmaDonationHTTP/GeneratedCode"
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const idField = "Id"
const secretCodeField = "SecretCode"
const nameField = "Name"
const addressField = "Address"
const phoneNumField = "PhoneNum"
const userTypeField = "UserType"
const diseaseDesField = "DiseaseDes"
const requestUsersField = "RequestUsers"
const pendingUsersField = "PendingUsers"
const connectUsersField = "ConnectUsers"
const emailField = "Email"

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		loginRequest := &pb.LoginRequest{}
		if err := data.Decode(loginRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		cloudReq := struct {
			Email             string `json:"email,omitempty"`
			Password          string `json:"password,omitempty"`
			ReturnSecureToken bool   `json:"returnSecureToken,omitempty"`
		}{
			Email:             loginRequest.Email,
			Password:          loginRequest.Password,
			ReturnSecureToken: true,
		}
		jsonData, err := json.Marshal(cloudReq)
		if err != nil {
			log.Println(err)
			return
		}
		postUrl := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="
		resp, err := http.Post(postUrl+common.ApiKey, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println(err)
			return
		}
		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		if res["error"] != nil {
			log.Println(common.AuthInvalidMsg)
			return
		}
		loginResponse := pb.Success{Name: "Login Success"}
		if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		createUserRequest := &pb.UserDetails{}
		if err := data.Decode(createUserRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		fireAuth, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		params := (&auth.UserToCreate{}).
			Email(createUserRequest.Email).
			Password("Test123")
		authUser, appErr := fireAuth.CreateUser(ctx, params)
		if appErr != nil {
			log.Println(appErr)
			return
		}
		if _, err := fireClient.Collection(common.CollectionUsers).Doc(authUser.UID).Set(ctx, map[string]interface{}{
			idField:           authUser.UID,
			emailField:        createUserRequest.Email,
			addressField:      createUserRequest.Address,
			nameField:         createUserRequest.Name,
			pendingUsersField: createUserRequest.PendingUsers,
			connectUsersField: createUserRequest.ConnectUsers,
			secretCodeField:   strconv.Itoa(rand.Intn(10000000)),
			diseaseDesField:   createUserRequest.DiseaseDes,
			userTypeField:     createUserRequest.UserType,
			phoneNumField:     createUserRequest.PhoneNum,
			requestUsersField: createUserRequest.RequestUsers,
		}); err != nil {
			log.Println(common.AddErrorMsg)
			return
		}
		userDoc, err := common.GetUserDocument(fireClient, authUser.UID)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if err := json.NewEncoder(w).Encode(userDoc); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func getUser(w http.ResponseWriter, r *http.Request) {

}

func getAllDonors(w http.ResponseWriter, r *http.Request) {

}

func getAllPatients(w http.ResponseWriter, r *http.Request) {

}

func updateContactDetails(w http.ResponseWriter, r *http.Request) {

}

func sendRequest(w http.ResponseWriter, r *http.Request) {

}

func acceptRequest(w http.ResponseWriter, r *http.Request) {

}

func cancelRequest(w http.ResponseWriter, r *http.Request) {

}

func cancelConnection(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/login", login)
	http.HandleFunc("/getUser", getUser)
	http.HandleFunc("/getDonors", getAllDonors)
	http.HandleFunc("/getPatients", getAllPatients)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/updateContactDetails", updateContactDetails)
	http.HandleFunc("/sendRequest", sendRequest)
	http.HandleFunc("/acceptRequest", acceptRequest)
	http.HandleFunc("/cancelRequest", cancelRequest)
	http.HandleFunc("/cancelConnection", cancelConnection)
	log.Println(http.ListenAndServe("192.168.29.118:8080", nil))
}

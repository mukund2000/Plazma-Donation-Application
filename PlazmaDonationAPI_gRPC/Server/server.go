package services

import (
	common "PlazmaDonation/PlazmaDonationAPI_gRPC/Common"
	pb "PlazmaDonation/PlazmaDonationAPI_gRPC/Gen_code"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"net/http"
	"sync"

	"log"
	"math/rand"
	"strconv"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

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

var m sync.Mutex

func (s *Server) CreateUser(ctx context.Context, user *pb.UserDetails) (*pb.UserDetails, error) {
	fireAuth, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	params := (&auth.UserToCreate{}).
		Email(user.Email).
		Password("Test123")
	authUser, appErr := fireAuth.CreateUser(ctx, params)
	if appErr != nil {
		log.Println(appErr)
		return nil, errors.New(common.InvalidLoginErrorMsg)
	}
	if _, err := fireClient.Collection(common.CollectionUsers).Doc(authUser.UID).Set(ctx, map[string]interface{}{
		idField:           authUser.UID,
		emailField:        user.Email,
		addressField:      user.Address,
		nameField:         user.Name,
		pendingUsersField: user.PendingUsers,
		connectUsersField: user.ConnectUsers,
		secretCodeField:   strconv.Itoa(rand.Intn(10000000)),
		diseaseDesField:   user.DiseaseDes,
		userTypeField:     user.UserType,
		phoneNumField:     user.PhoneNum,
		requestUsersField: user.RequestUsers,
	}); err != nil {
		log.Println(err)
		return nil, errors.New(common.AddErrorMsg)
	}
	userDoc, err := common.GetUserDocument(fireClient, authUser.UID)
	if err != nil {
		log.Println(err)
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	return userDoc, nil
}

func (s *Server) Login(_ context.Context, request *pb.LoginRequest) (*pb.Success, error) {
	cloudReq := struct {
		Email             string `json:"email,omitempty"`
		Password          string `json:"password,omitempty"`
		ReturnSecureToken bool   `json:"returnSecureToken,omitempty"`
	}{
		Email:             request.Email,
		Password:          request.Password,
		ReturnSecureToken: true,
	}
	jsonData, err := json.Marshal(cloudReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	postUrl := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="
	resp, err := http.Post(postUrl+common.ApiKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Println(err)
		return nil, errors.New(common.DecodeErrorMsg)
	}
	if res["error"] != nil {
		return nil, errors.New(common.AuthInvalidMsg)
	}
	loginResponse := pb.Success{Name: "Login Success"}
	return &loginResponse, nil

}

func (s *Server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.Success, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	_, err = fireClient.Collection(common.CollectionUsers).Doc(in.UserId).Delete(ctx)
	if err != nil {
		return nil, errors.New(common.FirebaseErrorMsg)
	}
	_, err = common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	return &pb.Success{Name: "User deleted"}, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	requestUSerDoc, err := common.GetUserDocument(fireClient, in.RequestedUserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	connectedUsers := userDoc.ConnectUsers
	found := false
	for _, user := range connectedUsers {
		if user == in.RequestedUserId {
			found = true
		}
	}
	userDetails := pb.UserResponse{}
	userDetails.Id = requestUSerDoc.Id
	userDetails.Name = requestUSerDoc.Name
	if found {
		userDetails.Address = requestUSerDoc.Address
		userDetails.PhoneNum = requestUSerDoc.PhoneNum
	}
	return &userDetails, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *pb.UserDetails) (*pb.UserDetails, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.Id)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	updatedPhoneNum := in.PhoneNum
	updatedAddress := in.Address
	updatedDiseaseDes := in.DiseaseDes
	if updatedPhoneNum == common.EmptyString || updatedPhoneNum == userDoc.PhoneNum {
		updatedPhoneNum = userDoc.PhoneNum
	}
	if updatedAddress == common.EmptyString || updatedAddress == userDoc.Address {
		updatedAddress = userDoc.Address
	}
	if updatedDiseaseDes == common.EmptyString || updatedDiseaseDes == userDoc.DiseaseDes {
		updatedDiseaseDes = userDoc.DiseaseDes
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(in.Id).Update(ctx, []firestore.Update{
		{
			Path:  phoneNumField,
			Value: in.PhoneNum,
		},
		{
			Path:  addressField,
			Value: in.Address,
		},
		{
			Path:  diseaseDesField,
			Value: in.DiseaseDes,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	updatedUserDoc, err := common.GetUserDocument(fireClient, in.Id)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	return updatedUserDoc, nil
}

func (s *Server) GetAllDonors(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	//user := in
	//if userStruct, found := UsersByCode[user.SecretCode]; found {
	//	var userData []*pb.UserResponse
	//	tempUsertype := userStruct.UserType
	//	if tempUsertype != pb.Type_donor {
	//		for key, value := range Donors {
	//			userDetails := pb.UserResponse{}
	//			userDetails.Id = value.Id
	//			userDetails.Name = value.Name
	//			if _, found := userStruct.ConnectUsers[key]; found {
	//				userDetails.Address = value.Address
	//				userDetails.PhoneNum = value.PhoneNum
	//			}
	//			userData = append(userData, &userDetails)
	//		}
	//		return &pb.ListUser{Users: userData}, nil
	//	}
	//}
	return nil, nil
}

func (s *Server) GetAllPatients(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	//user := in
	//if userStruct, found := UsersByCode[user.SecretCode]; found {
	//	var userData []*pb.UserResponse
	//	tempUsertype := userStruct.UserType
	//	if tempUsertype != pb.Type_patient {
	//		for key, value := range Patients {
	//			userDetails := pb.UserResponse{}
	//			userDetails.Id = value.Id
	//			userDetails.Name = value.Name
	//			if _, found := userStruct.ConnectUsers[key]; found {
	//				userDetails.Address = value.Address
	//				userDetails.PhoneNum = value.PhoneNum
	//			}
	//			userData = append(userData, &userDetails)
	//		}
	//		return &pb.ListUser{Users: userData}, nil
	//	}
	//}
	return nil, nil
}

func (s *Server) SendRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	//user := in
	//if sender, found := UsersByCode[user.Id]; found {
	//	if receiver, match := UsersById[user.Id]; match {
	//		m.Lock()
	//		receiverId := user.Id
	//		senderId := sender.Id
	//		senderType := sender.UserType
	//		receiverType := receiver.UserType
	//		if senderType != receiverType {
	//			// updation of userReg map
	//			sender.RequestUsers[receiverId] = 1
	//			secretC := receiver.SecretCode
	//			UsersByCode[secretC].PendingUsers[senderId] = 1
	//			return &pb.Success{Name: "Request sent Successfully."}, nil
	//		}
	//		m.Unlock()
	//		return &pb.Success{Name: "Invalid Request."}, nil
	//	}
	//	return &pb.Success{Name: "Invalid Request."}, nil
	//}
	return &pb.Success{Name: "Invalid Request."}, nil
}

func (s *Server) AcceptRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	//user := in
	//if sender, found := UsersByCode[user.Id]; found {
	//	if senderD, match := sender.PendingUsers[user.Id]; match {
	//		if senderD == 1 {
	//			delete(sender.PendingUsers, user.Id)
	//			sender.ConnectUsers[user.Id] = 1
	//		}
	//	} else {
	//		return &pb.Success{Name: "Sender Request not found."}, nil
	//	}
	//	receiver := UsersById[user.Id]
	//	senderId := UsersByCode[user.Id].Id
	//	if receiverD, match := receiver.RequestUsers[senderId]; match {
	//		if receiverD == 1 {
	//			delete(receiver.RequestUsers, senderId)
	//			receiver.ConnectUsers[senderId] = 1
	//		}
	//	} else {
	//		return &pb.Success{Name: "Receiver Request not found."}, nil
	//	}
	//} else {
	//	return &pb.Success{Name: "Invalid Request."}, nil
	//}
	return &pb.Success{Name: "Request Accepted."}, nil
}

func (s *Server) CancelRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	//user := in
	//if sender, found := UsersByCode[user.Id]; found {
	//	if senderD, match := sender.PendingUsers[user.Id]; match {
	//		if senderD == 1 {
	//			delete(sender.PendingUsers, user.Id)
	//		}
	//	} else {
	//		return &pb.Success{Name: "Sender Request not found."}, nil
	//	}
	//	receiver := UsersById[user.Id]
	//	senderId := UsersByCode[user.Id].Id
	//	if receiverD, match := receiver.RequestUsers[senderId]; match {
	//		if receiverD == 1 {
	//			delete(receiver.RequestUsers, senderId)
	//		}
	//	} else {
	//		return &pb.Success{Name: "Receiver Request not found."}, nil
	//	}
	//} else {
	//	return &pb.Success{Name: "Invalid Request."}, nil
	//}
	return &pb.Success{Name: "Request Rejected."}, nil
}

func (s *Server) CancelConnection(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	//user := in
	//if sender, found := UsersByCode[user.Id]; found {
	//	if senderD, match := sender.ConnectUsers[user.Id]; match {
	//		if senderD == 1 {
	//			delete(sender.ConnectUsers, user.Id)
	//		}
	//	} else {
	//		return &pb.Success{Name: "Sender Request not found."}, nil
	//	}
	//	receiver := UsersById[user.Id]
	//	senderId := UsersByCode[user.Id].Id
	//	if receiverD, match := receiver.ConnectUsers[senderId]; match {
	//		if receiverD == 1 {
	//			delete(receiver.ConnectUsers, senderId)
	//		}
	//	} else {
	//		return &pb.Success{Name: "Receiver Request not found."}, nil
	//	}
	//} else {
	//	return &pb.Success{Name: "Invalid Request."}, nil
	//}
	return &pb.Success{Name: "Connection Removed."}, nil
}

//func main() {
//	UsersByCode = make(map[string]*pb.UserDetails)
//	UsersById = make(map[string]*pb.UserDetails)
//	Donors = make(map[string]*pb.UserDetails)
//	Patients = make(map[string]*pb.UserDetails)
//	lis, err := net.Listen("tcp", ":3000")
//	if err != nil {
//		log.Println("Failed to listen server")
//	}
//	ser := grpc.NewServer()
//	pb.RegisterUserServiceServer(ser, &Server{})
//
//	if err := ser.Serve(lis); err != nil {
//		log.Println("failed to serve")
//	}
//}

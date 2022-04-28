package services

import (
	common "PlazmaDonation/PlazmaDonationAPI_gRPC/Common"
	pb "PlazmaDonation/PlazmaDonationAPI_gRPC/Gen_code"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"net/http"
	"sync"
	//grpc "grpc-go"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"strconv"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

var UsersById map[string]*pb.UserDetails
var UsersByCode map[string]*pb.UserDetails
var Donors map[string]*pb.UserDetails
var Patients map[string]*pb.UserDetails
var patient = "Patient"
var donor = "Donor"

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

func (s *Server) DeleteUser(_ context.Context, in *pb.UserDetails) (*pb.Success, error) {
	user := in
	if userStruct, found := UsersByCode[user.SecretCode]; found {
		tempId := userStruct.Id
		tempUsertype := userStruct.UserType
		if tempUsertype == donor {
			if _, found := Donors[tempId]; found {
				delete(Donors, tempId)
				//log.Println("Donor is Deleted whose name is ", tName)
				return &pb.Success{Name: "Donor is successfully deleted."}, nil
			} else {
				return &pb.Success{Name: "Donor is not present."}, nil
			}
		} else if tempUsertype == patient {
			if _, found := Patients[tempId]; found {
				delete(Patients, tempId)
				return &pb.Success{Name: "Patient is successfully deleted."}, nil
			} else {
				return &pb.Success{Name: "Patient is not present."}, nil
			}
		}
		delete(UsersByCode, user.SecretCode)
		delete(UsersById, tempId)
	}
	return nil, nil
}

func (s *Server) GetUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := in
	if userD, match := UsersByCode[user.SecretCode]; match {
		if userStruct, found := UsersById[user.Id]; found {
			checkConnection1 := false
			if _, found := userD.ConnectUsers[user.Id]; found {
				checkConnection1 = true
			}
			tempUsertype := userD.UserType
			userDetails := pb.UserResponse{}
			if userStruct.UserType != tempUsertype {
				userDetails.Id = userStruct.Id
				userDetails.Name = userStruct.Name
				if checkConnection1 {
					userDetails.Address = userStruct.Address
					userDetails.PhoneNum = userStruct.PhoneNum
				}
				return &userDetails, nil
			}
		}
	}
	return nil, nil
}

func (s *Server) UpdateUser(_ context.Context, in *pb.UserDetails) (*pb.UserDetails, error) {
	user := in

	if userStruct, found := UsersByCode[user.SecretCode]; found {
		m.Lock()
		userStruct.PhoneNum = user.PhoneNum
		userStruct.Address = user.Address
		m.Unlock()
		return &pb.UserDetails{Name: userStruct.Name, Address: userStruct.Address, PhoneNum: userStruct.PhoneNum}, nil
	}
	return nil, nil
}

func (s *Server) GetAllDonors(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	user := in
	if userStruct, found := UsersByCode[user.SecretCode]; found {
		var userData []*pb.UserResponse
		tempUsertype := userStruct.UserType
		if tempUsertype != donor {
			for key, value := range Donors {
				userDetails := pb.UserResponse{}
				userDetails.Id = value.Id
				userDetails.Name = value.Name
				if _, found := userStruct.ConnectUsers[key]; found {
					userDetails.Address = value.Address
					userDetails.PhoneNum = value.PhoneNum
				}
				userData = append(userData, &userDetails)
			}
			return &pb.ListUser{Users: userData}, nil
		}
	}
	return nil, nil
}

func (s *Server) GetAllPatients(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	user := in
	if userStruct, found := UsersByCode[user.SecretCode]; found {
		var userData []*pb.UserResponse
		tempUsertype := userStruct.UserType
		if tempUsertype != patient {
			for key, value := range Patients {
				userDetails := pb.UserResponse{}
				userDetails.Id = value.Id
				userDetails.Name = value.Name
				if _, found := userStruct.ConnectUsers[key]; found {
					userDetails.Address = value.Address
					userDetails.PhoneNum = value.PhoneNum
				}
				userData = append(userData, &userDetails)
			}
			return &pb.ListUser{Users: userData}, nil
		}
	}
	return nil, nil
}

func (s *Server) SendRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	user := in
	if sender, found := UsersByCode[user.SecretCode]; found {
		if receiver, match := UsersById[user.Id]; match {
			m.Lock()
			receiverId := user.Id
			senderId := sender.Id
			senderType := sender.UserType
			receiverType := receiver.UserType
			if senderType != receiverType {
				// updation of userReg map
				sender.RequestUsers[receiverId] = 1
				secretC := receiver.SecretCode
				UsersByCode[secretC].PendingUsers[senderId] = 1
				return &pb.Success{Name: "Request sent Successfully."}, nil
			}
			m.Unlock()
			return &pb.Success{Name: "Invalid Request."}, nil
		}
		return &pb.Success{Name: "Invalid Request."}, nil
	}
	return &pb.Success{Name: "Invalid Request."}, nil
}

func (s *Server) AcceptRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	user := in
	if sender, found := UsersByCode[user.SecretCode]; found {
		if senderD, match := sender.PendingUsers[user.Id]; match {
			if senderD == 1 {
				delete(sender.PendingUsers, user.Id)
				sender.ConnectUsers[user.Id] = 1
			}
		} else {
			return &pb.Success{Name: "Sender Request not found."}, nil
		}
		receiver := UsersById[user.Id]
		senderId := UsersByCode[user.SecretCode].Id
		if receiverD, match := receiver.RequestUsers[senderId]; match {
			if receiverD == 1 {
				delete(receiver.RequestUsers, senderId)
				receiver.ConnectUsers[senderId] = 1
			}
		} else {
			return &pb.Success{Name: "Receiver Request not found."}, nil
		}
	} else {
		return &pb.Success{Name: "Invalid Request."}, nil
	}
	return &pb.Success{Name: "Request Accepted."}, nil
}

func (s *Server) CancelRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	user := in
	if sender, found := UsersByCode[user.SecretCode]; found {
		if senderD, match := sender.PendingUsers[user.Id]; match {
			if senderD == 1 {
				delete(sender.PendingUsers, user.Id)
			}
		} else {
			return &pb.Success{Name: "Sender Request not found."}, nil
		}
		receiver := UsersById[user.Id]
		senderId := UsersByCode[user.SecretCode].Id
		if receiverD, match := receiver.RequestUsers[senderId]; match {
			if receiverD == 1 {
				delete(receiver.RequestUsers, senderId)
			}
		} else {
			return &pb.Success{Name: "Receiver Request not found."}, nil
		}
	} else {
		return &pb.Success{Name: "Invalid Request."}, nil
	}
	return &pb.Success{Name: "Request Rejected."}, nil
}

func (s *Server) CancelConnection(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
	user := in
	if sender, found := UsersByCode[user.SecretCode]; found {
		if senderD, match := sender.ConnectUsers[user.Id]; match {
			if senderD == 1 {
				delete(sender.ConnectUsers, user.Id)
			}
		} else {
			return &pb.Success{Name: "Sender Request not found."}, nil
		}
		receiver := UsersById[user.Id]
		senderId := UsersByCode[user.SecretCode].Id
		if receiverD, match := receiver.ConnectUsers[senderId]; match {
			if receiverD == 1 {
				delete(receiver.ConnectUsers, senderId)
			}
		} else {
			return &pb.Success{Name: "Receiver Request not found."}, nil
		}
	} else {
		return &pb.Success{Name: "Invalid Request."}, nil
	}
	return &pb.Success{Name: "Connection Removed."}, nil
}

func main() {
	UsersByCode = make(map[string]*pb.UserDetails)
	UsersById = make(map[string]*pb.UserDetails)
	Donors = make(map[string]*pb.UserDetails)
	Patients = make(map[string]*pb.UserDetails)
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Println("Failed to listen server")
	}
	ser := grpc.NewServer()
	pb.RegisterUserServiceServer(ser, &Server{})

	if err := ser.Serve(lis); err != nil {
		log.Println("failed to serve")
	}
}

package services

import (
	common "PlazmaDonation/Common"
	pb "PlazmaDonation/Gen_code"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"google.golang.org/api/iterator"
	"log"
	"math/rand"
	"net/http"
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

func (s *Server) GetAllDonors(ctx context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.Id)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType != pb.Type_patient {
		return nil, errors.New(common.DonorErrorMsg)
	}
	var usersArr []*pb.UserResponse
	iter := fireClient.Collection(common.CollectionUsers).Where(userTypeField, "==", pb.Type_donor).Documents(ctx)
	for {
		tempDoc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New(common.FirebaseErrorMsg)
		}
		userData, err := common.GetUserDocument(fireClient, tempDoc.Ref.ID)
		if err != nil {
			return nil, errors.New(common.ErrorGettingUserDoc)
		}
		userDetails := &pb.UserResponse{}
		userDetails.Id = userData.Id
		userDetails.Name = userData.Name
		for _, id := range userData.ConnectUsers {
			if id == userDoc.Id {
				userDetails.Address = userData.Address
				userDetails.PhoneNum = userData.PhoneNum
			}
		}
		usersArr = append(usersArr, userDetails)
	}
	return &pb.ListUser{Users: usersArr}, nil
}

func (s *Server) GetAllPatients(ctx context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.Id)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType != pb.Type_donor {
		return nil, errors.New(common.PatientErrorMsg)
	}
	var usersArr []*pb.UserResponse
	iter := fireClient.Collection(common.CollectionUsers).Where(userTypeField, "==", pb.Type_patient).Documents(ctx)
	for {
		tempDoc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New(common.FirebaseErrorMsg)
		}
		userData, err := common.GetUserDocument(fireClient, tempDoc.Ref.ID)
		if err != nil {
			return nil, errors.New(common.ErrorGettingUserDoc)
		}
		userDetails := &pb.UserResponse{}
		userDetails.Id = userData.Id
		userDetails.Name = userData.Name
		for _, id := range userData.ConnectUsers {
			if id == userDoc.Id {
				userDetails.Address = userData.Address
				userDetails.PhoneNum = userData.PhoneNum
			}
		}
		usersArr = append(usersArr, userDetails)
	}
	return &pb.ListUser{Users: usersArr}, nil
}

func (s *Server) SendRequest(ctx context.Context, in *pb.UserRequest) (*pb.Success, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	requestUserDoc, err := common.GetUserDocument(fireClient, in.RequestedUserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType == requestUserDoc.UserType {
		return nil, errors.New(common.UnAuthErrorMsg)
	}
	requestUsers := userDoc.RequestUsers
	requestUsers = append(requestUsers, requestUserDoc.Id)
	pendingUsers := requestUserDoc.PendingUsers
	pendingUsers = append(pendingUsers, userDoc.Id)
	_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  requestUsersField,
			Value: requestUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  pendingUsersField,
			Value: pendingUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	return &pb.Success{Name: "Request sent Successfully."}, nil
}

func (s *Server) AcceptRequest(ctx context.Context, in *pb.UserRequest) (*pb.Success, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	requestUserDoc, err := common.GetUserDocument(fireClient, in.RequestedUserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType == requestUserDoc.UserType {
		return nil, errors.New(common.UnAuthErrorMsg)
	}
	pendingUsers := userDoc.PendingUsers
	found := false
	for index, id := range userDoc.PendingUsers {
		if id == requestUserDoc.Id {
			found = true
			if index == 0 {
				pendingUsers = pendingUsers[index+1:]
			} else if index+1 != len(pendingUsers) {
				pendingUsers = append(pendingUsers[0:index], pendingUsers[index+1:]...)
			} else {
				pendingUsers = pendingUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	requestUsers := requestUserDoc.RequestUsers
	found = false
	for index, id := range requestUsers {
		if id == userDoc.Id {
			found = true
			if index == 0 {
				requestUsers = requestUsers[index+1:]
			} else if index+1 != len(pendingUsers) {
				requestUsers = append(requestUsers[0:index], requestUsers[index+1:]...)
			} else {
				requestUsers = requestUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	connectedUser := userDoc.ConnectUsers
	requestConnectedUsers := requestUserDoc.ConnectUsers
	connectedUser = append(connectedUser, requestUserDoc.Id)
	requestConnectedUsers = append(requestConnectedUsers, userDoc.Id)
	_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  pendingUsersField,
			Value: pendingUsers,
		},
		{
			Path:  connectUsersField,
			Value: connectedUser,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  requestUsersField,
			Value: requestUsers,
		},
		{
			Path:  connectUsersField,
			Value: requestConnectedUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	return &pb.Success{Name: "Request Accepted."}, nil
}

func (s *Server) CancelRequest(ctx context.Context, in *pb.UserRequest) (*pb.Success, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	requestUserDoc, err := common.GetUserDocument(fireClient, in.RequestedUserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType == requestUserDoc.UserType {
		return nil, errors.New(common.UnAuthErrorMsg)
	}
	pendingUsers := userDoc.PendingUsers
	found := false
	for index, id := range userDoc.PendingUsers {
		if id == requestUserDoc.Id {
			found = true
			if index == 0 {
				pendingUsers = pendingUsers[index+1:]
			} else if index+1 != len(pendingUsers) {
				pendingUsers = append(pendingUsers[0:index], pendingUsers[index+1:]...)
			} else {
				pendingUsers = pendingUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	requestUsers := requestUserDoc.RequestUsers
	found = false
	for index, id := range requestUsers {
		if id == userDoc.Id {
			found = true
			if index == 0 {
				requestUsers = requestUsers[index+1:]
			} else if index+1 != len(pendingUsers) {
				requestUsers = append(requestUsers[0:index], requestUsers[index+1:]...)
			} else {
				requestUsers = requestUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  pendingUsersField,
			Value: pendingUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  requestUsersField,
			Value: requestUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	return &pb.Success{Name: "Request Deleted."}, nil
}

func (s *Server) CancelConnection(ctx context.Context, in *pb.UserRequest) (*pb.Success, error) {
	_, fireClient, err := common.GetFireAuthFireClient(ctx)
	if err != nil {
		return nil, errors.New(common.InternalErrorMsg)
	}
	defer common.HandleFirebaseClientError(fireClient)
	userDoc, err := common.GetUserDocument(fireClient, in.UserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	requestUserDoc, err := common.GetUserDocument(fireClient, in.RequestedUserId)
	if err != nil {
		return nil, errors.New(common.ErrorGettingUserDoc)
	}
	if userDoc.UserType == requestUserDoc.UserType {
		return nil, errors.New(common.UnAuthErrorMsg)
	}
	connectedUsers := userDoc.ConnectUsers
	found := false
	for index, id := range userDoc.ConnectUsers {
		if id == requestUserDoc.Id {
			found = true
			if index == 0 {
				connectedUsers = connectedUsers[index+1:]
			} else if index+1 != len(connectedUsers) {
				connectedUsers = append(connectedUsers[0:index], connectedUsers[index+1:]...)
			} else {
				connectedUsers = connectedUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	requestConnectedUsers := requestUserDoc.ConnectUsers
	found = false
	for index, id := range requestConnectedUsers {
		if id == userDoc.Id {
			found = true
			if index == 0 {
				requestConnectedUsers = requestConnectedUsers[index+1:]
			} else if index+1 != len(requestConnectedUsers) {
				requestConnectedUsers = append(requestConnectedUsers[0:index], requestConnectedUsers[index+1:]...)
			} else {
				requestConnectedUsers = requestConnectedUsers[0:index]
			}
			break
		}
	}
	if !found {
		return nil, errors.New(common.RequestNotFound)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  connectUsersField,
			Value: connectedUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
		{
			Path:  connectUsersField,
			Value: requestConnectedUsers,
		},
	})
	if err != nil {
		return nil, errors.New(common.UpdateErrorMsg)
	}
	return &pb.Success{Name: "Connection Removed."}, nil
}

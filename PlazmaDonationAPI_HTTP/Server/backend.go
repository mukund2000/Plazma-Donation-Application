package Server

import (
	common "PlazmaDonationHTTP/Common"
	pb "PlazmaDonationHTTP/GeneratedCode"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"google.golang.org/api/iterator"
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

func Login(w http.ResponseWriter, r *http.Request) {
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
		loginResponse := &pb.Success{Name: "Login Success"}
		if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		deleteUserRequest := &pb.DeleteUserRequest{}
		if err := data.Decode(deleteUserRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		_, err = fireClient.Collection(common.CollectionUsers).Doc(deleteUserRequest.UserId).Delete(ctx)
		if err != nil {
			log.Println(common.FirebaseErrorMsg)
			return
		}
		_, err = common.GetUserDocument(fireClient, deleteUserRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		response := &pb.Success{Name: "User deleted"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserRequest{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		requestUSerDoc, err := common.GetUserDocument(fireClient, userRequest.RequestedUserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		connectedUsers := userDoc.ConnectUsers
		found := false
		for _, user := range connectedUsers {
			if user == userRequest.RequestedUserId {
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
		if err := json.NewEncoder(w).Encode(&userDetails); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}
func GetAllDonors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserDetails{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.Id)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType != pb.Type_patient {
			log.Println(common.DonorErrorMsg)
			return
		}
		var usersArr []*pb.UserResponse
		iter := fireClient.Collection(common.CollectionUsers).Where(userTypeField, "==", pb.Type_donor).Documents(ctx)
		for {
			tempDoc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Println(common.FirebaseErrorMsg)
				return
			}
			userData, err := common.GetUserDocument(fireClient, tempDoc.Ref.ID)
			if err != nil {
				log.Println(common.ErrorGettingUserDoc)
				return
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
		response := &pb.ListUser{Users: usersArr}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserDetails{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.Id)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType != pb.Type_donor {
			log.Println(common.PatientErrorMsg)
			return
		}
		var usersArr []*pb.UserResponse
		iter := fireClient.Collection(common.CollectionUsers).Where(userTypeField, "==", pb.Type_patient).Documents(ctx)
		for {
			tempDoc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Println(common.FirebaseErrorMsg)
				return
			}
			userData, err := common.GetUserDocument(fireClient, tempDoc.Ref.ID)
			if err != nil {
				log.Println(common.ErrorGettingUserDoc)
				return
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
		response := &pb.ListUser{Users: usersArr}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func UpdateContactDetails(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserDetails{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.Id)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		updatedPhoneNum := userRequest.PhoneNum
		updatedAddress := userRequest.Address
		updatedDiseaseDes := userRequest.DiseaseDes
		if updatedPhoneNum == common.EmptyString || updatedPhoneNum == userDoc.PhoneNum {
			updatedPhoneNum = userDoc.PhoneNum
		}
		if updatedAddress == common.EmptyString || updatedAddress == userDoc.Address {
			updatedAddress = userDoc.Address
		}
		if updatedDiseaseDes == common.EmptyString || updatedDiseaseDes == userDoc.DiseaseDes {
			updatedDiseaseDes = userDoc.DiseaseDes
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(userRequest.Id).Update(ctx, []firestore.Update{
			{
				Path:  phoneNumField,
				Value: userRequest.PhoneNum,
			},
			{
				Path:  addressField,
				Value: userRequest.Address,
			},
			{
				Path:  diseaseDesField,
				Value: userRequest.DiseaseDes,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		updatedUserDoc, err := common.GetUserDocument(fireClient, userRequest.Id)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if err := json.NewEncoder(w).Encode(&updatedUserDoc); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}
func SendRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserRequest{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		requestUserDoc, err := common.GetUserDocument(fireClient, userRequest.RequestedUserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType == requestUserDoc.UserType {
			log.Println(common.UnAuthErrorMsg)
			return
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
			log.Println(common.UpdateErrorMsg)
			return
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
			{
				Path:  pendingUsersField,
				Value: pendingUsers,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		response := &pb.Success{Name: "Request sent Successfully."}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserRequest{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		requestUserDoc, err := common.GetUserDocument(fireClient, userRequest.RequestedUserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType == requestUserDoc.UserType {
			log.Println(common.UnAuthErrorMsg)
			return
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
			log.Println(common.RequestNotFound)
			return
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
			log.Println(common.RequestNotFound)
			return
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
			log.Println(common.UpdateErrorMsg)
			return
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
			log.Println(common.UpdateErrorMsg)
			return
		}
		response := &pb.Success{Name: "Request Accepted."}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func CancelRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserRequest{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		requestUserDoc, err := common.GetUserDocument(fireClient, userRequest.RequestedUserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType == requestUserDoc.UserType {
			log.Println(common.UnAuthErrorMsg)
			return
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
			log.Println(common.RequestNotFound)
			return
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
			log.Println(common.RequestNotFound)
			return
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
			{
				Path:  pendingUsersField,
				Value: pendingUsers,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
			{
				Path:  requestUsersField,
				Value: requestUsers,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		response := &pb.Success{Name: "Request Deleted."}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

func CancelConnection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := json.NewDecoder(r.Body)
		userRequest := &pb.UserRequest{}
		if err := data.Decode(userRequest); err != nil {
			log.Println(common.DecodeErrorMsg)
			return
		}
		ctx := context.Background()
		_, fireClient, err := common.GetFireAuthFireClient(ctx)
		if err != nil {
			log.Println(common.InternalErrorMsg)
			return
		}
		defer common.HandleFirebaseClientError(fireClient)
		userDoc, err := common.GetUserDocument(fireClient, userRequest.UserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		requestUserDoc, err := common.GetUserDocument(fireClient, userRequest.RequestedUserId)
		if err != nil {
			log.Println(common.ErrorGettingUserDoc)
			return
		}
		if userDoc.UserType == requestUserDoc.UserType {
			log.Println(common.UnAuthErrorMsg)
			return
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
			log.Println(common.RequestNotFound)
			return
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
			log.Println(common.RequestNotFound)
			return
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(userDoc.Id).Update(ctx, []firestore.Update{
			{
				Path:  connectUsersField,
				Value: connectedUsers,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		_, err = fireClient.Collection(common.CollectionUsers).Doc(requestUserDoc.Id).Update(ctx, []firestore.Update{
			{
				Path:  connectUsersField,
				Value: requestConnectedUsers,
			},
		})
		if err != nil {
			log.Println(common.UpdateErrorMsg)
			return
		}
		response := &pb.Success{Name: "Connection Removed."}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(common.EncodeErrorMsg)
			return
		}
	default:
		http.Error(w, common.InternalErrorMsg, http.StatusBadRequest)
	}
}

//func main() {
//	http.HandleFunc("/createUser", CreateUser)
//	http.HandleFunc("/login", login)
//	http.HandleFunc("/getUser", getUser)
//	http.HandleFunc("/getDonors", getAllDonors)
//	http.HandleFunc("/getPatients", getAllPatients)
//	http.HandleFunc("/deleteUser", deleteUser)
//	http.HandleFunc("/updateContactDetails", updateContactDetails)
//	http.HandleFunc("/sendRequest", sendRequest)
//	http.HandleFunc("/acceptRequest", acceptRequest)
//	http.HandleFunc("/cancelRequest", cancelRequest)
//	http.HandleFunc("/cancelConnection", cancelConnection)
//	log.Println(http.ListenAndServe("192.168.29.118:8080", nil))
//}

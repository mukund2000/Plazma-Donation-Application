package main

import (
	"context"
	"sync"

	//grpc "grpc-go"
	"google.golang.org/grpc"
	pb "grpc/Grpc/Gen_code"
	"log"
	"math/rand"
	"net"
	"strconv"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

var UsersById map[string]*pb.UserDetails
var UsersByCode map[string]*pb.UserDetails
var Donors map[string]*pb.UserDetails
var Patients map[string]*pb.UserDetails
var patient = "Patient"
var donor = "Donor"

var m sync.Mutex

func (s *server) CreateUser(_ context.Context, in *pb.UserDetails) (*pb.UserDetails, error) {
	user := in
	m.Lock()
	user.Id = strconv.Itoa(rand.Intn(10000000))
	user.SecretCode = strconv.Itoa(rand.Intn(10000000))
	user.RequestUsers = make(map[string]int32)
	user.ConnectUsers = make(map[string]int32)
	user.PendingUsers = make(map[string]int32)
	UsersById[user.Id] = user
	UsersByCode[user.SecretCode] = user
	if user.UserType == patient {
		Patients[user.Id] = user
	} else if user.UserType == donor {
		Donors[user.Id] = user
	}
	m.Unlock()
	return user, nil
}

func (s *server) Login(_ context.Context, in *pb.UserDetails) (*pb.UserDetails, error) {
	user := in
	if userStruct, found := UsersByCode[user.SecretCode]; found {
		return userStruct, nil
	} else {
		return nil, nil
	}
}

func (s *server) DeleteUser(_ context.Context, in *pb.UserDetails) (*pb.Success, error) {
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

func (s *server) GetUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
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

func (s *server) UpdateUser(_ context.Context, in *pb.UserDetails) (*pb.UserDetails, error) {
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

func (s *server) GetAllDonors(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
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

func (s *server) GetAllPatients(_ context.Context, in *pb.UserDetails) (*pb.ListUser, error) {
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

func (s *server) SendRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
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

func (s *server) AcceptRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
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

func (s *server) CancelRequest(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
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

func (s *server) CancelConnection(_ context.Context, in *pb.UserRequest) (*pb.Success, error) {
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
	pb.RegisterUserServiceServer(ser, &server{})

	if err := ser.Serve(lis); err != nil {
		log.Println("failed to serve")
	}
}

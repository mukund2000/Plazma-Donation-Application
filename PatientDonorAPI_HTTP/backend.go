package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Id           string   `json:"id"`
	SecretCode   string   `json:"secretCode"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	PhoneNum     string   `json:"phoneNum"`
	UserType     string   `json:"usertype"`
	DiseaseDes   string   `json:"diseaseDes"`
	RequestIds   []string `json:"requestedIds"`
	PendingIds   []string `json:"pendingIds"`
	ConnectedIds []string `json:"connectedIds"`
}

type ShowUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PhoneNum string `json:"phoneNum"`
}

var usersReg = map[string]*User{}

var donors = map[string]*User{}

var patients = map[string]*User{}

var usersAcc = map[string]*User{}

var m sync.Mutex

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login is called")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("Headers")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("post is called")
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		log.Println("Encoded")
		if userStruct, found := usersReg[userData.SecretCode]; found {
			fmt.Println(userStruct)
			er := json.NewEncoder(w).Encode(userStruct)
			log.Println("User encoded")
			if er != nil {
				log.Println(er)
			}
		} else {

			userData = User{Name: "Failed"}
			json.NewEncoder(w).Encode(userData)
			log.Println("user not found")
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creatuser called")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		log.Println("Method is created")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		m.Lock()
		userData.Id = strconv.Itoa(rand.Intn(10000000))
		log.Println("id is created")
		userData.SecretCode = strconv.Itoa(rand.Intn(10000000))
		log.Println("Secret code is created")
		usersReg[userData.SecretCode] = &userData
		usersAcc[userData.Id] = &userData
		if userData.UserType == "Donor" {
			donors[userData.Id] = &userData
			log.Println("Donor is created")
		} else if userData.UserType == "Patient" {
			patients[userData.Id] = &userData
			log.Println("Patient is created")
		} else {
			log.Println("undefined user type!!")
		}
		m.Unlock()
		log.Println(usersAcc[userData.Id])
		er := json.NewEncoder(w).Encode(usersAcc[userData.Id])
		log.Println("Encoded")
		if er != nil {
			log.Println(er)
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		m.Lock()
		tempData := User{Name: "User data Deleted"}
		if userStruct, found := usersReg[userData.SecretCode]; found {
			temp_id := userStruct.Id
			temp_userType := userStruct.UserType
			if temp_userType == "Donor" {
				if donorD, found := donors[temp_id]; found {
					t_name := donorD.Name
					delete(donors, temp_id)
					log.Println("Donor is Deleted whose name is ", t_name)
				} else {
					log.Println("Donor is not Present.")
				}
			} else if temp_userType == "Patient" {
				if patientD, found := patients[temp_id]; found {
					t_name := patientD.Name
					delete(patients, temp_id)
					log.Println("Patient is Deleted whose name is ", t_name)
				} else {
					log.Println("Patient is not Present.")
				}
			}
			delete(usersReg, userData.SecretCode)
			delete(usersAcc, temp_id)
			json.NewEncoder(w).Encode(tempData)
		} else {
			tempData = User{Name: "User not found"}
			json.NewEncoder(w).Encode(tempData)
		}
		m.Unlock()
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get user envoked")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("headers envoked")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("Method envoked")
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		log.Println("user unmarshal")
		if userD, match := usersReg[userData.SecretCode]; match {
			if userStruct, found := usersAcc[userData.Id]; found {
				log.Println("Get user found")
				checkConnection1 := false
				for _, ids := range userD.ConnectedIds {
					if ids == userData.Id {
						checkConnection1 = true
					}
				}
				tempUsertype := userD.UserType
				userDetails := ShowUser{}
				if userStruct.UserType != tempUsertype {
					log.Println("Get user success")
					userDetails.Id = userStruct.Id
					userDetails.Name = userStruct.Name
					if checkConnection1 {
						userDetails.Address = userStruct.Address
						userDetails.PhoneNum = userStruct.PhoneNum
					}
					log.Println("USerName: ", userStruct.Name, "USer Phone Num: ", userStruct.PhoneNum, "User Address: ", userStruct.Address)
					er := json.NewEncoder(w).Encode(userDetails)
					if er != nil {
						log.Println(er)
					}
					log.Println("Get user encoded")
				} else {
					userDetails = ShowUser{Name: "Same Type user error."}
					json.NewEncoder(w).Encode(userDetails)
					log.Println("Same Type user error.")
				}
			}
		} else {
			tempData := ShowUser{Name: "User not found error."}
			json.NewEncoder(w).Encode(tempData)
			log.Println("User not found error.")
		}
	}

}

func getAllDonors(w http.ResponseWriter, r *http.Request) {
	log.Println("donors envoked")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		log.Println("post envoked")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", PhoneNum: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		log.Println("body envoked")
		if userStruct, found := usersReg[userData.SecretCode]; found {
			log.Println("donors found")
			var userData = []ShowUser{}
			temp_userType := userStruct.UserType
			donor := "Donor"
			if temp_userType != donor {
				log.Println("donors allowed")
				for key, value := range donors {
					userDetails := ShowUser{}
					userDetails.Id = value.Id
					userDetails.Name = value.Name
					for _, ids := range userStruct.ConnectedIds {
						if key == ids {
							userDetails.Address = value.Address
							userDetails.PhoneNum = value.PhoneNum
						}
					}
					log.Println("Donor id: ", key)
					log.Println("Donor Name: ", value.Name, "Phone Num: ", value.PhoneNum)
					userData = append(userData, userDetails)
					log.Println("donor created")
				}
				er := json.NewEncoder(w).Encode(userData)
				if er != nil {
					log.Println(er)
				}
				log.Println("donors encoded")
			} else {
				log.Println("Only Patients can access Donors data, Sorry!!")
			}
		}
	}
}

func getAllPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		if userStruct, found := usersReg[userData.SecretCode]; found {
			var userData = []ShowUser{}
			temp_userType := userStruct.UserType
			patient := "Patient"
			if temp_userType != patient {
				for key, value := range patients {
					userDetails := ShowUser{}
					userDetails.Id = value.Id
					userDetails.Name = value.Name
					for _, ids := range userStruct.ConnectedIds {
						if key == ids {
							userDetails.Address = value.Address
							userDetails.PhoneNum = value.PhoneNum
						}
					}
					log.Println("Patient id: ", key)
					log.Println("Patient Name: ", value.Name, "Phone Num: ", value.PhoneNum)
					userData = append(userData, userDetails)
				}
				er := json.NewEncoder(w).Encode(userData)
				if er != nil {
					log.Println(er)
				}
			} else {
				log.Println("Only Donors can access Patients data, Sorry!!")
			}
		}
	}
}

func updateContactDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		m.Lock()
		if userStruct, found := usersReg[userData.SecretCode]; found {
			temp_id := userStruct.Id
			temp_userType := userStruct.UserType
			userStruct.PhoneNum = userData.PhoneNum
			userStruct.Address = userData.Address
			log.Println("Registered User data has been updated")
			if userA, match := usersAcc[temp_id]; match {
				userA.Address = userData.Address
				userA.PhoneNum = userData.PhoneNum
				log.Println("Registered Accessed User data has been updated")
			}
			if temp_userType == "Donor" {
				if donorD, match := donors[temp_id]; match {
					donorD.Address = userData.Address
					donorD.PhoneNum = userData.PhoneNum
					log.Println("Registered Donor data has been updated")
				}
			} else if temp_userType == "Patient" {
				if patientD, match := patients[temp_id]; match {
					patientD.Address = userData.Address
					patientD.PhoneNum = userData.PhoneNum
					log.Println("Registered Patient data has been updated")
				}
			}
			log.Println("Contact Details have been updated")
			er := json.NewEncoder(w).Encode(userStruct)
			if er != nil {
				log.Println(er)
			}
		}
		m.Unlock()
	}
}

func sendRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("request envoked")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("headers envoked")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("method envoked")
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		log.Println("request unmarshal")
		if sender, found := usersReg[userData.SecretCode]; found {
			if receiver, match := usersAcc[userData.Id]; match {
				log.Println("request find")
				receiverId := userData.Id
				senderId := sender.Id
				senderType := sender.UserType
				receiverType := receiver.UserType
				if senderType != receiverType {
					// updation of userReg map
					m.Lock()
					sender.RequestIds = append(sender.RequestIds, receiverId)
					secretC := receiver.SecretCode
					usersReg[secretC].PendingIds = append(usersReg[secretC].PendingIds, senderId)
					m.Unlock()
					log.Println("Request Id:", receiverId, " of sender is made.")
					log.Println("request success")
				}
				log.Println("RequestIds of sender:", sender.RequestIds)
				log.Println("Pending Request User Ids of receiver:", receiver.PendingIds)
				er := json.NewEncoder(w).Encode(sender)
				if er != nil {
					log.Println(er)
				}
				log.Println("request encoded")
			}

		}

	}
}

func acceptRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		if sender, found := usersReg[userData.SecretCode]; found {
			m.Lock()
			for index, ids := range sender.PendingIds {
				if userData.Id == ids {
					log.Println("Request Id:", ids, "of sender is accepted.")
					sender.PendingIds = append(sender.PendingIds[:index], sender.PendingIds[index+1:]...)
					sender.ConnectedIds = append(sender.ConnectedIds, ids)
				}
			}
			receiver := usersAcc[userData.Id]
			senderId := usersReg[userData.SecretCode].Id
			for index2, reqIds := range receiver.RequestIds {
				if senderId == reqIds {
					log.Println("Request Id:", reqIds, " of receiver is accepted.")
					receiver.RequestIds = append(receiver.RequestIds[:index2], receiver.RequestIds[index2+1:]...)
					receiver.ConnectedIds = append(receiver.ConnectedIds, reqIds)
				}
			}
			m.Unlock()
			log.Println("sender Pending Ids:", sender.PendingIds)
			log.Println("sender Connected Ids:", sender.ConnectedIds)
			log.Println("Receiver Pending Ids:", receiver.PendingIds)
			log.Println("Receiver Connected Ids:", receiver.ConnectedIds)
			er := json.NewEncoder(w).Encode(sender)
			if er != nil {
				log.Println(er)
			}
		}
	}
}

func cancelRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		if sender, found := usersReg[userData.SecretCode]; found {
			m.Lock()
			for index, ids := range sender.PendingIds {
				if userData.Id == ids {
					log.Println("Request Id:", ids, "of sender is rejected.")
					sender.PendingIds = append(sender.PendingIds[:index], sender.PendingIds[index+1:]...)
				}
			}
			receiver := usersAcc[userData.Id]
			senderId := usersReg[userData.SecretCode].Id
			for index2, reqIds := range receiver.RequestIds {
				if senderId == reqIds {
					log.Println("Request Id:", reqIds, " of receiver is rejected.")
					receiver.RequestIds = append(receiver.RequestIds[:index2], receiver.RequestIds[index2+1:]...)
				}
			}
			m.Unlock()
			log.Println("sender Pending Ids:", sender.PendingIds)
			log.Println("sender Connected Ids:", sender.ConnectedIds)
			log.Println("Receiver request Ids:", receiver.RequestIds)
			log.Println("Receiver Connected Ids:", receiver.ConnectedIds)
			er := json.NewEncoder(w).Encode(sender)
			if er != nil {
				log.Println(er)
			}
		}
	}
}

func cancelConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		userData := User{Id: "", SecretCode: "", Name: "", Address: "", UserType: "", DiseaseDes: "", RequestIds: []string{}, PendingIds: []string{}, ConnectedIds: []string{}}
		if err := json.Unmarshal(reqBody, &userData); err != nil {
			log.Println(err)
		}
		if sender, found := usersReg[userData.SecretCode]; found {
			m.Lock()
			for index, ids := range sender.ConnectedIds {
				if userData.Id == ids {
					log.Println("Connection Id:", ids, "of sender is removed.")
					sender.ConnectedIds = append(sender.ConnectedIds[:index], sender.ConnectedIds[index+1:]...)
				}
			}
			receiver := usersAcc[userData.Id]
			senderId := usersReg[userData.SecretCode].Id
			for index2, reqIds := range receiver.ConnectedIds {
				if senderId == reqIds {
					log.Println("Connection Id:", reqIds, " of receiver is removed.")
					receiver.ConnectedIds = append(receiver.ConnectedIds[:index2], receiver.ConnectedIds[index2+1:]...)
				}
			}
			m.Unlock()
			log.Println("sender Connected Ids:", sender.ConnectedIds)
			log.Println("Receiver Connected Ids:", receiver.ConnectedIds)
			er := json.NewEncoder(w).Encode(sender)
			if er != nil {
				log.Println(er)
			}
		}
	}
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

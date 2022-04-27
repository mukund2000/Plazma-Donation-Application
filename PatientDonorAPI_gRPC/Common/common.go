package Common

import (
	pb "PlazmaDonation/PatientDonorAPI_gRPC/Gen_code"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
	"log"
)

const CollectionUsers = "User"
const ApiKey = ""
const credentialPath = "C:\\go_app\\src\\MyProjects\\Plazma-Donation-Application\\plazma-donation-application-firebase-adminsdk-e9wug-b2646babdf.json"
const UnableToGetInstance = "unable To Get Instance"
const UnableToGetAuth = "unable to get auth"
const UnableToGetClient = "unable to get client"
const AuthorizationHeaderName = "authorization"
const EmptyString = ""
const AuthInvalidMsg = "error: Authorization Id invalid"
const AuthTokenNotFound = "auth Token not provided"
const TokenExpired = "token Expired"
const InternalErrorMsg = "something went wrong in Server"
const InvalidLoginErrorMsg = "invalid user email or password"
const AddErrorMsg = "unable to append data in collection"
const ErrorGettingUserDoc = "error Getting User Doc"

func GetFirebaseInstance() (*firebase.App, error) {
	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile(credentialPath)
	return firebase.NewApp(ctx, nil, serviceAccount)
}

func GetFireAuthFireClient(ctx context.Context) (*auth.Client, *firestore.Client, error) {
	fireApp, err := GetFirebaseInstance()
	if err != nil {
		return nil, nil, errors.New(UnableToGetInstance)
	}
	ctx = context.Background()
	fireAuth, err := fireApp.Auth(ctx)
	if err != nil {
		return nil, nil, errors.New(UnableToGetAuth)
	}
	fireClient, err := fireApp.Firestore(ctx)
	if err != nil {
		return nil, nil, errors.New(UnableToGetClient)
	}
	return fireAuth, fireClient, nil
}

func HandleFirebaseClientError(fireClient *firestore.Client) {
	err := fireClient.Close()
	if err != nil {
		log.Println(err)
	}
}
func GetTokenFromContext(ctx context.Context) (string, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata: No proper metadata present")
	}
	authToken := ""
	for key, value := range headers {
		if key == AuthorizationHeaderName && len(value) != 0 {
			authToken = value[0]
			break
		}
	}
	if authToken == EmptyString {
		return "", errors.New(AuthInvalidMsg)
	}

	return authToken, nil
}
func GetUserDocument(fireClient *firestore.Client, userId string) (*pb.UserDetails, error) {
	ctx := context.Background()
	foundDoc, err := fireClient.Collection(CollectionUsers).Doc(userId).Get(ctx)
	if err != nil {
		return nil, err
	}
	userDoc := &pb.UserDetails{}
	jsonBody, err := json.Marshal(foundDoc.Data())
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonBody, userDoc); err != nil {
		return nil, err
	}
	return userDoc, nil
}

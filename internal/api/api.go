package api

import (
	"cmd/main.go/internal/service/user_request"
	desc "cmd/main.go/pkg/my-api"
)

const (
	createUserLogTag     = "CreateUser"
	listUserLogTag       = "ListUser"
	removeUserLogTag     = "RemoveUser"
	updateUserByIdLogTag = "UpdateUserById"
)

type Implementation struct {
	desc.UnimplementedApiServiceServer
	userRequestService user_request.ServiceInterface
}

func NewApiService(userRequestService user_request.ServiceInterface) desc.ApiServiceServer {
	return &Implementation{userRequestService: userRequestService}
}

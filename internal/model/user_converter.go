package model

import (
	"errors"

	desc "cmd/main.go/pkg/my-api"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ErrUnableToConvertUserRequestStatus is a unable to convert model error
	ErrUnableToConvertUserRequestStatus = errors.New("unable to convert user request status")
)

// ConvertUserRequestToPb - convert UserRequest to protobuf UserRequest message
func ConvertUserRequestToPb(userRequest *UserRequest) (*desc.CreateUserRequest, error) {

	return &desc.CreateUserRequest{
		IdUser:    userRequest.ID_user,
		Name:      userRequest.Name,
		Email:     userRequest.Email,
		CreatedAt: timestamppb.New(userRequest.CreatedAt),
		UpdatedAt: timestamppb.New(userRequest.UpdatedAt.Time),
		DoneAt:    timestamppb.New(userRequest.DoneAt.Time),
		DeletedAt: timestamppb.New(userRequest.DeletedAt.Time),
	}, nil
}

// ConvertRepeatedUserRequestsToPb - convert slice of UserRequest to slice of protobuf UserRequest messages
func ConvertRepeatedUserRequestsToPb(userRequests []UserRequest) ([]*desc.User, error) {
	var userRequestsPb []*desc.User

	for i := range userRequests {
		userRequest, err := ConvertUserToPb(&userRequests[i])
		if err != nil {
			return nil, err
		}
		userRequestsPb = append(userRequestsPb, userRequest)
	}

	return userRequestsPb, nil
}

// ConvertPbToUserRequest - convert protobuf UserRequest message to UserRequest
func ConvertPbToUserRequest(userRequest *desc.CreateUserRequest) (*UserRequest, error) {

	return &UserRequest{
		ID_user:   userRequest.IdUser,
		Name:      userRequest.Name,
		Email:     userRequest.Email,
		CreatedAt: userRequest.CreatedAt.AsTime(),
		UpdatedAt: ConvertPbTimeToNullableTime(userRequest.UpdatedAt),
		DoneAt:    ConvertPbTimeToNullableTime(userRequest.DoneAt),
		DeletedAt: ConvertPbTimeToNullableTime(userRequest.DeletedAt),
	}, nil
}

func ConvertUserToPb(userRequest *UserRequest) (*desc.User, error) {

	return &desc.User{
		Id:        userRequest.ID_user,
		Name:      userRequest.Name,
		Email:     userRequest.Email,
		CreatedAt: timestamppb.New(userRequest.CreatedAt),
		UpdatedAt: timestamppb.New(userRequest.UpdatedAt.Time),
		DoneAt:    timestamppb.New(userRequest.DoneAt.Time),
		DeletedAt: timestamppb.New(userRequest.DeletedAt.Time),
	}, nil
}

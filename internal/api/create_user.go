package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"cmd/main.go/internal/model"
	desc "cmd/main.go/pkg/my-api"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {

	req.CreatedAt = timestamppb.New(time.Now())

	if err := req.Validate(); err != nil {
		log.Print(fmt.Sprintf("%s: invalid argument", createUserLogTag),
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// ch, err := i.userRequestService.CheckExistsUserRequest(ctx, req.IdUser)
	// if err != nil {
	// 	log.Print(ctx, fmt.Sprintf("%s: User already exist", createUserLogTag),
	// 		"err", err,
	// 	)
	// }

	newItem := desc.CreateUserRequest{
		IdUser:    req.GetIdUser(),
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.GetUpdatedAt(),
		DeletedAt: req.GetDeletedAt(),
		DoneAt:    req.GetDoneAt(),
	}

	User, err := model.ConvertPbToUserRequest(&newItem)

	if err != nil {
		log.Print(fmt.Sprintf("%s: unable to convert Pb message to EquipmentRequest", createUserLogTag),
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	id, err := i.userRequestService.CreateUserRequest(ctx, User)

	if err != nil {
		log.Print(fmt.Sprintf("%s: userRequestService.CreateUserRequest failed", createUserLogTag),
			"err", err,
			"iduser", req.IdUser,
			"name", req.Name,
			"email", req.Email,
			"createdAt", req.CreatedAt,
			"updatedAt", req.UpdatedAt,
			"deletedAt", req.DeletedAt,
			"doneAt", req.DoneAt,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if id == 0 {
		log.Print(fmt.Sprintf("%s: : userRequestService.CreateUserRequest failed", createUserLogTag),
			"err", "unable to get created user",
			"iduser", req.IdUser,
			"name", req.Name,
			"email", req.Email,
			"createdAt", req.CreatedAt,
			"updatedAt", req.UpdatedAt,
			"deletedAt", req.DeletedAt,
			"doneAt", req.DoneAt,
		)

		return nil, status.Error(codes.Internal, "unable to get created user")
	}

	log.Print(fmt.Sprintf("%s: success ", createUserLogTag),
		"iduser ", id,
	)

	return &desc.CreateUserResponse{
		IdUser: id,
	}, nil

}

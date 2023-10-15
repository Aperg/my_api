package api

import (
	"cmd/main.go/internal/model"
	desc "cmd/main.go/pkg/my-api"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pkg/errors"
)

func (i *Implementation) GetUserById(ctx context.Context, req *desc.GetUserByIdRequest) (*desc.GetUserByIdResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, err := i.userRequestService.GetUserByIdRequest(ctx, req.GetIdsUser())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	us := make([]*desc.User, len(users))
	for idx := range users {
		us[idx], err = model.ConvertUserToPb(&users[idx])
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	return &desc.GetUserByIdResponse{User: us}, nil
}

package api

import (
	"cmd/main.go/internal/logger"
	"cmd/main.go/internal/model"
	desc "cmd/main.go/pkg/my-api"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetUserById(ctx context.Context, req *desc.GetUserByIdRequest) (*desc.GetUserByIdResponse, error) {
	err := req.Validate()
	if err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: invalid argument", GetUserByIdLogTag),
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, err := i.userRequestService.GetUserByIdRequest(ctx, req.GetIdsUser())
	if err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: userRequestService.GetUserByIdRequest failed", GetUserByIdLogTag),
			"err", err,
			"equipmentRequestId", req.GetIdsUser(),
		)

		return nil, status.Error(codes.Internal, err.Error())

	}

	us := make([]*desc.User, len(users))
	for idx := range users {
		us[idx], err = model.ConvertUserToPb(&users[idx])
		if err != nil {
			logger.ErrorKV(ctx, fmt.Sprintf("%s: unable to convert User to Pb message", GetUserByIdLogTag),
				"err", err,
			)

			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	logger.Info(ctx, fmt.Sprintf("%s: success", GetUserByIdLogTag))
	return &desc.GetUserByIdResponse{User: us}, nil
}

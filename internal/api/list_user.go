package api

import (
	"context"
	"fmt"

	"cmd/main.go/internal/logger"
	"cmd/main.go/internal/model"
	desc "cmd/main.go/pkg/my-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListUser(ctx context.Context, req *desc.ListUserRequest) (*desc.ListUserResponse, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: invalid argument", listUserLogTag),
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userRequests, err := i.userRequestService.ListUserRequest(ctx, req.Limit, req.Offset)

	if err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: userRequestService.ListUserRequest failed", listUserLogTag),
			"err", err,
			"limit", req.Limit,
			"offset", req.Offset,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if userRequests == nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: userRequestService.ListUserRequest failed", listUserLogTag),
			"err", "unable to get list of user requests",
			"limit", req.Limit,
			"offset", req.Offset,
		)

		return nil, status.Error(codes.NotFound, "unable to get list of user requests")
	}

	userRequestPb, err := model.ConvertRepeatedUserRequestsToPb(userRequests)

	if err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: unable to convert list of UserRequests to Pb message", listUserLogTag),
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info(ctx, fmt.Sprintf("%s: success", listUserLogTag))

	return &desc.ListUserResponse{
		Items: userRequestPb,
	}, nil
}

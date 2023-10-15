package api

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	desc "cmd/main.go/pkg/my-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateUserById(ctx context.Context, req *desc.UpdateUserByIdRequest) (*desc.UpdateUserByIdResponse, error) {

	if err := req.Validate(); err != nil {
		log.Print(fmt.Sprintf("%s: invalid argument", updateUserByIdLogTag),
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := i.userRequestService.UpdateUserByIdRequest(ctx, req.GetIdUser(), req.GetName(), req.GetEmail())

	if err != nil {
		log.Print(fmt.Sprintf("%s: userRequestService.UpdateUserByIdRequest failed", updateUserByIdLogTag),
			"err", err,
			"userRequestId", req.GetIdUser(),
			"name", req.GetName(),
			"email", req.GetEmail(),
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if !result {
		log.Print(fmt.Sprintf("%s: userRequestService.UpdateUseByIdRequest failed", updateUserByIdLogTag),
			"err", "unable to update user of user request, no rows affected",
			"userRequestId", req.GetIdUser(),
			"name", req.GetName(),
			"email", req.GetEmail(),
		)

		return nil, status.Error(codes.Internal, "unable to update user of user request")
	}

	log.Print(fmt.Sprintf("%s: success", updateUserByIdLogTag))

	return &desc.UpdateUserByIdResponse{
		Updated: result,
	}, nil
}

package api

import (
	"context"
	"fmt"

	desc "cmd/main.go/pkg/my-api"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) RemoveUser(ctx context.Context, req *desc.RemoveUserRequest) (*desc.RemoveUserResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := i.userRequestService.RemoveUserRequest(ctx, req.GetIdsUser())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err != nil {
		log.Info().Msgf(fmt.Sprintf("%s: userRequestService.RemoveUserRequest failed", removeUserLogTag),
			"err", err,
			"UserRequestId", req.GetIdsUser(),
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if !result {
		log.Info().Msgf(fmt.Sprintf("%s: userRequestService.RemoveUsertRequest failed", removeUserLogTag),
			"err", "unable to remove user request, no rows affected",
			"usertRequestId", req.GetIdsUser(),
		)

		return nil, status.Error(codes.Internal, "unable to remove user request")
	}

	log.Info().Msgf(fmt.Sprintf("%s: success", removeUserLogTag))

	return &desc.RemoveUserResponse{
		Removed: result}, nil
}

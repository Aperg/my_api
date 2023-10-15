package user_request

import (
	"context"

	"cmd/main.go/internal/database"
	"cmd/main.go/internal/model"
	"cmd/main.go/internal/repo"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type service struct {
	db                *sqlx.DB
	requestRepository repo.UserRequestRepo
}

// ServiceInterface is a interface for User request service
type ServiceInterface interface {
	CreateUserRequest(ctx context.Context, userRequest *model.UserRequest) (uint64, error)
	GetUserByIdRequest(ctx context.Context, IDs []uint64) ([]model.UserRequest, error)
	ListUserRequest(ctx context.Context, limit uint64, offset uint64) ([]model.UserRequest, error)
	RemoveUserRequest(ctx context.Context, IDs []uint64) (bool, error)
	CheckExistsUserRequest(ctx context.Context, ID uint64) (bool, error)
	UpdateUserByIdRequest(ctx context.Context, userRequestID uint64, name, email string) (bool, error)
}

// New is a function to create a new service
func New(db *sqlx.DB, requestRepository repo.UserRequestRepo) ServiceInterface {
	return service{
		db:                db,
		requestRepository: requestRepository,
	}
}

// ErrNoExistsUserRequest is a "User request not founded" error
var ErrNoExistsUserRequest = errors.New("user request with this id does not exist")

// ErrNoCreatedUserRequest is a "unable to create User request" error
var ErrNoCreatedUserRequest = errors.New("unable to create user request")

// ErrNoListUserRequest is a "unable to get list of User requests" error
var ErrNoListUserRequest = errors.New("unable to get list of user requests")

// ErrNoRemovedUserRequest is a "unable to remove User request" error
var ErrNoRemovedUserRequest = errors.New("unable to remove user request")

// ErrNoUpdatedUserIDUserRequest is a "unable to update User id of User request" error
var ErrNoUpdatedUserIDUserRequest = errors.New("unable to update user of user request")

func (s service) CreateUserRequest(ctx context.Context, userRequest *model.UserRequest) (uint64, error) {

	createdRequestID, txErr := database.WithTxReturnUint64(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) (uint64, error) {
		id, err := s.requestRepository.CreateUserRequest(ctx, userRequest, tx)
		if err != nil {
			return 0, errors.Wrap(err, "requestRepository.CreateUserRequest")
		}

		if id == 0 {
			return 0, ErrNoCreatedUserRequest
		}
		userRequest.ID_user = id

		return id, nil
	})

	if txErr != nil {
		return createdRequestID, txErr
	}

	return createdRequestID, nil
}

func (s service) GetUserByIdRequest(ctx context.Context, IDs []uint64) ([]model.UserRequest, error) {

	userRequests, err := s.requestRepository.GetUserByIdRequest(ctx, IDs)
	if err != nil {
		return nil, errors.Wrap(err, "repository.GetUserByIdtRequest")
	}

	if userRequests == nil {
		return nil, ErrNoListUserRequest
	}

	return userRequests, nil
}

func (s service) CheckExistsUserRequest(ctx context.Context, ID uint64) (bool, error) {

	exists, err := s.requestRepository.Exists(ctx, ID)

	if err != nil {
		return false, errors.Wrap(err, "repository.RemoveUserRequest")
	}

	if !exists {
		return false, ErrNoExistsUserRequest
	}

	return exists, nil
}

func (s service) ListUserRequest(ctx context.Context, limit uint64, offset uint64) ([]model.UserRequest, error) {

	userRequests, err := s.requestRepository.ListUserRequest(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "repository.ListUserRequest")
	}

	if userRequests == nil {
		return nil, ErrNoListUserRequest
	}

	return userRequests, nil
}

func (s service) RemoveUserRequest(ctx context.Context, IDs []uint64) (bool, error) {

	deleted, txErr := database.WithTxReturnBool(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) (bool, error) {
		result, err := s.requestRepository.RemoveUserRequest(ctx, IDs, tx)
		if err != nil {
			return false, errors.Wrap(err, "repository.RemoveUserRequest")
		}

		if !result {
			return false, ErrNoRemovedUserRequest
		}

		return result, nil
	})

	if txErr != nil {
		return deleted, txErr
	}

	return deleted, nil
}

func (s service) UpdateUserByIdRequest(ctx context.Context, userRequestID uint64, name, email string) (bool, error) {

	updated, txErr := database.WithTxReturnBool(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) (bool, error) {
		result, err := s.requestRepository.UpdateUserByIdRequest(ctx, userRequestID, name, email, tx)
		if err != nil {
			return false, errors.Wrap(err, "repository.UpdateUserByIdRequest")
		}

		if !result {
			return false, ErrNoUpdatedUserIDUserRequest
		}

		return result, nil
	})

	if txErr != nil {
		return updated, txErr
	}

	return updated, nil
}

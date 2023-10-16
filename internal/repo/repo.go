package repo

import (
	"context"
	"database/sql"
	"time"

	"cmd/main.go/internal/database"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"cmd/main.go/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	userRequestTable             = "users"
	userRequestIDColumn          = "id_user"
	userRequestNameColumn        = "name"
	userRequestEmailColumn       = "email"
	userRequestUpdatedAtColumn   = "updated_at"
	userRequestCreatedAtColumn   = "created_at"
	userRequestDoneAtColumn      = "done_at"
	userRequestDeletedAtAtColumn = "deleted_at"
)

// EuserRequestRepo is DAO for Euser Request
type UserRequestRepo interface {
	CreateUserRequest(ctx context.Context, userRequest *model.UserRequest, tx *sqlx.Tx) (uint64, error)
	GetUserByIdRequest(ctx context.Context, IDs []uint64) ([]model.UserRequest, error)
	ListUserRequest(ctx context.Context, limit uint64, offset uint64) ([]model.UserRequest, error)
	RemoveUserRequest(ctx context.Context, IDs []uint64, tx *sqlx.Tx) (bool, error)
	Exists(ctx context.Context, userRequestID uint64) (bool, error)
	UpdateUserByIdRequest(ctx context.Context, uerRequestID uint64, name, email string, tx *sqlx.Tx) (bool, error)
}

type userRequestRepo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewEuserRequestRepo returns Repo interface
func NewUserRequestRepo(db *sqlx.DB, batchSize uint) *userRequestRepo {
	return &userRequestRepo{
		db:        db,
		batchSize: batchSize,
	}
}

func (r *userRequestRepo) CreateUserRequest(ctx context.Context, userRequest *model.UserRequest, tx *sqlx.Tx) (uint64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.CreateUserRequest")
	defer span.Finish()

	sb := database.StatementBuilder.
		Insert(userRequestTable).
		Columns(
			userRequestIDColumn,
			userRequestNameColumn,
			userRequestEmailColumn,
			userRequestCreatedAtColumn,
			userRequestUpdatedAtColumn,
			userRequestDeletedAtAtColumn,
			userRequestDoneAtColumn).
		Values(
			userRequest.ID_user,
			userRequest.Name,
			userRequest.Email,
			userRequest.CreatedAt,
			userRequest.UpdatedAt,
			userRequest.DeletedAt,
			userRequest.DoneAt,
		).Suffix("RETURNING " + userRequestIDColumn)

	query, args, err := sb.ToSql()
	if err != nil {
		return 0, err
	}

	var queryer sqlx.QueryerContext
	if tx == nil {
		queryer = r.db
	} else {
		queryer = tx
	}

	var id uint64
	err = queryer.QueryRowxContext(ctx, query, args...).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, errors.Wrap(err, "db.QueryRowxContext()")
	}

	return id, nil
}

func (r *userRequestRepo) GetUserByIdRequest(ctx context.Context, IDs []uint64) ([]model.UserRequest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.GetUserByIdRequest")
	defer span.Finish()
	sb := database.StatementBuilder.
		Select("*").
		From(userRequestTable).
		Where(sq.Eq{userRequestIDColumn: IDs})

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}
	var userRequests []model.UserRequest
	err = r.db.SelectContext(ctx, &userRequests, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return userRequests, nil
}

func (r *userRequestRepo) ListUserRequest(ctx context.Context, limit uint64, offset uint64) ([]model.UserRequest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.ListUserRequest")
	defer span.Finish()
	sb := database.StatementBuilder.
		Select("*").
		From(userRequestTable).
		Where(sq.Eq{userRequestDeletedAtAtColumn: nil}).
		OrderBy(userRequestIDColumn).
		Limit(limit).
		Offset(offset)

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var userRequests []model.UserRequest
	err = r.db.SelectContext(ctx, &userRequests, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return userRequests, nil
}

func (r *userRequestRepo) RemoveUserRequest(ctx context.Context, IDs []uint64, tx *sqlx.Tx) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.RemoveUserRequest")
	defer span.Finish()
	sb := database.StatementBuilder.
		Update(userRequestTable).
		Set(userRequestDeletedAtAtColumn, time.Now()).
		Where(sq.And{
			sq.Eq{userRequestIDColumn: IDs},
			sq.Eq{userRequestDeletedAtAtColumn: nil}})

	query, args, err := sb.ToSql()
	if err != nil {
		return false, err
	}

	var queryer sqlx.ExecerContext
	if tx == nil {
		queryer = r.db
	} else {
		queryer = tx
	}

	var result sql.Result
	result, err = queryer.ExecContext(ctx, query, args...)

	if err != nil {
		return false, errors.Wrap(err, "db.SelectContext()")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "repo.RowsAffected()")
	}

	if affected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *userRequestRepo) Exists(ctx context.Context, userRequestID uint64) (bool, error) {

	sb := database.StatementBuilder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(userRequestTable).
		Where(sq.And{
			sq.Eq{userRequestIDColumn: userRequestID},
			sq.Eq{userRequestDeletedAtAtColumn: nil}}).
		Suffix(")")

	query, args, err := sb.ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRowxContext(ctx, query, args...).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "db.QueryRowxContext()")
	}

	return exists, nil
}

// // nolint:dupl
func (r *userRequestRepo) UpdateUserByIdRequest(ctx context.Context, userRequestID uint64, name, email string, tx *sqlx.Tx) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.UpdateUserByIdRequest")
	defer span.Finish()
	sb := database.StatementBuilder.
		Update(userRequestTable).
		Set(userRequestUpdatedAtColumn, time.Now()).
		Set(userRequestNameColumn, name).
		Set(userRequestEmailColumn, email).
		Where(sq.And{
			sq.Eq{userRequestIDColumn: userRequestID},
			sq.Eq{userRequestDeletedAtAtColumn: nil}})

	query, args, err := sb.ToSql()
	if err != nil {
		return false, err
	}

	var queryer sqlx.ExecerContext
	if tx == nil {
		queryer = r.db
	} else {
		queryer = tx
	}

	var result sql.Result
	result, err = queryer.ExecContext(ctx, query, args...)

	if err != nil {
		return false, errors.Wrap(err, "db.SelectContext()")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "repo.RowsAffected()")
	}

	if affected == 0 {
		return false, nil
	}

	return true, nil
}

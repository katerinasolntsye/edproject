package repository

import (
	"context"

	"github.com/katerinasolntsye/fulleng/internal/repository/model"
)

type Repository interface {
	GetIncomingPostbacks(ctx context.Context) ([]model.IncomingPostback, error)
	GetTrackers(ctx context.Context) ([]model.Tracker, error)
	GetSendPostbacks(ctx context.Context) ([]model.SendPostback, error)
	CreateUser(ctx context.Context, email, hashedPassword, name, timezone string) error
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateUserCredentials(ctx context.Context, user *model.User) error
	GetCredentials(ctx context.Context, email string) (*model.Credentials, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
}

package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/katerinasolntsye/fulleng/internal/repository"
	"github.com/katerinasolntsye/fulleng/internal/repository/model"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetIncomingPostbacks(ctx context.Context) ([]model.IncomingPostback, error) {
	return s.repo.GetIncomingPostbacks(ctx)
}

func (s *Service) GetTrackers(ctx context.Context) ([]model.Tracker, error) {
	return s.repo.GetTrackers(ctx)
}

func (s *Service) GetSendPostbacks(ctx context.Context) ([]model.SendPostback, error) {
	return s.repo.GetSendPostbacks(ctx)
}

func (s *Service) CreateUser(ctx context.Context, user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, user.Email, string(hashedPassword), user.Name, user.Timezone)
}

func (s *Service) UpdateUser(ctx context.Context, user *model.User) error {
	user.Password = ""

	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) UpdateUserCredentials(ctx context.Context, user *model.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return s.repo.UpdateUserCredentials(ctx, user)
}

func (s *Service) CheckUser(ctx context.Context, creds model.Credentials) error {
	user, err := s.repo.GetCredentials(ctx, creds.Email)
	if err != nil {
		return err
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))

	if err != nil {
		return fmt.Errorf("Wrong password")
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	return s.repo.GetUser(ctx, userId)
}

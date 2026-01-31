package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/katerinasolntsye/fulleng/internal/repository"
	"github.com/katerinasolntsye/fulleng/internal/repository/model"
)

type Service struct {
	repo       repository.Repository
	jwtService *JWTService
}

func NewService(repo repository.Repository, jwtService *JWTService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
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

// AuthenticateUser проверяет учетные данные и возвращает токены
func (s *Service) AuthenticateUser(ctx context.Context, creds model.Credentials) (accessToken, refreshToken string, userId int64, err error) {
	// Получаем учетные данные пользователя
	user, err := s.repo.GetCredentials(ctx, creds.Email)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid credentials")
	}

	// Получаем userId по email
	userId, err = s.repo.GetUserIdByEmail(ctx, creds.Email)
	if err != nil {
		return "", "", 0, fmt.Errorf("user not found")
	}

	// Генерируем токены
	accessToken, err = s.jwtService.GenerateAccessToken(userId, creds.Email)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token")
	}

	refreshToken, err = s.jwtService.GenerateRefreshToken(userId)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token")
	}

	// Сохраняем refresh token в БД
	err = s.repo.SaveRefreshToken(ctx, userId, refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to save refresh token")
	}

	return accessToken, refreshToken, userId, nil
}

// RefreshTokens обновляет токены используя refresh token
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// Валидируем refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Проверяем, что это именно refresh token
	if claims.Type != RefreshToken {
		return "", "", fmt.Errorf("invalid token type")
	}

	// Проверяем, что токен совпадает с сохраненным в БД
	savedToken, err := s.repo.GetRefreshToken(ctx, claims.UserId)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found")
	}

	if savedToken != refreshToken {
		return "", "", fmt.Errorf("refresh token mismatch")
	}

	// Получаем пользователя для получения email
	user, err := s.repo.GetUser(ctx, claims.UserId)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	// Генерируем новые токены
	newAccessToken, err = s.jwtService.GenerateAccessToken(claims.UserId, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token")
	}

	newRefreshToken, err = s.jwtService.GenerateRefreshToken(claims.UserId)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token")
	}

	// Сохраняем новый refresh token в БД
	err = s.repo.SaveRefreshToken(ctx, claims.UserId, newRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to save refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout удаляет refresh token из БД
func (s *Service) Logout(ctx context.Context, userId int64) error {
	return s.repo.DeleteRefreshToken(ctx, userId)
}

func (s *Service) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	return s.repo.GetUser(ctx, userId)
}

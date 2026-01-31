package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/katerinasolntsye/fulleng/internal/repository"
	"github.com/katerinasolntsye/fulleng/internal/repository/model"
)

// SQL query constants
const (
	queryGetIncomingPostbacks     = "SELECT * FROM incoming_postback"
	queryGetTrackers              = "SELECT * FROM tracker"
	queryGetSendPostbacks         = "SELECT * FROM send_postback"
	queryCreateUser               = "INSERT INTO users(email, password, name, timezone) VALUES ($1, $2, $3, $4)"
	queryUpdateUser               = "UPDATE users SET name = $2, surname = $3, birth_date = $4, country_id = $5, city_id = $6, google_id = $7, vkontakte_id = $8, telegram_id = $9, timezone = '+3' WHERE id = $1"
	queryUpdateUserCredentials    = "UPDATE users SET email = $2, tel = $3 WHERE id = $1"
	queryUpdateUserCredentialsPwd = "UPDATE users SET email = $2, password = $3, tel = $4 WHERE id = $1"
	queryGetCredentials           = "SELECT email, password FROM users WHERE email = $1"
	queryGetUser                  = "SELECT id, email, tel, name, surname, birth_date, country_id, city_id, google_id, vkontakte_id, telegram_id, created_at, timezone FROM users WHERE id = $1"
	queryGetUserIdByEmail         = "SELECT id FROM users WHERE email = $1"
	querySaveRefreshToken         = "UPDATE users SET refresh_token = $2 WHERE id = $1"
	queryGetRefreshToken          = "SELECT refresh_token FROM users WHERE id = $1"
	queryDeleteRefreshToken       = "UPDATE users SET refresh_token = NULL WHERE id = $1"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) repository.Repository {
	return &Repository{db: db}
}

func (r *Repository) GetIncomingPostbacks(ctx context.Context) ([]model.IncomingPostback, error) {
	var incoming []model.IncomingPostback

	res, err := r.db.Query(ctx, queryGetIncomingPostbacks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var postback model.IncomingPostback
		err := res.Scan(&postback.Id, &postback.TrackerId, &postback.CnvStatus, &postback.Payout, &postback.Currency, &postback.UrlQuery, &postback.RequestIp, &postback.CreatedAt, &postback.ClickId)
		if err != nil {
			log.Printf("Error scanning incoming postback: %v", err)
			continue
		}
		incoming = append(incoming, postback)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return incoming, nil
}

func (r *Repository) GetTrackers(ctx context.Context) ([]model.Tracker, error) {
	var trackers []model.Tracker

	res, err := r.db.Query(ctx, queryGetTrackers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var tracker model.Tracker
		err := res.Scan(&tracker.Id, &tracker.IsActive, &tracker.TrackerName, &tracker.PostbackTemplate)
		if err != nil {
			log.Printf("Error scanning tracker: %v", err)
			continue
		}
		trackers = append(trackers, tracker)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return trackers, nil
}

func (r *Repository) GetSendPostbacks(ctx context.Context) ([]model.SendPostback, error) {
	var sendPostbacks []model.SendPostback

	res, err := r.db.Query(ctx, queryGetSendPostbacks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var postback model.SendPostback
		err := res.Scan(&postback.Id, &postback.IncomingPostbackId, &postback.TrackerId, &postback.RequestUrl, &postback.ResponseBody, &postback.ResponseCode, &postback.CreatedAt)
		if err != nil {
			log.Printf("Error scanning send postback: %v", err)
			continue
		}
		sendPostbacks = append(sendPostbacks, postback)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return sendPostbacks, nil
}

func (r *Repository) CreateUser(ctx context.Context, email, hashedPassword, name, timezone string) error {
	_, err := r.db.Exec(ctx, queryCreateUser, email, hashedPassword, name, timezone)
	return err
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	println(user.Surname.String)

	_, err := r.db.Exec(ctx, queryUpdateUser, user.Id, user.Name, user.Surname, user.BirthDate, user.CountryId, user.CityId, user.GoogleId, user.VkontakteId, user.TelegramId)
	return err
}

func (r *Repository) UpdateUserCredentials(ctx context.Context, user *model.User) error {
	if user.Password != "" {
		_, err := r.db.Exec(ctx, queryUpdateUserCredentialsPwd, user.Id, user.Email, user.Password, user.Tel)
		return err
	}
	_, err := r.db.Exec(ctx, queryUpdateUserCredentials, user.Id, user.Email, user.Tel)
	return err
}

func (r *Repository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	res, err := r.db.Query(ctx, queryGetUser, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	defer res.Close()

	var user model.User

	for res.Next() {
		err := res.Scan(&user.Id, &user.Email, &user.Tel, &user.Name, &user.Surname, &user.BirthDate, &user.CountryId, &user.CityId, &user.GoogleId, &user.VkontakteId, &user.TelegramId, &user.CreatedAt, &user.Timezone)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		// return &user, nil
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	if user.Email != "" {
		return &user, nil
	}

	return nil, fmt.Errorf("User not found")
}

func (r *Repository) GetCredentials(ctx context.Context, email string) (*model.Credentials, error) {
	res, err := r.db.Query(ctx, queryGetCredentials, email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	defer res.Close()

	var creds model.Credentials

	for res.Next() {
		err := res.Scan(&creds.Email, &creds.Password)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		// return &creds, nil
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	if creds.Email != "" {
		return &creds, nil
	}

	return nil, fmt.Errorf("User not found")
}

func (r *Repository) GetUserIdByEmail(ctx context.Context, email string) (int64, error) {
	var userId int64
	err := r.db.QueryRow(ctx, queryGetUserIdByEmail, email).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}
	return userId, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userId int64, refreshToken string) error {
	_, err := r.db.Exec(ctx, querySaveRefreshToken, userId, refreshToken)
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, userId int64) (string, error) {
	var refreshToken sql.NullString
	err := r.db.QueryRow(ctx, queryGetRefreshToken, userId).Scan(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	if !refreshToken.Valid {
		return "", fmt.Errorf("refresh token not found")
	}
	return refreshToken.String, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, queryDeleteRefreshToken, userId)
	return err
}

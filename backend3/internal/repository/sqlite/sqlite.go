package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/katerinasolntsye/fulleng/internal/repository"
	"github.com/katerinasolntsye/fulleng/internal/repository/model"
	_ "github.com/mattn/go-sqlite3"
)

// SQL query constants
const (
	queryGetIncomingPostbacks     = "SELECT * FROM incoming_postback"
	queryGetTrackers              = "SELECT * FROM tracker"
	queryGetSendPostbacks         = "SELECT * FROM send_postback"
	queryCreateUser               = "INSERT INTO users(email, password, name, timezone) VALUES (?, ?, ?, ?)"
	queryUpdateUser               = "UPDATE users SET name = ?, surname = ?, birth_date = ?, country_id = ?, city_id = ?, google_id = ?, vkontakte_id = ?, telegram_id = ?, timezone = '+3' WHERE id = ?"
	queryUpdateUserCredentials    = "UPDATE users SET email = ?, tel = ? WHERE id = ?"
	queryUpdateUserCredentialsPwd = "UPDATE users SET email = ?, password = ?, tel = ? WHERE id = ?"
	queryGetCredentials           = "SELECT email, password FROM users WHERE email = ?"
	queryGetUser                  = "SELECT id, email, tel, name, surname, birth_date, country_id, city_id, google_id, vkontakte_id, telegram_id, created_at, timezone FROM users WHERE id = ?"
	queryGetUserIdByEmail         = "SELECT id FROM users WHERE email = ?"
	querySaveRefreshToken         = "UPDATE users SET refresh_token = ? WHERE id = ?"
	queryGetRefreshToken          = "SELECT refresh_token FROM users WHERE id = ?"
	queryDeleteRefreshToken       = "UPDATE users SET refresh_token = NULL WHERE id = ?"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository.Repository {
	return &Repository{db: db}
}

func (r *Repository) GetIncomingPostbacks(ctx context.Context) ([]model.IncomingPostback, error) {
	var incoming []model.IncomingPostback

	rows, err := r.db.QueryContext(ctx, queryGetIncomingPostbacks)
	if err != nil {
		log.Printf("QueryContext failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postback model.IncomingPostback
		err := rows.Scan(&postback.Id, &postback.TrackerId, &postback.CnvStatus, &postback.Payout, &postback.Currency, &postback.UrlQuery, &postback.RequestIp, &postback.CreatedAt, &postback.ClickId)
		if err != nil {
			log.Printf("Error scanning incoming postback: %v", err)
			continue
		}
		incoming = append(incoming, postback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incoming, nil
}

func (r *Repository) GetTrackers(ctx context.Context) ([]model.Tracker, error) {
	var trackers []model.Tracker

	rows, err := r.db.QueryContext(ctx, queryGetTrackers)
	if err != nil {
		log.Printf("QueryContext failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tracker model.Tracker
		err := rows.Scan(&tracker.Id, &tracker.IsActive, &tracker.TrackerName, &tracker.PostbackTemplate)
		if err != nil {
			log.Printf("Error scanning tracker: %v", err)
			continue
		}
		trackers = append(trackers, tracker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trackers, nil
}

func (r *Repository) GetSendPostbacks(ctx context.Context) ([]model.SendPostback, error) {
	var sendPostbacks []model.SendPostback

	rows, err := r.db.QueryContext(ctx, queryGetSendPostbacks)
	if err != nil {
		log.Printf("QueryContext failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postback model.SendPostback
		err := rows.Scan(&postback.Id, &postback.IncomingPostbackId, &postback.TrackerId, &postback.RequestUrl, &postback.ResponseBody, &postback.ResponseCode, &postback.CreatedAt)
		if err != nil {
			log.Printf("Error scanning send postback: %v", err)
			continue
		}
		sendPostbacks = append(sendPostbacks, postback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sendPostbacks, nil
}

func (r *Repository) CreateUser(ctx context.Context, email, hashedPassword, name, timezone string) error {
	_, err := r.db.ExecContext(ctx, queryCreateUser, email, hashedPassword, name, timezone)
	return err
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, queryUpdateUser, user.Name, user.Surname, user.BirthDate, user.CountryId, user.CityId, user.GoogleId, user.VkontakteId, user.TelegramId, user.Id)
	return err
}

func (r *Repository) UpdateUserCredentials(ctx context.Context, user *model.User) error {
	if user.Password != "" {
		_, err := r.db.ExecContext(ctx, queryUpdateUserCredentialsPwd, user.Email, user.Password, user.Tel, user.Id)
		return err
	}
	_, err := r.db.ExecContext(ctx, queryUpdateUserCredentials, user.Email, user.Tel, user.Id)
	return err
}

func (r *Repository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var user model.User

	err := r.db.QueryRowContext(ctx, queryGetUser, id).Scan(
		&user.Id, &user.Email, &user.Tel, &user.Name, &user.Surname,
		&user.BirthDate, &user.CountryId, &user.CityId, &user.GoogleId,
		&user.VkontakteId, &user.TelegramId, &user.CreatedAt, &user.Timezone,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetCredentials(ctx context.Context, email string) (*model.Credentials, error) {
	var creds model.Credentials

	err := r.db.QueryRowContext(ctx, queryGetCredentials, email).Scan(&creds.Email, &creds.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &creds, nil
}

func (r *Repository) GetUserIdByEmail(ctx context.Context, email string) (int64, error) {
	var userId int64
	err := r.db.QueryRowContext(ctx, queryGetUserIdByEmail, email).Scan(&userId)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userId int64, refreshToken string) error {
	_, err := r.db.ExecContext(ctx, querySaveRefreshToken, refreshToken, userId)
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, userId int64) (string, error) {
	var refreshToken sql.NullString
	err := r.db.QueryRowContext(ctx, queryGetRefreshToken, userId).Scan(&refreshToken)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	if !refreshToken.Valid {
		return "", fmt.Errorf("refresh token not found")
	}
	return refreshToken.String, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, userId int64) error {
	_, err := r.db.ExecContext(ctx, queryDeleteRefreshToken, userId)
	return err
}

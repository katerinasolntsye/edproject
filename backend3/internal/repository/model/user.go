package model

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

type Credentials struct {
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type User struct {
	Id           int64          `json:"id" db:"id"`
	Email        string         `json:"email" db:"email"`
	Password     string         `json:"password" db:"password"`
	Tel          sql.NullString `json:"tel" db:"tel"`
	Name         string         `json:"name" db:"name"`
	Surname      sql.NullString `json:"surname" db:"surname"`
	BirthDate    sql.NullTime   `json:"birthDate" db:"birth_date"`
	CountryId    int64          `json:"countryId" db:"country_id"`
	CityId       int64          `json:"cityId" db:"city_id"`
	GoogleId     sql.NullString `json:"googleId" db:"google_id"`
	VkontakteId  sql.NullString `json:"vkontakteId" db:"vkontakte_id"`
	TelegramId   sql.NullString `json:"telegramId" db:"telegram_id"`
	CreatedAt    time.Time      `json:"createdAt" db:"created_at"`
	Timezone     string         `json:"timezone" db:"timezone"`
	RefreshToken sql.NullString `json:"-" db:"refresh_token"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User

	payload := &struct {
		Tel       *string `json:"tel"`
		Surname   *string `json:"surname"`
		BirthDate *string `json:"birthDate"`
		// CountryId   *string `json:"countryId"`
		GoogleId    *string `json:"googleId"`
		VkontakteId *string `json:"vkontakteId"`
		TelegramId  *string `json:"telegramId"`
		*Alias
	}{
		Alias: (*Alias)(&u),
	}

	if u.Tel.Valid {
		payload.Tel = &u.Tel.String
	}

	if u.Surname.Valid {
		payload.Surname = &u.Surname.String
	}

	// if u.CountryId.Valid {
	// 	payload.CountryId = &u.CountryId.String
	// }

	if u.GoogleId.Valid {
		payload.GoogleId = &u.GoogleId.String
	}

	if u.VkontakteId.Valid {
		payload.VkontakteId = &u.VkontakteId.String
	}

	if u.TelegramId.Valid {
		payload.TelegramId = &u.TelegramId.String
	}

	if u.BirthDate.Valid {
		birthDate := u.BirthDate.Time.Format("2006-01-02")
		payload.BirthDate = &birthDate
	}

	return json.Marshal(payload)
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User

	payload := &struct {
		Tel       *string `json:"tel"`
		Surname   *string `json:"surname"`
		BirthDate *string `json:"birthDate"`
		// CountryId   *string `json:"countryId"`
		GoogleId    *string `json:"googleId"`
		VkontakteId *string `json:"vkontakteId"`
		TelegramId  *string `json:"telegramId"`
		CreatedAt   *string `json:"-"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	err := json.Unmarshal(data, &payload)

	if err != nil {
		return err
	}

	if payload.Tel != nil {
		u.Tel = sql.NullString{String: *payload.Tel, Valid: true}
	}

	if payload.Surname != nil {
		u.Surname = sql.NullString{String: *payload.Surname, Valid: true}
	}

	// if payload.CountryId != nil {
	// 	u.CountryId = sql.NullString{String: *payload.CountryId, Valid: true}
	// }

	if payload.GoogleId != nil {
		u.GoogleId = sql.NullString{String: *payload.GoogleId, Valid: true}
	}

	if payload.VkontakteId != nil {
		u.VkontakteId = sql.NullString{String: *payload.VkontakteId, Valid: true}
	}

	if payload.TelegramId != nil {
		u.TelegramId = sql.NullString{String: *payload.TelegramId, Valid: true}
	}

	if payload.BirthDate != nil {
		s := strings.Trim(*payload.BirthDate, `"`)
		birthDate, err := time.Parse("2006-01-02", s)
		if err != nil {
			birthDate, err = time.Parse("02.01.2006", s)
			if err != nil {
				return err
			}
			// return err
		}
		u.BirthDate = sql.NullTime{Time: birthDate, Valid: true}
	}

	return nil
}

package model

import (
	"database/sql"
	"time"
)

type IncomingPostback struct {
	Id        int       `db:"id"`
	TrackerId int       `db:"trackerid"`
	CnvStatus string    `db:"cnv_status"`
	Payout    float64   `db:"payout"`
	Currency  string    `db:"currency"`
	UrlQuery  string    `db:"url_query"`
	RequestIp string    `db:"request_ip"`
	CreatedAt time.Time `db:"created_at"`
	ClickId   string    `db:"clickid"`
}
type Tracker struct {
	Id               int    `db:"id"`
	IsActive         int    `db:"is_active"`
	TrackerName      string `db:"tracker_name"`
	PostbackTemplate string `db:"postback_template"`
}
type SendPostback struct {
	Id                 int       `db:"id"`
	IncomingPostbackId int       `db:"incoming_postback_id"`
	TrackerId          int       `db:"tracker_id"`
	RequestUrl         string    `db:"request_url"`
	ResponseBody       string    `db:"response_body"`
	ResponseCode       int       `db:"response_code"`
	CreatedAt          time.Time `db:"created_at"`
}
type SendPostbackFailed struct {
	TrackerName sql.NullString `db:"tracker_name"`
}

type Countries struct {
	Id          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}

type Cities struct {
	Id          int64  `json:"id" db:"id"`
	CountryId   int64  `json:"countryId" db:"country_id"`
	Description string `json:"description" db:"description"`
}

type Courses struct {
	Id         int64  `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Link       string `json:"link" db:"link"`
	Price      int64  `json:"price" db:"price"`
	Period     int64  `json:"period" db:"period"`
	PeriodType int64  `json:"periodType" db:"period_type"`
}

type Lessons struct {
	Id         int64  `json:"id" db:"id"`
	CourseId   int64  `json:"courseId" db:"course_id"`
	Title      string `json:"title" db:"title"`
	Link       string `json:"link" db:"link"`
	Number     int64  `json:"number" db:"number"`
	LessonText string `json:"lessonText" db:"lesson_text"`
}

type Tasks struct {
	Id       int64  `json:"id" db:"id"`
	CourseId int64  `json:"courseId" db:"course_id"`
	LessonId int64  `json:"lessonId" db:"lesson_id"`
	Title    string `json:"title" db:"title"`
	TaskType int64  `json:"type" db:"type"`
	Number   int64  `json:"number" db:"number"`
	Goal     int64  `json:"goal" db:"goal"`
	Data     string `json:"data" db:"data"`
}

type CompletedTasks struct {
	UserId        int64     `json:"userId" db:"user_id"`
	TaskId        int64     `json:"taskId" db:"task_id"`
	Completations int64     `json:"completations" db:"completations"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type Orders struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"userId" db:"user_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	TotalPrice  int64     `json:"totalPrice" db:"total_price"`
	PromocodeId string    `json:"promocodeId" db:"promocode_id"`
	PaymentType int64     `json:"paymentType" db:"payment_type"`
	PaidAt      time.Time `json:"paidAt" db:"paid_at"`
}

type OrderItems struct {
	Id        int64     `json:"id" db:"id"`
	OrderId   int64     `json:"orderId" db:"order_id"`
	Price     int64     `json:"price" db:"price"`
	CourseId  int64     `json:"courseId" db:"course_id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

type OrderTransactions struct {
	Id           int64     `json:"id" db:"id"`
	OrderId      int64     `json:"orderId" db:"order_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	Stage        int64     `json:"stage" db:"stage"`
	Status       int64     `json:"status" db:"status"`
	ProcessingId string    `json:"processingId" db:"processing_id"`
}

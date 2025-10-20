package storage

import (
	"database/sql"
	"time"
)

type Comic struct {
	Id               int
	Name             string
	Nsfw             bool
	SeriesId         int    `db:"series_id"`
	CoverImageBase64 string `db:"cover_image_base64"`
	NumPages         int    `db:"num_pages"`
	FilePath         string `db:"file_path"`
	Tags             []string
	LastPage         int            `db:"last_page"`
	OwnerID          int            `db:"ownerid"`
	OrderNum         int            `db:"ordernum"`
	Sha256sum        sql.NullString `db:"sha256sum"`
	Finished         bool           `db:"finished"`
	ArcId            int            `db:"arcid"`
}

type Series struct {
	Id               int
	Name             string
	CoverImageBase64 string `db:"cover_image_base64"`
	ComicsInSeries   int    `db:"comics_in_series"`
	Nsfw             bool   `db:"nsfw"`
	OrderNum         int    `db:"ordernum"`
	OwnerID          int    `db:"ownerid"`
}

type GnolJob struct {
	Id        int
	JobStatus int    `db:"job_status"`
	UserID    int    `db:"user_id"`
	JobType   int    `db:"job_type"`
	Data      string `db:"input_data"`
}

type GnolSession struct {
	SessionId  string    `db:"session_id"`
	ValidUntil time.Time `db:"valid_until"`
	UserId     int       `db:"user_id"`
}

type SeriesArc struct {
	Id          int            `db:"id"`
	OrderNum    int            `db:"ordernum"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	FandomLink  sql.NullString `db:"fandom_link"`
	SeriesId    int            `db:"series_id"`
}

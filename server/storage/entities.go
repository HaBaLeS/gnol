package storage

type Comic struct {
	Id               int
	Name             string
	Nsfw             bool
	SeriesId         int    `db:"series_id"`
	CoverImageBase64 string `db:"cover_image_base64"`
	NumPages         int    `db:"num_pages"`
	FilePath         string `db:"file_path"`
	Tags             []string
	LastPage         int `db:"last_page"`
	OwnerID          int `db:"ownerID"`
}

type Series struct {
	Id               int
	Name             string
	CoverImageBase64 string `db:"cover_image_base64"`
	ComicsInSeries   int    `db:"comics_in_series"`
}

type GnolJob struct {
	Id        int
	JobStatus int    `db:"job_status"`
	UserID    int    `db:"user_id"`
	JobType   int    `db:"job_type"`
	Data      string `db:"input_data"`
}

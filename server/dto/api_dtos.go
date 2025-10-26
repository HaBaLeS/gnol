package dto

import (
	"github.com/HaBaLeS/gnol/server/database"
)

type ComicDTO struct {
	Id        int
	Name      string
	Series_id int
	Sname     string
	Nsfw      bool
	Num_pages int
	Sha256sum string
}

type ArcDTO struct {
	*database.SeriesArc
	Comics []*database.Comic
}

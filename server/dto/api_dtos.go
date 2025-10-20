package dto

import "github.com/HaBaLeS/gnol/server/storage"

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
	*storage.SeriesArc
	Comics []*storage.Comic
}

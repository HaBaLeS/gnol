package dto

type ComicEntry struct {
	Id        int
	Name      string
	Series_id int
	Sname     string
	Nsfw      bool
	Num_pages int
	Sha256sum string
}

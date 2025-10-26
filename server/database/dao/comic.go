package dao

import (
	"fmt"

	"github.com/HaBaLeS/gnol/server/database"
)

const (
	ADD_USER_2_COMIC    = "insert into user_to_comic (user_id, comic_id) values ($1, $2)"
	ALL_COMICS_FOR_USER = `
	select  c.*, utc.last_page from comic as c
    	join user_to_comic utc on c.id = utc.comic_id
    	where utc.user_id = $1
	`

	COMICS_FOR_USER_IN_SERIES = `
select  c.*, utc.last_page, utc.finished from comic as c
    join user_to_comic utc on c.id = utc.comic_id
    where utc.user_id = $1 and series_id = $2 
	order by c.ordernum asc
`

	CREATE_COMIC = "insert into comic (id, name, nsfw, series_id, cover_image_base64, num_pages, file_path, ordernum ) values ($1, $2, $3, $4, $5, $6, $7, $8)"
)

func (dao *DAO) AddComicToUser(comicID string, userID string) error {
	_, err := dao.DB.Exec(ADD_USER_2_COMIC, userID, comicID)
	return err
}

func (dao *DAO) ComicById(id string) *database.Comic {
	c := &database.Comic{}
	err := dao.DB.Get(c, "select * from comic where id = $1", id)
	if err != nil {
		panic(fmt.Errorf("get comic by id failed: %v", err))
	}
	return c
}

func (dao *DAO) ComicFilenameForId(id string) string {
	var retVal string
	err := dao.DB.Get(&retVal, "select file_path from comic where id = $1", id)
	if err != nil {
		panic(fmt.Errorf("get comic by id failed: %v", err))
	}
	return retVal
}

func (dao *DAO) ComicsForUser(id int) []*database.Comic {
	retList := make([]*database.Comic, 0)
	if err := dao.DB.Select(&retList, ALL_COMICS_FOR_USER, id); err != nil {
		dao.log.Printf("SQL Error, %v", err)
	}

	//--  dont knwo if we should do that lazy or direct \o/
	for _, c := range retList {
		q := "select tag as Tags from tags join tag_to_comic ttc on tags.Id = ttc.tag_id where ttc.comic_id =$1"
		err := dao.DB.Select(&c.Tags, q, c.Id)
		if err != nil {
			panic(err)
		}
	}

	return retList
}

func (dao *DAO) ComicsForUserInSeries(id int, seriesID string) []*database.Comic {

	retList := make([]*database.Comic, 0)
	if err := dao.DB.Select(&retList, COMICS_FOR_USER_IN_SERIES, id, seriesID); err != nil {
		dao.log.Printf("SQL Errror, %v", err)
	}

	//--  dont know if we should do that lazy or direct \o/
	for _, c := range retList {
		q := "select tag as Tags from tags join tag_to_comic ttc on tags.Id = ttc.tag_id where ttc.comic_id =$1"
		err := dao.DB.Select(&c.Tags, q, c.Id)
		if err != nil {
			panic(err)
		}
	}
	return retList
}

func (dao *DAO) SaveComic(c *database.Comic) int {
	//insert into commic
	var newID int
	dao.DB.Get(&newID, "select nextval('comic_id_seq')")
	dao.DB.MustExec(CREATE_COMIC, newID, c.Name, c.Nsfw, c.SeriesId, c.CoverImageBase64, c.NumPages, c.FilePath, c.OrderNum)
	return newID
}

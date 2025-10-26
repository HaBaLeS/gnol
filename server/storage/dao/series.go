package dao

import (
	"fmt"

	"github.com/HaBaLeS/gnol/server/storage"
)

const (
	CREATE_SERIES = "insert into series (name, cover_image_base64) values ($1,$2)"

	SERIES_FOR_USER      = "select s.*, count(s.id) as comics_in_series from comic join series s on s.id = comic.series_id left join user_to_comic utc on comic.id = utc.comic_id where utc.user_id = $1 and s.nsfw = false group by s.id order by s.ordernum asc"
	SERIES_FOR_USER_NSFW = "select s.*, count(s.id) as comics_in_series from comic join series s on s.id = comic.series_id left join user_to_comic utc on comic.id = utc.comic_id where utc.user_id = $1 and s.nsfw = true group by s.id order by s.ordernum asc"
)

func (dao *DAO) AllSeries(includeNSFW bool) []*storage.Series {
	retList := make([]*storage.Series, 0)
	err := dao.DB.Select(&retList, "select id, name from series where nsfw = $1 or nsfw = false order by name", includeNSFW)
	if err != nil {
		panic(err)
	}
	return retList
}

func (dao *DAO) CreateSeries(name, imageB64 string) (int, error) {
	//fixme how do we do duplicates Names?
	res := dao.DB.MustExec(CREATE_SERIES, name, imageB64)
	id, err := res.LastInsertId()
	return int(id), err
}

func (dao *DAO) SeriesForUser(id int) []*storage.Series {
	retList := make([]*storage.Series, 0)
	if err := dao.DB.Select(&retList, SERIES_FOR_USER, id); err != nil {
		dao.log.Printf("SQL Errror, %v", err)
	}
	return retList
}

func (dao *DAO) NSFWSeriesForUser(id int) []*storage.Series {
	retList := make([]*storage.Series, 0)
	if err := dao.DB.Select(&retList, SERIES_FOR_USER_NSFW, id); err != nil {
		dao.log.Printf("SQL Errror, %v", err)
	}
	return retList
}

func (dao *DAO) SeriesByIdAndUser(seriesId string, userId int) (*storage.Series, bool) {
	retSerie := &storage.Series{}
	err := dao.DB.Get(retSerie, "select * from series s where s.id = $1 and s.ownerid  = $2", seriesId, userId)
	if err != nil {
		//You are not the owner!
		return nil, false
	}
	return retSerie, true
}

func (dao *DAO) SeriesInfoById(id string) *storage.Series {
	retSerie := &storage.Series{}
	err := dao.DB.Get(retSerie, "select id, name, nsfw from series s where s.id = $1", id)
	if err != nil {
		panic(fmt.Errorf("request for non existent series id %s", id))
	}
	return retSerie
}

func (dao *DAO) AddSeriesArc(seriesId, name string) {
	dao.DB.MustExec("INSERT INTO series_arc (series_id, name) values ($1, $2)", seriesId, name)
}

func (dao *DAO) ListSeriesArcs(seriesId string) []*storage.SeriesArc {
	retList := make([]*storage.SeriesArc, 0)

	err := dao.DB.Select(&retList, "select * from series_arc where series_id = $1 order by ordernum, name asc", seriesId)
	if err != nil {
		panic(err)
	}

	return retList
}

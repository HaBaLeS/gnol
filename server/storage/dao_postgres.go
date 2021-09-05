package storage

import (
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const(
	SERIES_FOR_USER = "select s.*, count(s.id) as comics_in_series from comic join series s on s.id = comic.series_id left join user_to_comic utc on comic.id = utc.comic_id where utc.user_id = ? group by s.id"
	COMICS_FOR_USER = "select c.* from comic as c join user_to_comic utc on c.id = utc.comic_id and utc.user_id = ?"
	CREATE_COMIC = "insert into comic (name, nsfw, series_id, cover_image_base64, num_pages, file_path) values ($1, $2, $3, $4, $5, $6)"
	ADD_USER_2_COMIC = "insert into user_to_comic (user_id, comic_id) values ($1, $2)"
)


var schema = `

DROP TABLE IF EXISTS "gnoluser";
DROP TABLE IF EXISTS "comic";
DROP TABLE IF EXISTS "series";
DROP TABLE IF EXISTS "user_to_comic";

CREATE TABLE "gnoluser"(
   id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
   name text NOT NULL,
   password_hash bytea,
   salt bytea
);

CREATE TABLE "comic"(
   id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
   name text NOT NULL,
   series_id INTEGER,
   nsfw bool,
   cover_image_base64 TEXT,
   num_pages         INTEGER,
   file_path TEXT
);

CREATE TABLE "series"(
   id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
   name text NOT NULL,
   cover_image_base64 TEXT 
);

CREATE TABLE "user_to_comic"(
   user_id int NOT NULL,
   comic_id int NOT NULL,
   CONSTRAINT user_to_comic_pkey PRIMARY KEY (user_id,comic_id)
);


CREATE UNIQUE INDEX gnoluser_name_key ON "gnoluser"(name);

`






type DAO struct {
	log *log.Logger
	DB *sqlx.DB
}

func NewDAO(cfg *util.ToolConfig) *DAO{
	dao := &DAO{}
	dao.init()
	return dao
}


func (dao *DAO) init() {
	db, err := sqlx.Connect("sqlite3", "__deleteme.db")
	if err != nil {
		log.Fatalln(err)
	}
	dao.DB = db
	dao.log = log.Default()

	db.MustExec(schema)
	dao.AddUser("falko","falko")
}

func (dao *DAO) ComicsForUser(id int) *[]Comic {
	retList := make([]Comic,0)
	if err := dao.DB.Select(&retList, COMICS_FOR_USER, id); err!= nil {
		dao.log.Printf("SQL Errror, %v", err)
	}
	return &retList
}

func (dao *DAO) SeriesForUser(id int) *[]Series {
	retList := make([]Series,0)
	if err := dao.DB.Select(&retList, SERIES_FOR_USER, id); err!= nil {
		dao.log.Printf("SQL Errror, %v", err)
	}
	return &retList
}

func (dao *DAO) SaveComic(c *Comic) (int, error){
	//insert into commic
	res:= dao.DB.MustExec(CREATE_COMIC, c.Name, c.Nsfw, c.SeriesId, c.CoverImageBase64, c.NumPages, c.FilePath)
	id, err := res.LastInsertId()
	return int(id), err
}

func (dao *DAO) AddComicToUser(comicID int, userID int) error{
	 _, err := dao.DB.Exec(ADD_USER_2_COMIC, userID, comicID)
	 return err
}

func (dao *DAO) ComicById(id string) *Comic {
	c := &Comic{}
	dao.DB.Get(c,"select * from comic where id = $1", id)
	return c
}

type User struct {
	Id int
	Name    string
	PasswordHash []byte `db:"password_hash"`
	Salt    []byte
}

type Comic struct {
	Id int
	Name string
	Nsfw bool
	SeriesId int `db:"series_id"`
	CoverImageBase64 string `db:"cover_image_base64"`
	NumPages         int `db:"num_pages"`
	FilePath string `db:"file_path"`
}

type Series struct {
	Id int
	Name string
	CoverImageBase64 string `db:"cover_image_base64"`
	ComicsInSeries int `db:"comics_in_series"`
}
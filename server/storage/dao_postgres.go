package storage

import (
	"database/sql"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const(
	SERIES_FOR_USER = "select s.*, count(s.id) as comics_in_series from comic join series s on s.id = comic.series_id left join user_to_comic utc on comic.id = utc.comic_id where utc.user_id = ? group by s.id"
	COMICS_FOR_USER = "select c.* from comic as c join user_to_comic utc on c.id = utc.comic_id and utc.user_id = ?"
	ADD_USER_2_COMIC = "insert into user_to_comic (user_id, comic_id) values ($1, $2)"
	UPDATE_VERSION = "insert into schema_version (version) values ($1)"
	CURRENT_VERSION = "select max(version) from schema_version"

	OLDEST_OPEN_JOB ="select * from gnoljobs where job_status = 0 order by id asc limit 1"
	UPDATE_JOB_STATUS = "update gnoljobs set job_status = $1 where id = $2"

	CREATE_COMIC = "insert into comic (name, nsfw, series_id, cover_image_base64, num_pages, file_path) values ($1, $2, $3, $4, $5, $6)"
	CREATE_JOB = "insert into gnoljobs (user_id, job_type, input_data) values ($1,$2,$3);"

)


var schema_1 = `

DROP TABLE IF EXISTS "schema_version";
DROP TABLE IF EXISTS "gnoluser";
DROP TABLE IF EXISTS "comic";
DROP TABLE IF EXISTS "series";
DROP TABLE IF EXISTS "user_to_comic";

CREATE TABLE "schema_version" (
    id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
    version INTEGER NOT NULL,
  	migration_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "gnoluser"(
   id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
   name text NOT NULL,
   password_hash bytea,
   salt bytea,
   webauthn bool DEFAULT false                  	
);

CREATE TABLE "comic"(
   id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
   name text NOT NULL,
   series_id INTEGER,
   nsfw bool,
   cover_image_base64 TEXT,
   num_pages INTEGER,
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

var schema_2 = `

DROP TABLE IF EXISTS "gnoljobs";
CREATE TABLE "gnoljobs" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	job_status INTEGER NOT NULL DEFAULT 0,
	user_id int NOT NULL,
	job_type int NOT NULL,
	input_data TEXT
);

`

var schema_3 = `

DROP TABLE IF EXISTS "webauthn_authenticator";
CREATE TABLE "webauthn_authenticator" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    aagu_id  bytea NOT NULL,
    signcount uint32 NOT NULL,
    clonewarning bool default false
);

DROP TABLE IF EXISTS "webauthn_credential";
CREATE TABLE "webauthn_credential" (
    id bytea PRIMARY KEY NOT NULL,
    publicKey bytea NOT NULL,
    attestationType TEXT,
    authenticator_id int NOT NULL,
    user_id int NOT NULL
);
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
	db, err := sqlx.Connect("sqlite3", "gnol_sqlite.db")
	if err != nil {
		log.Fatalln(err)
	}
	dao.DB = db
	dao.log = log.Default()

	var version int

	//version = getVersion
	noVersion := db.Get(&version,CURRENT_VERSION)
	if noVersion == sql.ErrNoRows {
		version = 0
	}

	if version < 1 {
		db.MustExec(schema_1)
		db.MustExec(UPDATE_VERSION, 1)
		dao.AddUser("falko","oklaf")
	}

	if version < 2 {
		db.MustExec(schema_2)
		db.MustExec(UPDATE_VERSION, 2)
	}

	if version < 3 {
		db.MustExec(schema_3)
		db.MustExec(UPDATE_VERSION, 3)
	}

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

func (dao *DAO) CreateJob(jtype, juser int, data string) error{
	_, err := dao.DB.Exec(CREATE_JOB, juser, jtype, data)
	return err
}

func (dao *DAO) GetOldestOpenJob() *GnolJob{
	job := new(GnolJob)
	err := dao.DB.Get(job, OLDEST_OPEN_JOB)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return job
}

func (dao *DAO) UpdatJobStatus(job *GnolJob, newStatus int) {
	dao.DB.MustExec(UPDATE_JOB_STATUS, newStatus, job.Id)
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


type GnolJob struct {
	Id int
	JobStatus int `db:"job_status"`
	UserID int `db:"user_id"`
	JobType int `db:"job_type"`
	Data string `db:"input_data"`
}


package dao

import (
	"fmt"
	"log"

	"github.com/HaBaLeS/gnol/server/database/migration"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DAO struct {
	log *log.Logger
	DB  *sqlx.DB
	cfg *util.ToolConfig
}

func NewDAO(cfg *util.ToolConfig) *DAO {
	dao := &DAO{
		cfg: cfg,
	}
	dao.init()
	return dao
}

func (dao *DAO) init() {

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dao.cfg.PostgresUser, dao.cfg.PostgresPass, dao.cfg.PostgresHost, dao.cfg.PostgresPort, dao.cfg.PostgresDB)
	connConfig, err := pgx.ParseConfig(url)
	if err != nil {
		panic(err)
	}
	dbc := stdlib.OpenDB(*connConfig)
	db := sqlx.NewDb(dbc, "pgx")

	dao.DB = db
	dao.log = log.Default()

	//put migrations in extra file
	migration.Migrate(db)

}

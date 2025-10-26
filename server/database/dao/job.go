package dao

import (
	"database/sql"

	"github.com/HaBaLeS/gnol/server/database"
)

const (
	OLDEST_OPEN_JOB   = "select * from gnoljobs where job_status = 0 order by id asc limit 1"
	UPDATE_JOB_STATUS = "update gnoljobs set job_status = $1 where id = $2"
	CREATE_JOB        = "insert into gnoljobs (user_id, job_type, input_data) values ($1,$2,$3);"
)

func (dao *DAO) SetFinished(userID int, comicID string) {
	dao.DB.MustExec("update user_to_comic set finished = true where user_id=$1 and comic_id = $2", userID, comicID)
}

func (dao *DAO) ToggleFinished(userID, comicID string) {
	dao.DB.MustExec("update user_to_comic set finished = !finished where user_id=$1 and comic_id = $2", userID, comicID)
}

func (dao *DAO) CreateJob(jtype, juser int, data string) error {
	_, err := dao.DB.Exec(CREATE_JOB, juser, jtype, data)
	return err
}

func (dao *DAO) GetOldestOpenJob() *database.GnolJob {
	job := new(database.GnolJob)
	err := dao.DB.Get(job, OLDEST_OPEN_JOB)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return job
}

func (dao *DAO) UpdatJobStatus(job *database.GnolJob, newStatus int) {
	dao.DB.MustExec(UPDATE_JOB_STATUS, newStatus, job.Id)
}

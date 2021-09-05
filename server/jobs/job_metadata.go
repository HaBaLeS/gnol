package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/mholt/archiver/v3"
	"os"
)

//CreateNewArchiveJob create a prepared Job form processing a new CBR/CBZ/RAR/ZIP file
func (j *JobRunner) CreateNewArchiveJob(archive string, userID int) {
	bgjob := &BGJob{
		JobType:     ScanMeta,
		InputFile:   archive,
		DisplayName: "Scan Metadata",
		JobStatus:   NotStarted,
		ExtraData:   make(map[string]string, 10),
		BaseEntity:  storage.CreateBaseEntity(bucketJobOpen),
		UserID: userID,
	}

	j.save(bgjob)
}

func (j *JobRunner) scanMetaData(job *BGJob) error {

	f, err :=  os.Open(job.InputFile)
	if err != nil {
		return err
	}
	eIface,err := archiver.ByHeader(f)
	if err != nil {
		return err
	}

	e, ok := eIface.(archiver.Walker)
	if !ok {
		return fmt.Errorf("format specified by source filename is not an extractor format: %s (%T)", job.InputFile, eIface)
	}

	c := &storage.Comic{}
	err = e.Walk(job.InputFile, func(f archiver.File) error {
		if f.Name() == "gnol.json" {
			dec := json.NewDecoder(f)
			err := dec.Decode(c)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	c.FilePath = job.InputFile
	id, err :=  j.dao.SaveComic(c)
	if err != nil {
		return err
	}
	err = j.dao.AddComicToUser(id, job.UserID)
	return err
}
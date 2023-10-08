package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/mholt/archiver/v3"
	"os"
	"strings"
)

// CreateNewArchiveJob create a prepared Job form processing a new CBR/CBZ/RAR/ZIP file
func (j *JobRunner) CreateNewArchiveJob(archive string, userID int) {
	bgjob := &storage.GnolJob{
		JobType:   ScanMeta,
		Data:      archive,
		JobStatus: NotStarted,
		UserID:    userID,
	}
	j.save(bgjob)
}

func (j *JobRunner) scanMetaData(job *storage.GnolJob) error {

	f, err := os.Open(job.Data)
	if err != nil {
		return err
	}
	eIface, err := archiver.ByHeader(f)
	if err != nil {
		return err
	}

	e, ok := eIface.(archiver.Walker)
	if !ok {
		return fmt.Errorf("format specified by source filename is not an extractor format: %s (%T)", job.Data, eIface)
	}

	c := &storage.Comic{}
	err = e.Walk(job.Data, func(f archiver.File) error {
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
	c.FilePath = job.Data
	id := j.dao.SaveComic(c)
	if err != nil {
		return err
	}

	//get existing Tags
	var knownTags []string
	err = j.dao.DB.Select(&knownTags, "select lower(Tag) from tags")
	if err != nil {
		return err
	}
	for _, tag := range c.Tags {
		if tag == "" {
			continue
		}
		added := false
		for _, kt := range knownTags {
			if strings.ToLower(tag) == kt {
				j.dao.DB.MustExec("insert into tag_to_comic (comic_id, tag_id) values ($1, (select Id from tags where lower(tag) = lower($2)))", id, kt)
				added = true
				break
			}
		}
		if !added {
			//Tag does not exist, lets add it
			j.dao.DB.MustExec("insert into tags (Tag) values ($1)", tag)
			j.dao.DB.MustExec("insert into tag_to_comic (comic_id, tag_id) values ($1, (select Id from tags where lower(tag) = lower($2)))", id, tag)
		}
	}

	err = j.dao.AddComicToUser(id, job.UserID)
	return err
}

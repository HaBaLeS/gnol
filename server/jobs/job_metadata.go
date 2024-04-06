package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/mholt/archiver/v3"
	"os"
	"strconv"
	"strings"
)

type JobMeta struct {
	Filename string `json:"filename"`
	SeriesId int    `json:"seriesId"`
	OrderNum int    `json:"orderNum"`
	Nsfw     bool   `json:"nsfw"`
}

// CreateNewArchiveJob create a prepared Job form processing a new CBR/CBZ/RAR/ZIP file
func (j *JobRunner) CreateNewArchiveJob(jobMeta *JobMeta, userID int) error {
	data, err := json.Marshal(jobMeta)
	if err != nil {
		return err
	}
	bgjob := &storage.GnolJob{
		JobType:   ScanMeta,
		Data:      string(data),
		JobStatus: NotStarted,
		UserID:    userID,
	}
	j.save(bgjob)
	return nil
}

func (j *JobRunner) scanMetaData(jdesc *storage.GnolJob) error {

	jm := &JobMeta{}
	err := json.Unmarshal([]byte(jdesc.Data), jm)
	if err != nil {
		return err
	}

	f, err := os.Open(jm.Filename)
	if err != nil {
		return err
	}
	eIface, err := archiver.ByHeader(f)
	if err != nil {
		return err
	}

	e, ok := eIface.(archiver.Walker)
	if !ok {
		return fmt.Errorf("format specified by source filename is not an extractor format: %s (%T)", jm.Filename, eIface)
	}
	c := &storage.Comic{}
	err = e.Walk(jm.Filename, func(f archiver.File) error {
		if f.Name() == "gnol.json" {
			dec := json.NewDecoder(f)
			err := dec.Decode(c)
			return err
		}
		return nil
	})

	if c.OrderNum == 0 {
		c.OrderNum = jm.OrderNum //fixme only overwrite if certein conditions are meet
		if c.OrderNum == 0 {
			c.OrderNum = 100
		}
	}

	c.SeriesId = jm.SeriesId //fixme only overwrite if certein conditions are meet
	if jm.Nsfw {
		c.Nsfw = true
		c.Tags = append(c.Tags, "nsfw")
	}

	if err != nil {
		return err
	}
	c.FilePath = jm.Filename
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
				_, err := j.dao.DB.Exec("insert into tag_to_comic (comic_id, tag_id) values ($1, (select Id from tags where lower(tag) = lower($2)))", id, kt)
				if err != nil {
					//return err
					j.log.Printf("error writing Tag to comic: '%s' with error %v", kt, err)
				}
				added = true
				break
			}
		}
		if !added {
			//Tag does not exist, lets add it
			_, err := j.dao.DB.Exec("insert into tags (Tag) values ($1)", tag)
			if err != nil {
				//return err
				j.log.Printf("error creating Tag: '%s' with error %v", tag, err)
			}
			_, err = j.dao.DB.Exec("insert into tag_to_comic (comic_id, tag_id) values ($1, (select Id from tags where lower(tag) = lower($2)))", id, tag)
			if err != nil {
				//return err
				j.log.Printf("error writing Tag to comic: '%s' with error %v", tag, err)
			}
			knownTags = append(knownTags, tag)
		}
	}

	err = j.dao.AddComicToUser(strconv.Itoa(id), strconv.Itoa(jdesc.UserID))

	h, err := util.HashFile(c.FilePath)
	if err != nil {
		panic(err)
	}
	j.dao.DB.MustExec("update comic set sha256sum = $1 where id = $2", h, id)

	return err
}

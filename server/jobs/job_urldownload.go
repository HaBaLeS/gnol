package jobs

import (
	"github.com/HaBaLeS/gnol/server/storage"
	"io"
	"net/http"
	"os"
	"path"
)

func (j *JobRunner) CreateNewURLJob(url, uid string) {
	bgjob := &BGJob{
		JobType:     DownloadUrl,
		InputFile:   url,
		DisplayName: "Download URL",
		JobStatus:   NotStarted,
		ExtraData:   make(map[string]string, 10),
		BaseEntity:  storage.CreateBaseEntity(bucketJobOpen),
		UserID: 	 uid,
	}
	j.save(bgjob)
}

func (j *JobRunner) downloadFromUrl(job *BGJob) error {
	uri := job.InputFile
	r, err := http.Get(uri)
	if err != nil {
		j.log.Errorf("Error downloading %s with error, %s\n", uri,err)
		return err
	}
	outName := path.Join(j.cfg.DataDirectory, path.Base(uri))
	wto, cer := os.Create(outName)
	if cer != nil {
		j.log.Errorf("Error creating file %s with error, %s\n", outName,cer)
		return cer
	}

	w,ioe := io.Copy(wto,r.Body)
	if ioe != nil {
		j.log.Errorf("Error copying file, %s", ioe)
		return ioe
	}
	j.log.InfoF("Copyed %d bytes", w)

	//Create followup job after downlaod
	j.CreateNewArchiveJob(outName, job.UserID, "")

	return nil
}

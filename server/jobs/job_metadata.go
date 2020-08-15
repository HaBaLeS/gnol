package jobs

import "github.com/HaBaLeS/gnol/server/storage"

//CreateNewArchiveJob create a prepared Job form processing a new CBR/CBZ/RAR/ZIP file
func (j *JobRunner) CreateNewArchiveJob(archive string, user string, public string) {
	bgjob := &BGJob{
		JobType:     ScanMeta,
		InputFile:   archive,
		DisplayName: "Scan Metadata",
		JobStatus:   NotStarted,
		ExtraData:   make(map[string]string, 10),
		BaseEntity:  storage.CreateBaseEntity(bucketJobOpen),
	}
	bgjob.ExtraData["public"] = public
	bgjob.ExtraData["uploadUser"] = user

	j.save(bgjob)
}

func (j *JobRunner) scanMetaData(job *BGJob) int {
	err, meta := storage.NewMetadata(job.InputFile)
	meta.UploadUser = job.ExtraData["uploadUser"]
	if job.ExtraData["public"] == "public" {
		meta.Public = true
	}
	if err != nil {
		return Error
	}
	err = meta.UpdateMeta()
	if err != nil {
		return Error
	}
	err = j.bs.Comic.SaveComicMeta(meta)

	if err != nil {
		return Error
	}
	//job.env.dao.AddComicToList(meta, meta.UploadUser)
	return Done
}

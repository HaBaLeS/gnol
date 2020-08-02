package conversion

import "github.com/HaBaLeS/gnol/server/dao"

func (j *JobRunner) CreateNewArchiveJob(archive string, user string, public string) {
	bgjob := &BGJob{
		JobType:     ScanMeta,
		InputFile:   archive,
		DisplayName: "Scan Metadata",
		JobStatus:   NotStarted,
		ExtraData:   make(map[string]string, 10),
		BaseEntity:  dao.CreateBaseEntity(),
	}
	bgjob.ExtraData["public"] = public
	bgjob.ExtraData["uploadUser"] = user

	j.save(bgjob)
}

func scanMetaData(job *BGJob) int {
	err, meta := dao.NewMetadata(job.InputFile)
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
	err = job.env.dao.SaveComicMeta(meta)

	if err != nil {
		return Error
	}
	job.env.dao.AddComicToList(meta)
	return Done
}

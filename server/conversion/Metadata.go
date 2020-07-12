package conversion

import "github.com/HaBaLeS/gnol/server/dao"

func (j *JobRunner) CreateNewArchiveJob(archive string) {
	bgjob := &BGJob{
		JobType:     ScanMeta,
		InputFile:   archive,
		DisplayName: "Scan Metadata",
		JobStatus:   NotStarted,
	}
	bgjob.save()
}

func scanMetaData(job *BGJob, daoHandler *dao.DAOHandler){
	err, meta := dao.NewMetadata(job.InputFile)
	if err != nil {
		job.JobStatus = Error
		return
	}
	err = meta.Update()
	if err != nil {
		job.JobStatus = Error
		return
	}
	err = meta.Save()
	if err != nil {
		job.JobStatus = Error
		return
	}
	//FIXME Update DB
	daoHandler.AddComicToList(meta)

	job.JobStatus=Done
}

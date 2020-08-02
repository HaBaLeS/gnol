package conversion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/HaBaLeS/go-logger"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

var JOB_OPEN_BUCKET = []byte("jobs_open")
var JOB_DONE_BUCKET = []byte("jobs_done")
var JOB_ERROR_BUCKET = []byte("jobs_error")

const (
	PdfToCbz = iota
	ScanMeta
)

const (
	NotStarted = iota
	Error
	Done
)

type BGJob struct {
	dao.BaseEntity
	JobType     int    //Convert to CBR, Scan for Metadata, Scrape Meta, SortFolder, Clean Cache etc....
	JobStatus   int    //What's the status done, not started, error
	DisplayName string //Name the Job
	Duration    string //how long did it take
	InputFile   string
	ExtraData   map[string]string
	env         *JobRunner //give access to the jobrunner and Daos
}

//JOBS are defined by creating a name.job json in a special folder. jobs are processed one after the other by reading the directory where the jobs are and
//Takeing one job, reading the description and processing it. MEta is updated while processing ... only 1 job at a time

type JobRunner struct {
	running   bool
	jobLocked bool
	log       *logger.Logger
	dao       *dao.DAOHandler
}

//NewJobRunner Constructor
func NewJobRunner(dao *dao.DAOHandler) *JobRunner { //fixme give config instead of job path
	out, _ := os.Create("jobs.log")
	logger, _ := logger.NewLogger("GnolJob", 0, logger.InfoLevel, out)
	return &JobRunner{
		running:   false,
		jobLocked: false,
		log:       logger,
		dao:       dao,
	}
}

//StartMonitor creates a periodic job that watches the filesystem for Job Description files to process
func (j *JobRunner) StartMonitor() {
	j.log.Info("Starting Monitor")
	j.running = true
	ticker := time.NewTicker(time.Second * 10) //every 10 sec
	go func() {
		for {
			if !j.running {
				return
			}
			<-ticker.C
			if !j.jobLocked {
				j.log.Info("Running Job Detector")
				job := j.CheckForJobs()
				if job != nil {
					j.jobLocked = true
					go j.processJob(job)
				}
			} else {
				j.log.Info("Skipping run, Job is processing")
			}
		}
	}()
}

//checkForJobs scans folder for job metadata if there is at least one it is created and returned to be executed
func (j *JobRunner) CheckForJobs() *BGJob {
	job := j.FirstOpenJob()
	if job != nil {
		job.env = j
		return job
	}
	return nil
}

// Stop the periodic jobs checking
func (j *JobRunner) StopMonitor() {
	j.log.Info("Stop Job Monitor")
	j.running = false
}

func (j *JobRunner) processJob(job *BGJob) {
	newstatus := job.JobStatus
	switch job.JobType {
	case PdfToCbz:
		{
			newstatus = convertToPDF(job)
		}
	case ScanMeta:
		{
			//FIXME begin time
			newstatus = scanMetaData(job)
			//FIXME endTime
		}

	default:
		j.log.Errorf("Unsupported Job Type: %v", job.JobType)

	}

	j.UpdateJobStatus(job, newstatus)

	j.jobLocked = false
}

func (j *JobRunner) save(job *BGJob) {
	err := j.dao.Write(JOB_OPEN_BUCKET, job)
	if err != nil {
		panic(err)
	}
}

func (j *JobRunner) FirstOpenJob() *BGJob {
	r := new(BGJob)
	err := j.dao.Db.View(func(tx *bolt.Tx) error {
		t := tx.Bucket([]byte("jobs_open"))
		_, v := t.Cursor().First()
		if v == nil {
			return fmt.Errorf("No Jobs available")
		}

		dec := json.NewDecoder(bytes.NewReader(v))
		return dec.Decode(r)
	})
	if err != nil {
		//legal reason is not errors found
		return nil
	}
	return r
}

func (j *JobRunner) UpdateJobStatus(job *BGJob, newstatus int) {
	oldBucket := JOB_OPEN_BUCKET
	newBucket := JOB_OPEN_BUCKET
	switch job.JobStatus {
	case NotStarted:
		oldBucket = JOB_OPEN_BUCKET
	case Done:
		oldBucket = JOB_DONE_BUCKET
	case Error:
		oldBucket = JOB_ERROR_BUCKET
	default:
		panic(fmt.Errorf("Unknown Job Status %d", job.JobStatus))
	}
	job.JobStatus = newstatus
	switch job.JobStatus {
	case NotStarted:
		newBucket = JOB_OPEN_BUCKET
	case Done:
		newBucket = JOB_DONE_BUCKET
	case Error:
		newBucket = JOB_ERROR_BUCKET
	default:
		panic(fmt.Errorf("Unknown Job Status %d", job.JobStatus))
	}
	err := j.dao.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(oldBucket).Delete(job.IdBytes())
	})
	if err != nil {
		panic(err)
	}
	err2 := j.dao.Write(newBucket, job)
	if err2 != nil {
		panic(err)
	}
}

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


var bucketJobOpen = []byte("jobs_open")
var bucketJobDone = []byte("jobs_done")
var bucketJobError = []byte("jobs_error")

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
	env         *JobRunner //give access to the JobRunner and DAOs
}

//JOBS are defined by creating a name.job json in a special folder. jobs are processed one after the other by reading the directory where the jobs are and
//taking one job, reading the description and processing it. MEta is updated while processing ... only 1 job at a time

type JobRunner struct {
	running   bool
	jobLocked bool
	log       *logger.Logger
	dao       *dao.DAOHandler
}

//NewJobRunner Constructor
func NewJobRunner(dao *dao.DAOHandler) *JobRunner { //fixme give config instead of job path
	out, _ := os.Create("jobs.log")
	lg, _ := logger.NewLogger("GnolJob", 0, logger.InfoLevel, out)
	return &JobRunner{
		running:   false,
		jobLocked: false,
		log:       lg,
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
	newStatus := job.JobStatus
	switch job.JobType {
	case PdfToCbz:
		{
			newStatus = convertToPDF(job)
		}
	case ScanMeta:
		{
			//FIXME begin time
			newStatus = scanMetaData(job)
			//FIXME endTime
		}

	default:
		j.log.Errorf("Unsupported Job Type: %v", job.JobType)

	}

	j.UpdateJobStatus(job, newStatus)

	j.jobLocked = false
}

func (j *JobRunner) save(job *BGJob) {
	err := j.dao.Write(bucketJobOpen, job)
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
			return fmt.Errorf("no jobs available")
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

func (j *JobRunner) UpdateJobStatus(job *BGJob, newStatus int) {
	oldBucket := bucketJobOpen
	newBucket := bucketJobOpen
	switch job.JobStatus {
	case NotStarted:
		oldBucket = bucketJobOpen
	case Done:
		oldBucket = bucketJobDone
	case Error:
		oldBucket = bucketJobError
	default:
		panic(fmt.Errorf("unknown job status %d", job.JobStatus))
	}
	job.JobStatus = newStatus
	switch job.JobStatus {
	case NotStarted:
		newBucket = bucketJobOpen
	case Done:
		newBucket = bucketJobDone
	case Error:
		newBucket = bucketJobError
	default:
		panic(fmt.Errorf("unknown job status %d", job.JobStatus))
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

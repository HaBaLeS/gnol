package jobs

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/HaBaLeS/go-logger"
	"os"
	"path"
	"time"
)


var bucketJobOpen = []byte("jobs_open")
var bucketJobDone = []byte("jobs_done")
var bucketJobError = []byte("jobs_error")

const (
	PdfToCbz = iota
	ScanMeta
	DownloadUrl
)

const (
	NotStarted = iota
	Error
	Done
)


//JOBS are defined by creating a name.job json in a special folder. jobs are processed one after the other by reading the directory where the jobs are and
//taking one job, reading the description and processing it. MEta is updated while processing ... only 1 job at a time
type JobRunner struct {
	running   bool
	jobLocked bool
	log       *logger.Logger
	cfg		  *util.ToolConfig
	dao       *storage.DAO
}

//NewJobRunner Constructor
func NewJobRunner(dao *storage.DAO, cfg *util.ToolConfig) *JobRunner { //fixme give config instead of job path
	out, _ := os.Create(path.Join(cfg.TempDirectory,"jobs.log"))
	lg, _ := logger.NewLogger("GnolJob", 0, logger.InfoLevel, out)
	return &JobRunner{
		running:   false,
		jobLocked: false,
		log:       lg,
		dao: dao,
		cfg: cfg,
	}
}

//StartMonitor creates a periodic job that watches the filesystem for Job Description files to process
func (j *JobRunner) StartMonitor() {
	j.log.Info("Starting Monitor")
	j.running = true
	ticker := time.NewTicker(time.Second * 2) //every 10 sec
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
func (j *JobRunner) CheckForJobs() *storage.GnolJob {
	job := j.FirstOpenJob()
	if job != nil {
		return job
	}
	return nil
}

// Stop the periodic jobs checking
func (j *JobRunner) StopMonitor() {
	j.log.Info("Stop Job Monitor")
	j.running = false
}

func (j *JobRunner) processJob(job *storage.GnolJob) {
	var jobError error
	switch job.JobType {
	case PdfToCbz:
		{
			jobError = j.convertToPDF(job)
		}
	case ScanMeta:
		{
			jobError = j.scanMetaData(job)
		}
	case DownloadUrl:
		{
			jobError = j.downloadFromUrl(job)
		}

	default:
		j.log.Errorf("Unsupported Job Type: %v", job.JobType)
	}

	if jobError != nil{
		fmt.Printf("Error in job: %s\n", jobError)
		j.dao.UpdatJobStatus(job, Error)
	} else {
		//j.bs.Delete(job) //FIXME this makes the job run forever!! not good!!!
		j.dao.UpdatJobStatus(job, Done)
	}
	j.jobLocked = false
}

func (j *JobRunner) save(job *storage.GnolJob) {
	err := j.dao.CreateJob(job.JobType, job.UserID, job.Data)
	if err != nil {
		panic(err)
	}
}

func (j *JobRunner) FirstOpenJob() *storage.GnolJob {
	job := j.dao.GetOldestOpenJob()
	if job == nil {
		return nil
	}

	return job
}

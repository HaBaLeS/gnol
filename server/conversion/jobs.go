package conversion

import (
	"encoding/json"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/HaBaLeS/go-logger"
	"github.com/rs/xid"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	PdfToCbz = iota
	ScanMeta
	CleanCache
)

const (
	NotStarted = iota
	Error
	Done
)

var jobDir string

type BGJob struct {
	JobType     int    //Convert to CBR, Scan for Metadata, Scrape Meta, SortFolder, Clean Cache etc....
	JobStatus   int    //What's the status done, not started, error
	DisplayName string //Name the Job
	Duration    string //how long did it take
	MetaFile    string //reference to myself
	InputFile   string
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
func NewJobRunner(jbDir string, dao *dao.DAOHandler) *JobRunner { //fixme give config instead of job path
	jobDir = jbDir
	jgf := path.Join(jobDir, "jobs.log")
	out, _ := os.Create(jgf)
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
	files, err := ioutil.ReadDir(jobDir)
	if err != nil {
		//fixme handle better!!
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".job" {
			j.log.Debugf("skip: %s -> %s\n", f.Name(), filepath.Ext(f.Name()))
			continue //skip unknown files
		}
		dscf, oe := os.Open(path.Join(jobDir, f.Name()))
		if oe != nil {
			j.log.Errorf("Cannot Open: %s", path.Join(jobDir, f.Name()))
		}
		dec := json.NewDecoder(dscf)
		job := &BGJob{}
		de := dec.Decode(job)
		if de != nil {
			j.log.Errorf("Could not Decode Job: '%s'. Error %v", f.Name(), de)
			continue
		}
		if job.JobStatus == NotStarted {
			return job
		}
	}
	return nil
}

// Stop the periodic jobs checking
func (j *JobRunner) StopMonitor() {
	j.log.Info("Stop Job Monitor")
	j.running = false
}

func (j *JobRunner) processJob(job *BGJob) {
	switch job.JobType {
	case PdfToCbz:
		{
			convertToPDF(job)
		}
	case ScanMeta:
		{
			//FIXME begin time
			scanMetaData(job, j.dao)
			//FIXME endTime
		}

	default:
		j.log.Errorf("Unsupported Job Type: %v", job.JobType)
	}

	//FIXME learn to recover properly
	job.save()
	j.jobLocked = false
}

func (j *BGJob) save() {
	outfile := j.MetaFile
	if outfile == "" {
		outfile = filepath.Join(jobDir, xid.New().String()+".job")
		j.MetaFile = outfile
	}

	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jenc := json.NewEncoder(f)
	ence := jenc.Encode(j)
	if ence != nil {
		panic(ence)
	}
}

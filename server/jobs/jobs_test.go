package jobs

import (
	"testing"
)

var jr *JobRunner

func TestJobRunner_StartMonitor(t *testing.T) {
	jr.StartMonitor()
}

func TestJobRunner_StopMonitor(t *testing.T) {
	jr.StopMonitor()
	if jr.running {
		t.Errorf("Expected jr.running to be false")
	}
}

func TestJobRunner_CheckForJobs(t *testing.T) {
	bgjob := jr.CheckForJobs()

	if bgjob == nil {
		t.Error("Did not receive BGJob File")
	}

}

func TestMain(m *testing.M) {
	jr = NewJobRunner(nil)
	m.Run()
}

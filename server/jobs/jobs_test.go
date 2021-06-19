package jobs

import (
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/stretchr/testify/assert"
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

var config *util.ToolConfig

func TestJobRunner_CheckForJobs(t *testing.T) {
	assert.FailNow(t,"Not implemented")
	bgjob := jr.CheckForJobs()

	if bgjob == nil {
		t.Error("Did not receive BGJob File")
	}

}

func TestMain(m *testing.M) {
	config = &util.ToolConfig{

	}
	jr = NewJobRunner(nil, config)
	m.Run()
}

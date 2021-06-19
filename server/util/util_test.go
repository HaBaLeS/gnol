package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)



func TestMain(m *testing.M) {
	m.Run() //Dont forget to start the Test ;-)
}

func TestBla(t *testing.T){
	assert.Equal(t,1,2,"blafasel")
}


func TestReadConfig(t *testing.T) {

	_, err := ReadConfig("testdata/doesnotexist")
	assert.Errorf(t,err,"Reading non exisisting file should produce a error")
	cfg, _ := ReadConfig("testdata/test.cfg")
	assert.Equal(t, cfg.ListenPort,666,"Failure reading int form config file")
}

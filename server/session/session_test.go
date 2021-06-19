package session

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


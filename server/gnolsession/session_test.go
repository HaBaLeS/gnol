package gnolsession

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run() //Dont forget to start the Test ;-)
}

func TestBla(t *testing.T) {
	assert.Equal(t, 1, 2, "blafasel")
}

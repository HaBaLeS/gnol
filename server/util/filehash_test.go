package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashFile(t *testing.T) {

	h, err := HashFile("testdata/shasum")
	assert.NoError(t, err)
	assert.Equal(t, "29ad0edde0616bf92a58e81e7ce20a685ba1d0d0105d6beeb725b09218502f47", h)
}

package persistence_test

import (
	"github.com/HaBaLeS/gnol/cmd/leech-tool/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
)

type PersistenceSuite struct {
	suite.Suite
	container *persistence.DataContainer
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPersistenceSuite(t *testing.T) {
	ts := &PersistenceSuite{
		container: persistence.LoadDataInMemory("test.json"),
	}
	suite.Run(t, ts)
	ts.container.Commit()
}

func (s *PersistenceSuite) TestLeechDBForHost() {
	ldb := s.container.LeechDBForHost("habales.de")
	assert.NotNil(s.T(), ldb)
	assert.Equal(s.T(), ldb.Host, "habales.de")
}

func (s *PersistenceSuite) TestLeechDB_UrlData() {
	ldb := s.container.LeechDBForHost("habales.de")
	url, err := url.Parse("heise.de")
	assert.NoError(s.T(), err)
	data := ldb.UrlData(url)
	assert.NotEmpty(s.T(), data.Created)
	assert.NotEmpty(s.T(), data.DataUrl)
}

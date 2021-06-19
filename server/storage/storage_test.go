package storage

import (
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

var us *UserStore

func TestMain(m *testing.M) {
	tc := &util.ToolConfig{
		Database: "testdata/test.db",
	}
	bs := NewBoltStorage(tc)
	us = newUserStore(bs)
	m.Run() //Dont forget to start the Test ;-)

}

func TestUserStore_CreateUser(t *testing.T) {
	u := us.CreateUser("hansd","passwd")
	assert.Equal(t,"hands", u.Name, "Name not set")
	assert.NotEmpty(t, u.Id)
	assert.NotEmpty(t, u.PwdHash)
	assert.NotEmpty(t, u.Salt)

	us.bs.Close()
}

func TestBaseEntity_Id(t *testing.T) {
	be := CreateBaseEntity([]byte("testStorage"))
	id := be.Id
	if id == "" {
		t.Error("ID not Generated!")
	}
}

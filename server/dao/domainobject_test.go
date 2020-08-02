package dao

import (
	"path/filepath"
	"testing"
)

var dao *DaoService

func TestMain(m *testing.M) {
	absPath, _ := filepath.Abs("testdata/")
	dao = NewDaoService(absPath)
	m.Run() //Dont forget to start the Test ;-)
}

func TestCreate(t *testing.T) {
	u := CreateUser()

	if u.BaseEntity.id == "" {
		t.Error("BaseEntity not set up correct")
	}
}

func TestBaseEntity_Id(t *testing.T) {
	be := createBaseEntity("user")
	id := be.Id()
	if id == "" {
		t.Error("ID not Generated!")
	}
}

func TestBaseEntity_Save(t *testing.T) {

}

func TestLoad(t *testing.T) {

	utf := User{}
	dao.Load(utf, "1234")

	w := "gnol1"
	g := utf.Name

	if w != g {
		t.Errorf("Load: Want %s as name - got %s", w, g)
	}
}

package dao

import (
	"testing"
)


func TestMain(m *testing.M) {
	m.Run() //Dont forget to start the Test ;-)
}

func TestBaseEntity_Id(t *testing.T) {
	be := CreateBaseEntity()
	id := be.Id
	if id == "" {
		t.Error("ID not Generated!")
	}
}

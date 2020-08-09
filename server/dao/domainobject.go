package dao
/*
type DaoService struct {
	dataPath string
}

func NewDaoService(path string) *DaoService {
	return &DaoService{
		dataPath: path,
	}
}*/

type BaseEntity struct {
	Id string
}

type Entity interface {
	IdBytes() []byte
}

func (b *BaseEntity) IdBytes() []byte {
	return []byte(b.Id)
}

package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/boltdb/bolt"
	"github.com/rs/xid"
	"time"
)



type BoltStorage struct {
	db        *bolt.DB
	Comic     *ComicStorage
}

func NewBoltStorage(cfg *util.ToolConfig) *BoltStorage{
	bs := &BoltStorage{}
	bs.Init(cfg)
	bs.Comic = newComicStore(bs, cfg)
	return bs
}

type StorageFunc func(tx *bolt.Tx) error

func (ms *BoltStorage) Init(config *util.ToolConfig){
	db, err := bolt.Open(config.Database, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	ms.db = db

	err = ms.WriteRaw(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("meta"))
		_, err = tx.CreateBucketIfNotExists([]byte("jobs_open"))
		_, err = tx.CreateBucketIfNotExists([]byte("jobs_error"))
		_, err = tx.CreateBucketIfNotExists([]byte("jobs_done"))
		return err
	})
	if err != nil {
		panic(err)
	}
}

/*
func (ms *BoltStorage) Load( id string, into BaseEntity ) error {
	//find Bucket by ID

	//find item and load
	panic("not implemented")
	//deserialize item -> into
}*/

func (ms *BoltStorage) Write(from Entity) error {
	bucket := from.Bucket()
	return ms.write(bucket, from)
}

func (ms *BoltStorage) write(bucket []byte, data Entity) error {
	return ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		buf := new(bytes.Buffer)
		e := json.NewEncoder(buf)
		ece := e.Encode(data)
		if ece != nil {
			return ece
		}
		b.Put(data.IdBytes(), buf.Bytes())
		return nil
	})
}

func (ms *BoltStorage) Delete(entity Entity) error{
	return ms.db.Update(func(tx *bolt.Tx) error {
		 bck := tx.Bucket(entity.Bucket())
		 return bck.Delete(entity.IdBytes())
	})
}

func (ms *BoltStorage) ReadRaw(fn StorageFunc) error {
	return ms.db.View(fn)
}
func (ms *BoltStorage) WriteRaw(fn StorageFunc) error {
	return ms.db.Update(fn)
}

func (ms *BoltStorage) Close() {
	ms.db.Close()
}



func loadFromJson(i interface{}, v []byte) error {
	d := json.NewDecoder(bytes.NewReader(v))
	return d.Decode(i)
}


type Entity interface {
	IdBytes() []byte
	Bucket() []byte
}

type BaseEntity struct {
	Id string
}

func (b *BaseEntity) ChangeBucket(bucket []byte) {
	b.Id = fmt.Sprintf("%s|%s", b.ident(), bucket)
}

func (b *BaseEntity) IdBytes() []byte {
	return []byte(b.Id)
}

func (b *BaseEntity) ChangeId(ident string) {
	bu := b.Bucket()
	b.Id = fmt.Sprintf("%s|%s", ident, bu)
}

func (b *BaseEntity) Bucket() []byte {
	sl:= bytes.Split([]byte(b.Id),[]byte("|"))
	return sl[len(sl)-1]
}

func (b *BaseEntity) ident() []byte {
	sl:= bytes.Split([]byte(b.Id),[]byte("|"))
	return bytes.Join(sl[:len(sl)-1],[]byte("|"))
}

//CreateBaseEntity creates an entity, generates an Id and appends the bucket and
func CreateBaseEntity(bucket []byte) *BaseEntity {
	return &BaseEntity{
		Id: fmt.Sprintf("%s|%s",xid.New().String(),bucket),
	}
}


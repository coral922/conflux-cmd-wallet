package storage

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/tidwall/gjson"
	"os"
	"time"
)

type Storage struct {
	*bolt.DB
	dbFile string
}

func NewStorage(dbFile string) (*Storage, error) {
	db, err := bolt.Open(dbFile, 0666, &bolt.Options{
		Timeout: time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db, dbFile: dbFile}, nil
}

func (s *Storage) Destroy() error {
	_ = s.DB.Close()
	_ = os.Remove(s.dbFile + ".lock")
	return os.Remove(s.dbFile)
}

func (s *Storage) Get(bucket string, key string) (string, error) {
	var res string
	var err error
	err = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(key))
		res = string(v)
		return nil
	})
	return res, err
}

func (s *Storage) Keys(bucket string, filter ...Filter) ([]string, error) {
	res := make([]string, 0)
	var err error
	err = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var found = true
			for _, f := range filter {
				if !f(v) {
					found = false
					break
				}
			}
			if found {
				res = append(res, string(k))
			}
			return nil
		})
	})
	return res, err
}

func (s *Storage) HasBucket(bucket string) (bool, error) {
	var res bool
	var err error
	err = s.View(func(tx *bolt.Tx) error {
		res = tx.Bucket([]byte(bucket)) != nil
		return nil
	})
	return res, err
}

func (s *Storage) First(bucket string, filter ...Filter) (string, error) {
	var res string
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var found = true
			for _, f := range filter {
				if !f(v) {
					found = false
					break
				}
			}
			if found {
				res = string(v)
				break
			}
		}
		return nil
	})
	return res, err
}

func (s *Storage) All(bucket string, filter ...Filter) ([]string, error) {
	res := make([]string, 0)
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var found = true
			for _, f := range filter {
				if !f(v) {
					found = false
					break
				}
			}
			if found {
				res = append(res, string(v))
			}
		}
		return nil
	})
	return res, err
}

func (s *Storage) GetStruct(bucket string, key string, pointer interface{}) error {
	v, err := s.Get(bucket, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(v), pointer)
}

func (s *Storage) Set(bucket, key, value string) error {
	var err error
	err = s.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
	return err
}

func (s *Storage) Del(bucket, key string) error {
	var err error
	err = s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
	return err
}

func (s *Storage) SetStruct(bucket string, key string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.Set(bucket, key, string(j))
}

type Filter func(data []byte) bool

func JsonStrAttrFilter(attr, v string) Filter {
	return func(data []byte) bool {
		if !gjson.ValidBytes(data) {
			return false
		}
		return gjson.GetBytes(data, attr).String() == v
	}
}

func JsonBoolAttrFilter(attr string, v bool) Filter {
	return func(data []byte) bool {
		if !gjson.ValidBytes(data) {
			return false
		}
		return gjson.GetBytes(data, attr).Bool() == v
	}
}

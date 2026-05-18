package x

import (
	"bytes"
	"errors"

	"go.etcd.io/bbolt"
)

var (
	ErrBboltBucketNotFound = errors.New("bbolt bucket not found")
	ErrBboltKeyNotFound    = errors.New("bbolt key not found")
)

func BboltOpen(path string) (*bbolt.DB, error) {
	if path == "" {
		return nil, ErrInvalidArgs
	}
	if db, err := bbolt.Open(path, 0600, nil); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func BboltCreateBuckets(db *bbolt.DB, buckets ...[]byte) error {
	if db == nil || len(buckets) == 0 {
		return ErrInvalidArgs
	}
	return db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if len(b) > 0 {
				if _, err := tx.CreateBucketIfNotExists(b); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func BboltDeleteBuckets(db *bbolt.DB, buckets ...[]byte) error {
	if db == nil || len(buckets) == 0 {
		return ErrInvalidArgs
	}
	return db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if len(b) > 0 {
				if err := tx.DeleteBucket(b); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func BboltView(db *bbolt.DB, bucket []byte, fn func(*bbolt.Tx, *bbolt.Bucket) error) error {
	if db == nil || len(bucket) == 0 || fn == nil {
		return ErrInvalidArgs
	}
	return db.View(func(tx *bbolt.Tx) error {
		if b := tx.Bucket(bucket); b == nil {
			return ErrBboltBucketNotFound
		} else {
			return fn(tx, b)
		}
	})
}

func BboltUpdate(db *bbolt.DB, bucket []byte, fn func(*bbolt.Tx, *bbolt.Bucket) error) error {
	if db == nil || len(bucket) == 0 || fn == nil {
		return ErrInvalidArgs
	}
	return db.Update(func(tx *bbolt.Tx) error {
		if b := tx.Bucket(bucket); b == nil {
			return ErrBboltBucketNotFound
		} else {
			return fn(tx, b)
		}
	})
}

func BboltGet(db *bbolt.DB, bucket, key []byte) ([]byte, error) {
	if db == nil || len(bucket) == 0 || len(key) == 0 {
		return nil, ErrInvalidArgs
	}
	var ret []byte
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBboltBucketNotFound
		}
		val := b.Get(key)
		if val == nil {
			return ErrBboltKeyNotFound
		}
		ret = make([]byte, len(val))
		copy(ret, val)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func BboltPut(db *bbolt.DB, bucket, key, value []byte) error {
	if db == nil || len(bucket) == 0 || len(key) == 0 || len(value) == 0 {
		return ErrInvalidArgs
	}
	return db.Update(func(tx *bbolt.Tx) error {
		if b := tx.Bucket(bucket); b == nil {
			return ErrBboltBucketNotFound
		} else {
			return b.Put(key, value)
		}
	})
}

func BboltDelete(db *bbolt.DB, bucket, key []byte) error {
	if db == nil || len(bucket) == 0 || len(key) == 0 {
		return ErrInvalidArgs
	}
	return db.Update(func(tx *bbolt.Tx) error {
		if b := tx.Bucket(bucket); b == nil {
			return ErrBboltBucketNotFound
		} else {
			return b.Delete(key)
		}
	})
}

// prefix为nil或空表示迭代整个bucket
func BboltIter(db *bbolt.DB, bucket, prefix []byte, fn func([]byte, []byte) error) error {
	if db == nil || len(bucket) == 0 || fn == nil {
		return ErrInvalidArgs
	}
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBboltBucketNotFound
		}
		c := b.Cursor()
		if len(prefix) == 0 {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if err := fn(k, v); err != nil {
					return err
				}
			}
		} else {
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				if err := fn(k, v); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// prefix为nil或空表示统计整个bucket键值对数量
func BboltCount(db *bbolt.DB, bucket, prefix []byte) (int, error) {
	if db == nil || len(bucket) == 0 {
		return 0, ErrInvalidArgs
	}
	ret := 0
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBboltBucketNotFound
		}
		if len(prefix) == 0 {
			ret = b.Stats().KeyN
		} else {
			c := b.Cursor()
			for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
				ret++
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return ret, nil
}

package x

import (
	"errors"

	"go.etcd.io/bbolt"
)

var ErrBboltQueueEmpty = errors.New("bbolt queue is empty")

type BboltQueue struct {
	db     *bbolt.DB
	name   string
	bucket []byte
}

func NewBboltQueue(db *bbolt.DB, name string) *BboltQueue {
	if db == nil || name == "" {
		return nil
	}
	bucket := []byte(name)
	if err := BboltCreateBuckets(db, bucket); err == nil {
		return &BboltQueue{
			db:     db,
			name:   name,
			bucket: bucket,
		}
	}
	return nil
}

func (q *BboltQueue) PushFirst(data []byte) error {
	if len(data) == 0 {
		return ErrInvalidArgs
	}
	return BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		var id int64
		if key, _ := b.Cursor().First(); key != nil {
			id = Btoi64[int64](key)
		}
		return b.Put(I64tob(id-1), data)
	})
}

func (q *BboltQueue) PushLast(data []byte) error {
	if len(data) == 0 {
		return ErrInvalidArgs
	}
	return BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		if id, err := b.NextSequence(); err != nil {
			return err
		} else {
			return b.Put(I64tob(id), data)
		}
	})
}

func (q *BboltQueue) PopFirst() ([]byte, error) {
	var ret []byte
	err := BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		c := b.Cursor()
		if key, val := c.First(); key == nil {
			return ErrBboltQueueEmpty
		} else {
			ret = make([]byte, len(val))
			copy(ret, val)
			return c.Delete()
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *BboltQueue) PopLast() ([]byte, error) {
	var ret []byte
	err := BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		c := b.Cursor()
		if key, val := c.Last(); key == nil {
			return ErrBboltQueueEmpty
		} else {
			ret = make([]byte, len(val))
			copy(ret, val)
			return c.Delete()
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *BboltQueue) GetFirst() ([]byte, error) {
	var ret []byte
	err := BboltView(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		if key, val := b.Cursor().First(); key == nil {
			return ErrBboltQueueEmpty
		} else {
			ret = make([]byte, len(val))
			copy(ret, val)
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *BboltQueue) GetLast() ([]byte, error) {
	var ret []byte
	err := BboltView(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		if key, val := b.Cursor().Last(); key == nil {
			return ErrBboltQueueEmpty
		} else {
			ret = make([]byte, len(val))
			copy(ret, val)
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *BboltQueue) Len() int {
	if l, err := BboltCount(q.db, q.bucket, nil); err == nil {
		return l
	}
	return 0
}

func (q *BboltQueue) IsEmpty() bool {
	ret := true
	BboltView(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		if k, _ := b.Cursor().First(); k != nil {
			ret = false
		}
		return nil
	})
	return ret
}

func (q *BboltQueue) Clear() error {
	return BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		if err := tx.DeleteBucket(q.bucket); err != nil {
			return err
		} else {
			_, err = tx.CreateBucketIfNotExists(q.bucket)
			return err
		}
	})
}

func (q *BboltQueue) BatchPushFirst(data [][]byte) error {
	if len(data) == 0 {
		return ErrInvalidArgs
	}
	return BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		var id int64
		if key, _ := b.Cursor().First(); key != nil {
			id = Btoi64[int64](key)
		}
		for i := range data {
			if len(data[i]) != 0 {
				if err := b.Put(I64tob(id-int64(i)-1), data[i]); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (q *BboltQueue) BatchPushLast(data [][]byte) error {
	if len(data) == 0 {
		return ErrInvalidArgs
	}
	return BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		for i := range data {
			if len(data[i]) != 0 {
				if id, err := b.NextSequence(); err != nil {
					return err
				} else {
					return b.Put(I64tob(id), data[i])
				}
			}
		}
		return nil
	})
}

func (q *BboltQueue) BatchPopFirst(n int) ([][]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidArgs
	}
	ret := make([][]byte, 0, n)
	err := BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		c := b.Cursor()
		i := 0
		for key, val := c.First(); key != nil && i < n; key, val = c.Next() {
			data := make([]byte, len(val))
			copy(data, val)
			ret = append(ret, data)
			i++
		}
		if len(ret) == 0 {
			return ErrBboltQueueEmpty
		}
		for range len(ret) {
			if key, _ := c.First(); key != nil {
				if err := c.Delete(); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *BboltQueue) BatchPopLast(n int) ([][]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidArgs
	}
	ret := make([][]byte, 0, n)
	err := BboltUpdate(q.db, q.bucket, func(tx *bbolt.Tx, b *bbolt.Bucket) error {
		c := b.Cursor()
		i := 0
		for key, val := c.Last(); key != nil && i < n; key, val = c.Prev() {
			data := make([]byte, len(val))
			copy(data, val)
			ret = append(ret, data)
			i++
		}
		if len(ret) == 0 {
			return ErrBboltQueueEmpty
		}
		for range len(ret) {
			if key, _ := c.Last(); key != nil {
				if err := c.Delete(); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

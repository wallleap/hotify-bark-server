package gotifycompat

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Store persists messages and the client-token material for the
// gotify-compatible monitoring interface.
type Store interface {
	// Add persists a message, assigns a monotonic ID and returns it.
	Add(m *Message) (uint64, error)
	// Recent returns up to limit messages with ID < since, ordered by ID
	// descending (newest first). since==0 disables the filter.
	Recent(limit int, since uint64) ([]Message, error)

	TokenHash() ([]byte, error)
	SetTokenHash(h []byte) error
	AutoToken() (string, error)
	SetAutoToken(t string) error

	Close() error
}

// bbolt-backed store persisted under <data>/gotify.db.
type bboltStore struct {
	db  *bolt.DB
	max int
}

const (
	bucketMessages = "messages"
	bucketMeta     = "meta"
)

var (
	keyCount     = []byte("count")
	keyTokenHash = []byte("tokenHash")
	keyAutoToken = []byte("autoToken")
)

func openBboltStore(path string, max int) (*bboltStore, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketMessages)); err != nil {
			return err
		}
		_, err := tx.CreateBucketIfNotExists([]byte(bucketMeta))
		return err
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &bboltStore{db: db, max: max}, nil
}

func (s *bboltStore) Close() error {
	return s.db.Close()
}

func uint64Bytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func bytesUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func (s *bboltStore) Add(m *Message) (uint64, error) {
	var id uint64
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketMessages))
		mb := tx.Bucket([]byte(bucketMeta))

		seq, err := b.NextSequence()
		if err != nil {
			return err
		}
		id = seq
		m.ID = id

		v, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if err := b.Put(uint64Bytes(id), v); err != nil {
			return err
		}

		count := metaUint64(mb, keyCount) + 1
		if count > uint64(s.max) {
			c := b.Cursor()
			if k, _ := c.First(); k != nil {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
			count--
		}
		return mb.Put(keyCount, uint64Bytes(count))
	})
	return id, err
}

func (s *bboltStore) Recent(limit int, since uint64) ([]Message, error) {
	if limit <= 0 {
		return []Message{}, nil
	}
	out := make([]Message, 0, limit)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketMessages))
		c := b.Cursor()
		// Position just past the last key, then step back onto the newest.
		_, _ = c.Seek(uint64Bytes(^uint64(0)))
		k, _ := c.Prev()
		for len(out) < limit && k != nil {
			id := bytesUint64(k)
			if since != 0 && id >= since {
				k, _ = c.Prev()
				continue
			}
			var m Message
			if err := json.Unmarshal(b.Get(k), &m); err != nil {
				return err
			}
			m.ID = id
			out = append(out, m)
			k, _ = c.Prev()
		}
		return nil
	})
	return out, err
}

func (s *bboltStore) TokenHash() ([]byte, error) {
	var out []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucketMeta)).Get(keyTokenHash)
		out = append(out, v...)
		return nil
	})
	return out, err
}

func (s *bboltStore) SetTokenHash(h []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketMeta)).Put(keyTokenHash, h)
	})
}

func (s *bboltStore) AutoToken() (string, error) {
	var out string
	err := s.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucketMeta)).Get(keyAutoToken)
		out = string(v)
		return nil
	})
	return out, err
}

func (s *bboltStore) SetAutoToken(t string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketMeta)).Put(keyAutoToken, []byte(t))
	})
}

func metaUint64(b *bolt.Bucket, key []byte) uint64 {
	v := b.Get(key)
	if len(v) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(v)
}

// memoryStore is the degraded (non-persistent) fallback used when the data
// directory is unavailable. IDs remain monotonic for the lifetime of the
// process; the bridge tolerates an ID reset after a restart.
type memoryStore struct {
	mu    sync.Mutex
	seq   uint64
	max   int
	msgs  map[uint64]Message
	order []uint64
}

func newMemoryStore(max int) *memoryStore {
	return &memoryStore{max: max, msgs: make(map[uint64]Message)}
}

func (s *memoryStore) Close() error { return nil }

func (s *memoryStore) Add(m *Message) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	id := s.seq
	m.ID = id
	s.msgs[id] = *m
	s.order = append(s.order, id)
	if len(s.order) > s.max {
		old := s.order[0]
		s.order = s.order[1:]
		delete(s.msgs, old)
	}
	return id, nil
}

func (s *memoryStore) Recent(limit int, since uint64) ([]Message, error) {
	if limit <= 0 {
		return []Message{}, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Message, 0, limit)
	for i := len(s.order) - 1; i >= 0 && len(out) < limit; i-- {
		id := s.order[i]
		if since != 0 && id >= since {
			continue
		}
		out = append(out, s.msgs[id])
	}
	return out, nil
}

func (s *memoryStore) TokenHash() ([]byte, error) {
	return nil, nil
}

func (s *memoryStore) SetTokenHash(h []byte) error {
	return fmt.Errorf("in-memory store does not persist tokens")
}

func (s *memoryStore) AutoToken() (string, error) {
	return "", nil
}

func (s *memoryStore) SetAutoToken(t string) error {
	return fmt.Errorf("in-memory store does not persist tokens")
}

func openStore(dataDir string, max int) (Store, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("empty data dir")
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}
	return openBboltStore(filepath.Join(dataDir, "gotify.db"), max)
}

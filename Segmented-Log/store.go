package segmentedlog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var (
	errNoStoreSpaceLeft = errors.New("no store space left")
	errNoIndexSpaceLeft = errors.New("no index space left")
)

// store defines a storage abstraction for the log
// log is append only file
type store struct {
	file    *os.File
	size    uint64
	maxSize uint64
}

// newStore returns a new storage
func newStore(file string, cfg *Config) (*store, error) {
	st, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &store{
		file:    f,
		size:    uint64(st.Size()),
		maxSize: cfg.Segment.MaxStoreSizeBytes,
	}, nil
}

// read takes an offset in a file and returns a record
func (s *store) read(offset uint64) ([]byte, error) {
	// read the first 8 bytes to determine the size of the record
	b := make([]byte, 8)
	_, err := s.file.ReadAt(b, int64(offset))
	if err != nil {
		return nil, err
	}

	var size uint64
	err = binary.Read(bytes.NewReader(b), binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	b = make([]byte, size)
	_, err = s.file.ReadAt(b, int64(offset)+8)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// write append the record to the log and return
func (s *store) write(data []byte) (uint64, error) {
	if s.size+uint64(len(data)+8) > s.maxSize {
		return 0, errNoStoreSpaceLeft
	}

	b := make([]byte, 8+len(data))
	binary.BigEndian.PutUint64(b[0:8], uint64(len(data)))
	copy(b[8:], data)

	n, err := s.file.Write(b)
	if err != nil {
		return 0, err
	}

	if n != len(data)+8 {
		return 0, fmt.Errorf("can't write all data")
	}

	err = s.file.Sync()
	if err != nil {
		return 0, err
	}

	offset := s.size
	s.size += uint64(n)

	return offset, nil
}

func (s *store) close() error {
	return s.file.Close()
}

func (s *store) remove() error {
	return os.Remove(s.file.Name())
}

package datareader

import (
	"errors"
	"io"
	"os"
	"syscall"
)

// DataReader .
type DataReader interface {
	Read() bool
	Close()
	Data() []byte
	Err() error
}

// MMapReader is a concreate type of DataReader which implements MMap
// ChunkSize is a size
type MMapReader struct {
	FileDescriptor int
	ChunkSize      int
	data           []byte
	err            error
	backBuffer     []byte
	offset         int64
	file           *os.File
	fileSize       int64
	index          int
}

// NewMMapReader creates an instance of MMapReader which utilize syscall.Mmap
// It takes a chunkSize (option) which read portion or all of the data
// If chunkSize is greater than the size of pageSize, it must exact in multiple of page size
// chunkSize will be defaulted to pageSize if not being provided
// most common pageSize is 4KB
func NewMMapReader(filename string, chunkSize ...int) (DataReader, error) {
	pageSize := syscall.Getpagesize()
	size := pageSize
	if len(chunkSize) > 0 {
		size = chunkSize[0]
	}

	if size > pageSize && size%pageSize != 0 {
		return nil, errors.New("chunkSize must be a multiple of the page size if buffersize is greater than page size")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &MMapReader{
		FileDescriptor: int(file.Fd()),
		ChunkSize:      size,
		file:           file,
		fileSize:       fi.Size(),
	}, nil
}

// Read .
func (mr *MMapReader) Read() bool {
	if mr.offset >= mr.fileSize {
		// EOF
		mr.data = nil
		mr.err = io.EOF
		return false
	}
	if mr.fileSize-mr.offset < int64(mr.ChunkSize) {
		lastChunk := mr.fileSize - mr.offset
		mr.ChunkSize = int(lastChunk)
	}
	data, err := syscall.Mmap(mr.FileDescriptor, mr.offset, mr.ChunkSize, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		mr.data = nil
		mr.err = err
		return false
	}
	// recount the offset
	mr.offset += int64(mr.ChunkSize)

	// prepend data with previous backBuffer if any
	data = append(mr.backBuffer, data...)
	mr.backBuffer = nil

	if mr.offset >= mr.fileSize {
		// the last chunk
		mr.data = data

		return true
	}

	// backCounter to make sure data ends with line feed
	backCounter := mr.countBackBuffer(data)

	// put into backBuffer for next read
	mr.backBuffer = data[len(data)-backCounter:]

	// truncate data to make sure data ends with line feed
	data = data[:len(data)-backCounter]
	mr.data = data

	mr.index++
	return true
}

// Close .
func (mr *MMapReader) Close() {
	mr.file.Close()
}

// Data returns data
func (mr *MMapReader) Data() []byte {
	return mr.data
}

// Err returns error
func (mr *MMapReader) Err() error {
	return mr.err
}

// countBackBuffer counts back the index (if necessary)
// until it finds the line feed (\n)
func (mr *MMapReader) countBackBuffer(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	var back int
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '\n' {
			return back
		}
		back++
	}
	return back
}

package utils

import (
	"io"
	"sync"
)

// SynchronizedWriter provides WriteCloser interface with mutexed operations.
type SynchronizedWriter struct {
	writer io.WriteCloser
	lock   *sync.Mutex
}

// Write writes content to the writer mutually exclusive.
func (sw *SynchronizedWriter) Write(content []byte) (n int, err error) {
	sw.lock.Lock()
	defer sw.lock.Unlock()

	return sw.writer.Write(content)
}

// Close closes content to the writer mutually exclusive.
func (sw *SynchronizedWriter) Close() (err error) {
	sw.lock.Lock()
	defer sw.lock.Unlock()

	return sw.writer.Close()
}

// NewSynchronizedWriter makes writer synchronized.
func NewSynchronizedWriter(writer io.WriteCloser) (sw *SynchronizedWriter) {
	sw = new(SynchronizedWriter)
	sw.writer = writer
	sw.lock = new(sync.Mutex)

	return
}

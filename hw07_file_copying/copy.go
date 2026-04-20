package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var buffer []byte
	if limit > 0 {
		buffer = make([]byte, limit)
	} else {
		buffer = make([]byte, 1024)
	}

	fileFrom, openError := os.Open(fromPath)
	if openError != nil {
		return openError
	}
	defer fileFrom.Close()

	stat, statError := os.Stat(fromPath)
	if statError != nil {
		return ErrUnsupportedFile
	}

	fileSize := stat.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	remainingSize := fileSize - offset
	if limit < remainingSize && limit != 0 {
		remainingSize = limit
	}

	fileTo, createError := os.Create(toPath)
	if createError != nil {
		return createError
	}
	defer fileTo.Close()

	var currentPosition int64
	copiedBytes := 0

	for currentPosition < remainingSize {
		n, readError := fileFrom.ReadAt(buffer, offset+currentPosition)
		currentPosition += int64(n)
		if readError != nil && !errors.Is(readError, io.EOF) {
			return readError
		}

		writtenBytes, writeError := fileTo.Write(buffer[:n])
		copiedBytes += writtenBytes

		percentage := float64(copiedBytes) / float64(remainingSize) * 100
		fmt.Printf("Current percentage: %.2f%%\n", percentage)

		if writeError != nil {
			return writeError
		}
	}

	return nil
}

package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileNotExist          = errors.New("file doesn't exist")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	err := fileErrors(fromPath, offset)
	if err != nil {
		return err
	}

	// use limited reader instead default reader
	limit = changeZeroLimit(fromPath, limit)
	limitedReader, err := openLimitedReader(fromPath, offset, limit)
	if err != nil {
		return err
	}

	// define bar reader
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(limitedReader)
	defer barReader.Close()

	// define default writer
	writer, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, barReader)
	if err != nil {
		return err
	}

	return nil
}

func fileErrors(fromPath string, offset int64) error {
	fileSt, err := os.Stat(fromPath)
	if os.IsNotExist(err) {
		return ErrFileNotExist
	}
	// check other errors
	if err != nil {
		return err
	}

	fileSize := fileSt.Size()
	if fileSize < 0 {
		return ErrUnsupportedFile
	}
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func changeZeroLimit(fromPath string, limit int64) int64 {
	fileSt, err := os.Stat(fromPath)
	// stay limit unchanged if stat caught error
	if err != nil {
		return limit
	}

	if limit < 1 {
		limit = fileSt.Size()
	}

	return limit
}

func openLimitedReader(fromPath string, offset, limit int64) (io.Reader, error) {
	reader, err := os.Open(fromPath)
	if err != nil {
		return reader, err
	}

	_, err = reader.Seek(offset, io.SeekStart)
	if err != nil {
		return reader, err
	}
	limitedReader := io.LimitReader(reader, limit)

	return limitedReader, nil
}

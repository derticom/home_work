package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3" //nolint: depguard // import is necessary
	"github.com/pkg/errors"     //nolint: depguard // import is necessary
)

const filePermission = 0o644

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrTheSamePath           = errors.New("source and destination are the same")
	ErrInvalidInput          = errors.New("invalid input")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateInput(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, filePermission)
	if err != nil {
		return errors.Wrap(err, "failed to os.OpenFile")
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to sourceFile.Stat")
	}

	// offset больше, чем размер файла - невалидная ситуация
	if sourceFileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	// программа может НЕ обрабатывать файлы, у которых неизвестна длина (например, /dev/urandom)
	if !sourceFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if limit > sourceFileInfo.Size()-offset || limit == 0 {
		limit = sourceFileInfo.Size() - offset
	}

	_, err = sourceFile.Seek(offset, 0)
	if err != nil {
		return errors.Wrap(err, "failed to sourceFile.Seek")
	}

	destination, err := os.Create(toPath)
	if err != nil {
		return errors.Wrap(err, "failed to os.Create")
	}
	defer destination.Close()

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(sourceFile)

	_, err = io.CopyN(destination, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.Wrap(err, "failed to io.CopyN")
	}

	bar.Finish()

	return nil
}

func validateInput(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrInvalidInput
	}

	absFromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return errors.Wrap(err, "failed to filepath.Abs")
	}

	absToPath, err := filepath.Abs(toPath)
	if err != nil {
		return errors.Wrap(err, "failed to filepath.Abs")
	}

	if absFromPath == absToPath {
		return ErrTheSamePath
	}

	if offset < 0 || limit < 0 {
		return ErrInvalidInput
	}
	return nil
}

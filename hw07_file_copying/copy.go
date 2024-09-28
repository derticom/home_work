package main

import (
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint: depguard // import is necessary
	"github.com/pkg/errors"     //nolint: depguard // import is necessary
)

const filePermission = 0o644

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, filePermission)
	if err != nil {
		return ErrUnsupportedFile
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
	if sourceFileInfo.Size() == 0 {
		return errors.New("unknown len source file")
	}

	if limit > sourceFileInfo.Size()-offset {
		limit = sourceFileInfo.Size() - offset
	}

	if limit == 0 {
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
		return errors.Wrap(err, "failed to copy data")
	}

	bar.Finish()

	return nil
}

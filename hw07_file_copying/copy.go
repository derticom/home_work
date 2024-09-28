package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
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

	buf := make([]byte, limit)

	var totalRead int64
	for totalRead != limit {
		n, err := sourceFile.Read(buf[totalRead:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to sourceFile.Read")
		}
		totalRead += int64(n)
	}

	err = os.WriteFile(toPath, buf, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to os.WriteFile")
	}

	return nil
}

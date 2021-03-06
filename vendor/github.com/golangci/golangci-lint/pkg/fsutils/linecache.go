package fsutils

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

type fileLinesCache [][]byte

type LineCache struct {
	files     map[string]fileLinesCache
	fileCache *FileCache
}

func NewLineCache(fc *FileCache) *LineCache {
	return &LineCache{
		files:     map[string]fileLinesCache{},
		fileCache: fc,
	}
}

// GetLine returns a index1-th (1-based index) line from the file on filePath
func (lc *LineCache) GetLine(filePath string, index1 int) (string, error) {
	if index1 == 0 { // some linters, e.g. gosec can do it: it really means first line
		index1 = 1
	}

	rawLine, err := lc.getRawLine(filePath, index1-1)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(rawLine, "\r")), nil
}

func (lc *LineCache) getRawLine(filePath string, index0 int) ([]byte, error) {
	fc, err := lc.getFileCache(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get file %s lines cache", filePath)
	}

	if index0 < 0 {
		return nil, fmt.Errorf("invalid file line index0 < 0: %d", index0)
	}

	if index0 >= len(fc) {
		return nil, fmt.Errorf("invalid file line index0 (%d) >= len(fc) (%d)", index0, len(fc))
	}

	return fc[index0], nil
}

func (lc *LineCache) getFileCache(filePath string) (fileLinesCache, error) {
	fc := lc.files[filePath]
	if fc != nil {
		return fc, nil
	}

	fileBytes, err := lc.fileCache.GetFileBytes(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get file %s bytes from cache", filePath)
	}

	fc = bytes.Split(fileBytes, []byte("\n"))
	lc.files[filePath] = fc
	return fc, nil
}

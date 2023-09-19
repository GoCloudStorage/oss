package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
)

var (
	Client *chunkStorage
)

type chunkStorage struct {
	root string
}

func Init(root string) {
	Client = &chunkStorage{
		root: root,
	}
}

func (c *chunkStorage) SaveChunk(key string, chunkNumber int, data io.Reader, offset int64) error {
	fileDir := path.Join(c.root, key)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir, err: %v", err)
	}
	filepath := path.Join(fileDir, strconv.Itoa(chunkNumber))
	return saveFile(filepath, data, offset)
}

func (c *chunkStorage) MergeChunk(key string, chunkNums int, totalSize int) error {
	fileDir := path.Join(c.root, key)
	dirs, err := os.ReadDir(fileDir)
	if err != nil {
		removeDir(fileDir)
		return fmt.Errorf("failed to read dir [%s], err: %v", fileDir, err)
	}
	if len(dirs) != chunkNums {
		removeDir(fileDir)
		return fmt.Errorf("file chunk not complete, need %d have %d", chunkNums, len(dirs))
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	// create final storage_engine file
	finalFilePath := path.Join(fileDir, "data")
	for i, part := range dirs {
		srcPath := path.Join(fileDir, part.Name())
		size, err := copyFileToFile(srcPath, finalFilePath)
		if err != nil {
			return fmt.Errorf("failed to copy [%d] chunk file to dst file, err: %v", i, err)
		}
		totalSize -= size
	}

	if totalSize != 0 {
		removeDir(fileDir)
		return fmt.Errorf("merge chunk not complete, %d", totalSize)
	}

	// remove chunk file
	for i, part := range dirs {
		if part.Name() != "data" {
			filePath := path.Join(fileDir, part.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove [%d] chunk file, err: %v", i, err)
			}
		}
	}
	return nil
}

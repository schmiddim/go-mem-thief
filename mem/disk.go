package mem

import (
	"fmt"
	"io"
	"os"
)

type IoOps struct {
}

// CreateFile deprecated
func (i *IoOps) CreateFile(filePath string, size int64) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	buffer := make([]byte, 1024*1024)
	for i := int64(0); i < size; i += int64(len(buffer)) {
		_, err := file.Write(buffer)
		if err != nil {
			return err
		}
	}

	return nil

}

func (i *IoOps) WriteFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (i *IoOps) ReadLargeFile(filePath string) ([]byte, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	buf := make([]byte, 4096) // 4 KB Buffer
	var data []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("%w", err)
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}
	return data, nil
}

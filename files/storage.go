package files

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

type Storage interface {
	Upload(filePath string, file io.Reader) (string, error)
	Read(filePath string) (*os.File, error)
}

func NewHostStorage(basePath string, baseUrl string) (*HostStorage, error) {
	storage := &HostStorage{
		basePath: strings.TrimRight(basePath, "/"),
		baseUrl:  strings.TrimRight(baseUrl, "/"),
	}

	if err := storage.initBasePath(); err != nil {
		return nil, err
	}

	return storage, nil
}

type HostStorage struct {
	basePath          string
	baseUrl           string
	ErrorOnFileExists bool
	FileNameEncoder   func(string) string
	FsPerm            os.FileMode
}

func (m HostStorage) initBasePath() error {
	if _, err := os.Stat(m.basePath); !os.IsNotExist(err) {
		return nil
	}
	if err := os.MkdirAll(m.basePath, m.FsPerm); err != nil {
		return err
	}

	return nil
}

func (m HostStorage) prepareFilePath(filePath string) (string, error) {
	dirPath, fileName := path.Split(filePath)
	dirPath = strings.Trim(dirPath, "/")

	fileExt := path.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(fileExt)]

	if m.FileNameEncoder != nil {
		fileName = m.FileNameEncoder(fileName)
	}

	uploadPath := fmt.Sprintf("%s/%s/%s%s",
		m.basePath,
		dirPath,
		fileName,
		fileExt,
	)

	if _, err := os.Stat(uploadPath); !os.IsNotExist(err) {
		if m.ErrorOnFileExists {
			return "", fmt.Errorf("File already exists: %s ", uploadPath)
		}

		uploadPath = fmt.Sprintf(
			"%s/%s/%s_%d%d%s",
			m.basePath,
			dirPath,
			fileName,
			time.Now().Nanosecond(),
			rand.Intn(10000),
			fileExt,
		)
	}

	return uploadPath, nil
}

func (m HostStorage) ForSubFolder(path string) (*HostStorage, error) {
	newStorage := &HostStorage{
		basePath:          fmt.Sprintf("%s/%s", m.basePath, strings.TrimLeft(path, "/")),
		baseUrl:           fmt.Sprintf("%s/%s", m.baseUrl, strings.TrimLeft(path, "/")),
		ErrorOnFileExists: m.ErrorOnFileExists,
		FileNameEncoder:   m.FileNameEncoder,
		FsPerm:            m.FsPerm,
	}

	if err := newStorage.initBasePath(); err != nil {
		return nil, err
	}

	return newStorage, nil
}

func (m HostStorage) BuildUrl(filePath string) string {
	return fmt.Sprintf("%s/%s", m.baseUrl, strings.TrimLeft(filePath, "/"))
}

func (m HostStorage) Upload(filePath string, dataReader io.Reader) (string, error) {
	uploadPath, err := m.prepareFilePath(filePath)
	if err != nil {
		return "", err
	}

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	err = os.MkdirAll(path.Dir(uploadPath), m.FsPerm)
	if err != nil {
		return "", err
	}

	out, cErr := os.OpenFile(uploadPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, m.FsPerm)
	if cErr != nil {
		return "", cErr
	}

	_, err = io.Copy(out, dataReader)
	if err != nil {
		return uploadPath, err
	}

	return uploadPath, out.Close()
}

func (m HostStorage) Read(filePath string) (*os.File, error) {
	realFilePath := fmt.Sprintf("%s/%s", m.basePath, strings.TrimLeft(filePath, "/"))
	if _, err := os.Stat(realFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s not exists: %s ", filePath, err)
	}

	return os.Open(realFilePath)
}

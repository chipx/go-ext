package files

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const basePath = "/tmp/manager_test"

func setUp() {
	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		err := os.RemoveAll(basePath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Remove - ", basePath)
	}
}

func TestHostStorage_PrepareFileName(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath: basePath,
	}
	fileName := fmt.Sprintf("%s/h11/sdasd/1.jpg", basePath)
	newFileName, err := manager.prepareFilePath("h11/sdasd/1.jpg")
	require.Nil(t, err)
	require.NotEmpty(t, newFileName)
	require.Equal(t, fileName, newFileName)
}

func TestHostStorage_PrepareFileName_ErrOnAlreadyExists(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath:          basePath,
		ErrorOnFileExists: true,
	}
	fileName := fmt.Sprintf("%s/h11/sdasd/1.jpg", basePath)
	os.MkdirAll(fmt.Sprintf("%s/h11/sdasd", basePath), os.ModePerm)
	f, createErr := os.Create(fileName)
	if createErr != nil {
		panic(createErr)
	}
	f.WriteString("Test")
	_ = f.Close()

	newFileName, err := manager.prepareFilePath("h11/sdasd/1.jpg")
	require.NotNil(t, err)
	require.Empty(t, newFileName)
}

func TestHostStorage_PrepareFileName_AddSuffixOnAlreadyExists(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath:          basePath,
		ErrorOnFileExists: false,
	}
	fileName := fmt.Sprintf("%s/h11/sdasd/1.jpg", basePath)
	os.MkdirAll(fmt.Sprintf("%s/h11/sdasd", basePath), os.ModePerm)
	f, createErr := os.Create(fileName)
	if createErr != nil {
		panic(createErr)
	}
	f.WriteString("Test")
	_ = f.Close()

	newFileName, err := manager.prepareFilePath("h11/sdasd/1.jpg")
	require.Nil(t, err)
	require.Regexp(t, fmt.Sprintf(`^%s\/h11\/sdasd\/1_\d+\.jpg$`, strings.Replace(basePath, "/", "\\/", 10)), newFileName)
}

func TestHostStorage_PrepareFileName_UseEncoder(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath: basePath,
		FileNameEncoder: func(s string) string {
			require.Equal(t, "1", s)
			return fmt.Sprintf("%s-%s", s, s)
		},
	}

	newFileName, err := manager.prepareFilePath("h11/sdasd/1.jpg")
	require.Nil(t, err)
	require.NotEmpty(t, newFileName)
	require.Equal(t, fmt.Sprintf("%s/h11/sdasd/1-1.jpg", basePath), newFileName)
}

func TestHostStorage_Upload(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath: basePath,
		FsPerm:   os.ModePerm,
	}

	fileData := "Hello world"
	reader := bytes.NewReader([]byte(fileData))
	uploadedPath, err := manager.Upload("h11/sdasd/1.jpg", reader)
	filePath := fmt.Sprintf("%s/h11/sdasd/1.jpg", basePath)
	require.Nil(t, err)
	require.Equal(t, filePath, uploadedPath)

	s, _ := os.Stat(filePath)
	require.Equal(t, s.Mode().Perm(), os.ModePerm)

	uploadedData, rErr := ioutil.ReadFile(uploadedPath)
	if rErr != nil {
		panic(rErr)
	}

	require.Equal(t, fileData, string(uploadedData))
}

func TestManger_Read(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath: basePath,
		FsPerm:   os.ModePerm,
	}

	fileName := fmt.Sprintf("%s/h11/sdasd/1.jpg", basePath)
	mErr := os.MkdirAll(fmt.Sprintf("%s/h11/sdasd", basePath), os.ModePerm)
	if mErr != nil {
		panic(mErr)
	}

	f, createErr := os.Create(fileName)
	if createErr != nil {
		panic(createErr)
	}
	f.WriteString("Test")
	fErr := f.Close()
	if fErr != nil {
		panic(fErr)
	}

	reader, err := manager.Read("h11/sdasd/1.jpg")
	require.Nil(t, err)
	dataFromManager, rErr := ioutil.ReadAll(reader)
	if rErr != nil {
		panic(rErr)
	}
	require.Equal(t, "Test", string(dataFromManager))
}

func TestManger_ForSubFolder(t *testing.T) {
	setUp()
	manager := &hostStorage{
		basePath: basePath,
		baseUrl:  "http://fs.local/",
		FsPerm:   os.ModePerm,
	}

	newManager, err := manager.ForSubFolder("test")
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s/test", manager.basePath), newManager.basePath)
	require.Equal(t, fmt.Sprintf("%s/test", manager.baseUrl), newManager.baseUrl)
}

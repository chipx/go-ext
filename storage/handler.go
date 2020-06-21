package storage

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strings"
)

func NewHttpHandler(storage Storage, baseUrl string, accessCallback func(string, *http.Request) bool) *httpHandler {
	return &httpHandler{
		storage:        storage,
		baseUrl:        baseUrl,
		accessCallback: accessCallback,
	}
}

type httpHandler struct {
	storage        Storage
	accessCallback func(string, *http.Request) bool
	baseUrl        string
	Logger         logrus.FieldLogger
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimLeft(r.URL.Path, h.baseUrl)
	if h.accessCallback != nil && !h.accessCallback(filePath, r) {
		w.WriteHeader(403)
		return
	}

	file, err := h.storage.Read(filePath)
	if err != nil {
		h.logError("File %s for url %s not found: %s", filePath, r.URL.Path, err)
		w.WriteHeader(404)
		return
	}

	fInfo, fInfoErr := file.Stat()
	if fInfoErr != nil {
		h.logError("Get stat for file %s filed: %s", filePath, fInfoErr)
		w.WriteHeader(500)
		return
	}

	http.ServeContent(w, r, path.Base(filePath), fInfo.ModTime(), file)
}

func (h *httpHandler) logError(message string, args ...interface{}) {
	if h.Logger == nil {
		return
	}

	h.Logger.Errorf(message, args...)
}

package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jeremybastin1207/mindia-core/internal/apikey"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	mindialog "github.com/jeremybastin1207/mindia-core/internal/logging"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/settings"
	"github.com/jeremybastin1207/mindia-core/internal/task"
	pathutils "github.com/jeremybastin1207/mindia-core/pkg/path"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Tasks struct {
	ClearCache                  task.ClearCacheTask
	NamedTransformationOperator task.NamedTransformationOperator
	ApiKeyOperator              task.ApiKeyOperator
	AnalyticsOperator           task.AnalyticsOperator
	TaskOperator                task.TaskOperator
	GetMedia                    task.GetMediaTask
	DownloadMedia               task.DownloadMediaTask
	UploadMedia                 task.UploadMediaTask
	DeleteMedia                 task.DeleteMediaTask
	MoveMedia                   task.MoveMediaTask
	CopyMedia                   task.CopyMediaTask
	TagMedia                    task.TagMediaTask
	ColorizeMedia               task.ColorizeMediaTask
}

type ApiServer struct {
	host          string
	port          int
	masterKey     string
	apikeyStorage apikey.Storer
	logger        mindialog.Logger
	tasks         Tasks
}

func NewApiServer(
	masterKey string,
	host string,
	port int,
	apikeyStorage apikey.Storer,
	logger mindialog.Logger,
	tasks Tasks,
) ApiServer {
	return ApiServer{
		host,
		port,
		masterKey,
		apikeyStorage,
		logger,
		tasks,
	}
}

func (s *ApiServer) Serve() {
	r := mux.NewRouter().PathPrefix("").Subrouter()
	r.Methods("GET", "OPTIONS").PathPrefix("/metrics").Handler(promhttp.Handler())

	apir := r.PathPrefix("/" + settings.ApiVersion).Subrouter()

	sr := apir.PathPrefix("/cache").Subrouter()
	sr.Use(apiKeyMiddleware(s.apikeyStorage, s.masterKey))
	sr.Methods("DELETE", "OPTIONS").Path("/clear").HandlerFunc(apiHandler(s.handleClearCache))

	sr = apir.PathPrefix("/named_transformation").Subrouter()
	sr.Use(apiKeyMiddleware(s.apikeyStorage, s.masterKey))
	sr.Methods("GET", "OPTIONS").HandlerFunc(apiHandler(s.handleReadNamedTransformations))
	sr.Methods("POST", "OPTIONS").HandlerFunc(apiHandler(s.handleCreateNamedTransformation))
	sr.Methods("PATCH", "OPTIONS").Path("/{name}").HandlerFunc(apiHandler(s.handleUpdateNamedTransformation))
	sr.Methods("DELETE", "OPTIONS").Path("/{name}").HandlerFunc(apiHandler(s.handleDeleteNamedTransformation))
	sr.Methods("DELETE", "OPTIONS").HandlerFunc(apiHandler(s.handleDeleteAllNamedTransformations))

	sr = apir.PathPrefix("/api_key").Subrouter()
	sr.Use(masterKeyMiddleware(s.masterKey))
	sr.Methods("GET", "OPTIONS").HandlerFunc(apiHandler(s.handleReadKeys))
	sr.Methods("POST", "OPTIONS").HandlerFunc(apiHandler(s.handleCreateKey))
	sr.Methods("DELETE", "OPTIONS").Path("/{api_key}").HandlerFunc(apiHandler(s.handleDeleteKey))

	sr = apir.PathPrefix("/task").Subrouter()
	sr.Use(apiKeyMiddleware(s.apikeyStorage, s.masterKey))

	sr = apir.PathPrefix("/analytics").Subrouter()
	sr.Use(apiKeyMiddleware(s.apikeyStorage, s.masterKey))
	sr.Methods("GET", "OPTIONS").Path("/space").HandlerFunc(apiHandler(s.handleReadSpaceUsage))

	sr = apir.PathPrefix("/download").Subrouter()
	sr.Methods("GET", "OPTIONS").Path("/{path:.*}").HandlerFunc(apiHandler(s.handleDownloadMedia))

	sr = apir.PathPrefix("").Subrouter()
	sr.Use(apiKeyMiddleware(s.apikeyStorage, s.masterKey))
	sr.Methods("POST", "OPTIONS").Path("/download/archive").HandlerFunc(apiHandler(s.handleDownloadMultipleMedias))
	sr.Methods("POST", "OPTIONS").Path("/upload").HandlerFunc(apiHandler(s.handleUploadMedia))
	sr.Methods("POST", "OPTIONS").Path("/upload/{path:.*}").HandlerFunc(apiHandler(s.handleUploadMedia))
	sr.Methods("PUT", "OPTIONS").Path("/move").HandlerFunc(apiHandler(s.handleMoveMedia))
	sr.Methods("PUT", "OPTIONS").Path("/copy").HandlerFunc(apiHandler(s.handleCopyMedia))
	sr.Methods("POST", "OPTIONS").Path("/tag/{path:.*}").HandlerFunc(apiHandler(s.handleTagMedia))
	sr.Methods("POST", "OPTIONS").Path("/colorize/{path:.*}").HandlerFunc(apiHandler(s.handleColorizeMedia))
	sr.Methods("GET", "OPTIONS").Path("/file/{path:.*}").HandlerFunc(apiHandler(s.handleGetMedia))
	sr.Methods("GET", "OPTIONS").Path("/files/{path:.*}").HandlerFunc(apiHandler(s.handleGetMultipleMedias))
	sr.Methods("DELETE", "OPTIONS").Path("/delete_bulk").HandlerFunc(apiHandler(s.handleDeleteMultipleMedias))
	sr.Methods("DELETE", "OPTIONS").Path("/{path:.*}").HandlerFunc(apiHandler(s.handleDeleteMedia))

	s.logger.Info("server is starting...")

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	listenAddr := fmt.Sprintf("%s:%d", s.host, s.port)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(s.logger)(r)),
		ErrorLog:     s.logger.GetErrorLogger(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		s.logger.Info("server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Critical(fmt.Sprintf("could not gracefully shutdown the server: %v", err))
		}
		close(done)
	}()

	s.logger.Info(fmt.Sprintf("server is ready to handle requests at %s", listenAddr))
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Critical(fmt.Sprintf("could not listen on %s: %v", listenAddr, err))
	}

	<-done
	s.logger.Info("server stopped")
}

func (s *ApiServer) handleClearCache(w http.ResponseWriter, r *http.Request) error {
	err := s.tasks.ClearCache.ClearAll()
	if err != nil {
		return err
	}
	return writeMessage(w, "Successfully cleared cache")
}

func (s *ApiServer) handleReadNamedTransformations(w http.ResponseWriter, r *http.Request) error {
	namedTransformations, err := s.tasks.NamedTransformationOperator.GetAll()
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, namedTransformations))
}

func (s *ApiServer) handleCreateNamedTransformation(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Name               string `json:"name"`
		TransformationsStr string `json:"transformations"`
	}
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	t, err := s.tasks.NamedTransformationOperator.Create(b.Name, b.TransformationsStr)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, *t))
}

func (s *ApiServer) handleUpdateNamedTransformation(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Transformations string `json:"transformations"`
	}
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	t, err := s.tasks.NamedTransformationOperator.Update(mux.Vars(r)["name"], b.Transformations)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, *t))
}

func (s *ApiServer) handleDeleteNamedTransformation(w http.ResponseWriter, r *http.Request) error {
	var (
		vars = mux.Vars(r)
		name = vars["name"]
	)
	err := s.tasks.NamedTransformationOperator.Delete(name)
	if err != nil {
		return err
	}
	return writeMessage(w, "successfully deleted named transformation")
}

func (s *ApiServer) handleDeleteAllNamedTransformations(w http.ResponseWriter, r *http.Request) error {
	err := s.tasks.NamedTransformationOperator.DeleteAll()
	if err != nil {
		return err
	}
	return writeMessage(w, "successfully deleted all named transformations")
}

func (s *ApiServer) handleReadKeys(w http.ResponseWriter, r *http.Request) error {
	apiKeys, err := s.tasks.ApiKeyOperator.GetAll()
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, apiKeys))
}

func (s *ApiServer) handleCreateKey(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Name string
	}
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	apiKey, err := s.tasks.ApiKeyOperator.Create(b.Name)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, *apiKey))
}

func (s *ApiServer) handleDeleteKey(w http.ResponseWriter, r *http.Request) error {
	err := s.tasks.ApiKeyOperator.Delete(mux.Vars(r)["api_key"])
	if err != nil {
		return err
	}
	return writeMessage(w, "successfully deleted apikey")
}

func (s *ApiServer) handleReadSpaceUsage(w http.ResponseWriter, r *http.Request) error {
	spaceUsage, err := s.tasks.AnalyticsOperator.SpaceUsage()
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, *spaceUsage))
}

func (s *ApiServer) handleDownloadMultipleMedias(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Paths []string `json:"paths"`
	}
	paths := []media.Path{}
	b, err := parseBody[Body](w, r)
	if err != nil {
		return err
	}
	for _, p := range b.Paths {
		paths = append(paths, media.NewPath(p))
	}
	bytes, err := s.tasks.DownloadMedia.DownloadMultiple(paths)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(*bytes)
	return err
}

func (s *ApiServer) handleDownloadMedia(w http.ResponseWriter, r *http.Request) error {
	var (
		vars = mux.Vars(r)
		path = vars["path"]
	)
	transformations, imagePath := parsePath(path)
	body, contentType, err := s.tasks.DownloadMedia.Download(
		transformations,
		media.NewPath(imagePath),
	)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", *contentType)
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, body)
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

func (s *ApiServer) handleUploadMedia(w http.ResponseWriter, r *http.Request) error {
	var (
		vars                  = mux.Vars(r)
		path                  = vars["path"]
		transformations       string
		parsedTransformations []string
		body                  io.Reader
		filename              string
		contentType           string
	)

	transformations, imagePath := parsePath(path)
	if transformations != "" {
		parsedTransformations = append(parsedTransformations, transformations)
	}
	reader, err := r.MultipartReader()
	if err != nil {
		return &mindiaerr.Error{ErrCode: mindiaerr.ErrBadRequest, Msg: err}
	}

	for {
		p, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeInternal, Msg: err}
		}

		if p.FormName() == "transformations" {
			b := bytes.Buffer{}
			io.Copy(bufio.NewWriter(&b), bufio.NewReader(p))
			json.Unmarshal(b.Bytes(), &parsedTransformations)
		}
		if p.FormName() == "file" {
			body = bufio.NewReader(p)
			filename = p.FileName()
			contentType = p.Header.Get("Content-Type")
			break // Don't read next part otherwise body will be empty.
		}
	}

	if body == nil {
		return &mindiaerr.Error{
			ErrCode: mindiaerr.ErrBadRequest,
			Msg:     fmt.Errorf("form-data 'file' must be provided"),
		}
	}

	uploadedMedia, err := s.tasks.UploadMedia.Upload(
		pathutils.JoinPath(imagePath, uuid.New().String()+filepath.Ext(filename)),
		body,
		media.ContentType(contentType),
		0,
		parsedTransformations,
	)

	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, *uploadedMedia))
}

func (s *ApiServer) handleMoveMedia(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Src string
		Dst string
	}
	b, err := parseBody[Body](w, r)
	if err != nil {
		return err
	}
	m, err := s.tasks.MoveMedia.Move(media.NewPath(b.Src), media.NewPath(b.Dst))
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, m))
}

func (s *ApiServer) handleCopyMedia(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Src string
		Dst string
	}
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	m, err := s.tasks.CopyMedia.Copy(media.NewPath(b.Src), media.NewPath(b.Dst))
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, m))
}

func (s *ApiServer) handleTagMedia(w http.ResponseWriter, r *http.Request) error {
	path := media.NewPath(toAbsolutePath(mux.Vars(r)["path"]))

	result, err := s.tasks.TagMedia.Tag(path)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, result))
}

func (s *ApiServer) handleColorizeMedia(w http.ResponseWriter, r *http.Request) error {
	path := media.NewPath(toAbsolutePath(mux.Vars(r)["path"]))

	result, err := s.tasks.ColorizeMedia.Colorize(path)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, result))
}

func (s *ApiServer) handleGetMedia(w http.ResponseWriter, r *http.Request) error {
	path := media.NewPath(toAbsolutePath(mux.Vars(r)["path"]))

	file, err := s.tasks.GetMedia.Get(path)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, file))
}

func (s *ApiServer) handleGetMultipleMedias(w http.ResponseWriter, r *http.Request) error {
	type Body struct {
		Offset int    `schema:"offset"`
		Limit  int    `schema:"limit"`
		SortBy string `schema:"sort_by"`
	}
	query, err := parseQuery[Body](r.URL)
	if err != nil {
		return err
	}

	sortBy := strings.Split(query.SortBy, ":")
	var ascBool bool
	if sortBy[1] == "asc" {
		ascBool = true
	} else {
		ascBool = false
	}

	medias, err := s.tasks.GetMedia.GetMultiple(
		media.NewPath(toAbsolutePath(mux.Vars(r)["path"])),
		query.Offset,
		query.Limit,
		sortBy[0],
		ascBool,
	)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, medias))
}

func (s *ApiServer) handleDeleteMedia(w http.ResponseWriter, r *http.Request) error {
	path := media.NewPath(toAbsolutePath(mux.Vars(r)["path"]))

	err := s.tasks.DeleteMedia.Delete(path)
	if err != nil {
		return err
	}
	return writeMessage(w, "Successfully deleted file")
}

func (s *ApiServer) handleDeleteMultipleMedias(w http.ResponseWriter, r *http.Request) error {
	var paths []media.Path
	body, err := parseBody[[]*string](w, r)
	if err != nil {
		return err
	}
	for _, p := range *body {
		paths = append(paths, media.NewPath(*p))
	}
	err = s.tasks.DeleteMedia.DeleteMultiple(paths)
	if err != nil {
		return err
	}
	return writeJSON(w, encodeJSON(w, "Medias deleted successfully!"))
}

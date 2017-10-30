package backend

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/fxkr/openview/backend/cache"
	"github.com/fxkr/openview/backend/model"
	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
)

type Application struct {
	config  *Config
	router  chi.Router
	service Service
}

func NewApplication(config *Config) (*Application, error) {
	Initialize()
	defer Terminate()

	c, err := cache.NewFileCache(config.CacheDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mc, err := cache.NewMiniRedisCache(cache.MiniRedisCacheConfig{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	app := &Application{
		config:  config,
		router:  chi.NewRouter(),
		service: NewService(config.ImageDir, config.ResourceDir, c, mc),
	}

	r := app.router

	r.Get("/static/*", app.handleResource)
	r.Get("/*", app.handlePath)
	r.NotFound(handler.Status(http.StatusNotFound).ServeHTTP)

	return app, nil
}

func (app *Application) Run() error {
	err := http.ListenAndServe(":3000", app.router)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (app *Application) handleResource(w http.ResponseWriter, r *http.Request) {
	relativePath, err := safe.SafeNewRelativePath(chi.URLParam(r, "*"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}
	fullPath := app.config.ResourceDir.Join(relativePath)
	http.ServeFile(w, r, fullPath.String())
}

func (app *Application) handlePath(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	_, haveSize := r.URL.Query()["size"]
	if action == "" && haveSize {
		action = "thumb"
	}

	switch action {
	case "":
		app.handleFile(w, r)
	case "thumb":
		app.handleThumbnail(w, r)
	case "info":
		app.handleDirectoryInfo(w, r)
	case "image-info":
		app.handleImageInfo(w, r)
	default:
		handler.Status(http.StatusBadRequest).ServeHTTP(w, r)
	}
}

func (app *Application) handleFile(w http.ResponseWriter, r *http.Request) {
	path, err := safe.SafeNewRelativePath(strings.TrimSuffix(chi.URLParam(r, "*"), "/"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	app.service.Get(path).ServeHTTP(w, r)
}

func (app *Application) handleDirectoryInfo(w http.ResponseWriter, r *http.Request) {
	page, err := model.NewPage(r.URL.Query().Get("page_token"), r.URL.Query().Get("page_size"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	path, err := safe.SafeNewRelativePath(strings.TrimSuffix(chi.URLParam(r, "*"), "/"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	app.service.GetDirectory(path, page).ServeHTTP(w, r)
}

func (app *Application) handleImageInfo(w http.ResponseWriter, r *http.Request) {
	path, err := safe.SafeNewRelativePath(strings.TrimSuffix(chi.URLParam(r, "*"), "/"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	app.service.GetImage(path).ServeHTTP(w, r)
}

func (app *Application) handleThumbnail(w http.ResponseWriter, r *http.Request) {
	path, err := safe.SafeNewRelativePath(strings.TrimSuffix(chi.URLParam(r, "*"), "/"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	size, err := model.NewThumbSize(r.URL.Query().Get("size"))
	if err != nil {
		handler.StatusError(http.StatusBadRequest, errors.WithStack(err)).ServeHTTP(w, r)
		return
	}

	app.service.GetImageThumbnail(path, size).ServeHTTP(w, r)
}

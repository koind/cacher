package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/koind/cacher/internal/domain/repository"
	"github.com/koind/cacher/internal/domain/service"
	"io/ioutil"
	"net/http"
)

// HTTP сервер
type HTTPServer struct {
	http.Server
	router       *gin.Engine
	domain       string
	cacheService *service.CacheService
}

// Возвращает новый HTTP сервер
func NewHTTPServer(cacheService *service.CacheService, domain string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	hs := HTTPServer{router: r, domain: domain, cacheService: cacheService}

	hs.router.GET("/get/:key", hs.GetOneHandle)
	hs.router.GET("/list", hs.GetAllHandle)
	hs.router.POST("/upsert", hs.UpsertHandle)
	hs.router.DELETE("/delete/:key", hs.DeleteHandle)

	http.Handle("/", r)

	return &hs
}

// Запускает HTTP сервер
func (s *HTTPServer) Start() error {
	return http.ListenAndServe(s.domain, s.router)
}

// Возвращет одну запись по ключу
func (s *HTTPServer) GetOneHandle(c *gin.Context) {
	key := c.Param("key")

	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, responseError("ключ кэша не может быть пустым"))

		return
	}

	cache, err := s.cacheService.GetOneByKey(c, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(cache))
}

// Возвращет все записи
func (s *HTTPServer) GetAllHandle(c *gin.Context) {
	list, err := s.cacheService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err.Error()))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(list))
}

// Обновить запись, если существует, и создает, если нет
func (s *HTTPServer) UpsertHandle(c *gin.Context) {
	cache := repository.Cache{}

	if err := c.ShouldBindJSON(&cache); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))

		return
	}

	if err := validator.New().StructCtx(c, cache); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))

		return
	}

	newCache, err := s.cacheService.Upsert(c, cache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err.Error()))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(newCache))
}

// Удаляет одну запись по ключу
func (s *HTTPServer) DeleteHandle(c *gin.Context) {
	key := c.Param("key")

	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, responseError("ключ кэша не может быть пустым"))

		return
	}

	err := s.cacheService.Delete(c, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(nil))
}

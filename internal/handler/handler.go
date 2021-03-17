package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/koind/cacher/docs"
	"github.com/koind/cacher/internal/domain/repository"
	"github.com/koind/cacher/internal/domain/service"
	"github.com/pkg/errors"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io/ioutil"
	"net/http"
)

var (
	CacheKeyCannotBeEmptyErr = errors.New("ключ кэша не может быть пустым")
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
	hs.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
// @Summary Возвращет одну запись по ключу
// @Description Возвращет одну запись по ключу
// @Tags api
// @ID get-one-handle
// @Produce  json
// @Param key path string true "ключ кэша"
// @Success 200 {object} repository.Cache
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router /get/{key} [get]
func (s *HTTPServer) GetOneHandle(c *gin.Context) {
	key := c.Param("key")

	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, responseError(CacheKeyCannotBeEmptyErr))

		return
	}

	cache, err := s.cacheService.GetOneByKey(c, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError(err))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(cache))
}

// Возвращет все записи
// @Summary Возвращет все записи
// @Description Возвращет все записи
// @Tags api
// @ID get-all-handle
// @Produce  json
// @Success 200 {array} repository.Cache
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router /list [get]
func (s *HTTPServer) GetAllHandle(c *gin.Context) {
	list, err := s.cacheService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(list))
}

// Обновить запись, если существует, и создает, если нет
// @Summary Обновить запись, если существует, и создает, если нет
// @Description Обновить запись, если существует, и создает, если нет
// @Tags api
// @ID upsert-Handle
// @Accept  json
// @Produce  json
// @Param data body repository.Cache true "данные кэша"
// @Success 200 {object} repository.Cache
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router /upsert [post]
func (s *HTTPServer) UpsertHandle(c *gin.Context) {
	cache := repository.Cache{}

	if err := c.ShouldBindJSON(&cache); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err))

		return
	}

	if err := validator.New().StructCtx(c, cache); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err))

		return
	}

	newCache, err := s.cacheService.Upsert(c, cache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(newCache))
}

// Удаляет одну запись по ключу
// @Summary Удаляет одну запись по ключу
// @Description Удаляет одну запись по ключу
// @Tags api
// @ID delete-handle
// @Produce  json
// @Param key path string true "ключ кэша"
// @Success 200 {object} handler.response
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router /delete/{key} [delete]
func (s *HTTPServer) DeleteHandle(c *gin.Context) {
	key := c.Param("key")

	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, responseError(CacheKeyCannotBeEmptyErr))

		return
	}

	err := s.cacheService.Delete(c, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError(err))

		return
	}

	c.JSON(http.StatusOK, responseSuccess(nil))
}

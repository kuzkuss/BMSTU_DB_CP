package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	postUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/post/usecase"
)

type Delivery struct {
	PostUC postUsecase.UseCaseI
}

func (delivery *Delivery) CreatePosts(c echo.Context) error {
	posts := make([]*models.Post, 0, 10)
	err := c.Bind(&posts)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err = delivery.PostUC.CreatePosts(posts, c.Param("slug_or_id"))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrConflict):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusConflict, models.ErrConflict.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, posts)
}

func (delivery *Delivery) UpdatePost(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	var post models.Post
	err = c.Bind(&post)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	post.Id = id

	err = delivery.PostUC.UpdatePost(&post)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, post)
}

func (delivery *Delivery) SelectPost(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	queryRelated := c.QueryParam("related")
	var related []string

	if queryRelated != "" {
		related = strings.Split(queryRelated, ",")
		for _, elem := range related {
			if elem != "user" && elem != "forum" && elem != "thread" {
				c.Logger().Error(models.ErrBadRequest)
				return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
			}
		}
	}	

	post, err := delivery.PostUC.SelectPost(id, related)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, post)
}

func (delivery *Delivery) SelectThreadPosts(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	sinceStr := c.QueryParam("since")
	since, err := strconv.Atoi(sinceStr)
	if sinceStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	sort := c.QueryParam("sort")
	if sort != "flat" && sort != "tree" && sort != "parent_tree" {
		sort = "flat"
	}

	posts, err := delivery.PostUC.SelectThreadPosts(c.Param("slug_or_id"), limit, since, desc, sort)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, posts)
}

func NewDelivery(e *echo.Echo, postUC postUsecase.UseCaseI) {
	handler := &Delivery{
		PostUC: postUC,
	}

	e.GET("/api/post/:id/details", handler.SelectPost)
	e.POST("/api/post/:id/details", handler.UpdatePost)
	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	e.GET("/api/thread/:slug_or_id/posts", handler.SelectThreadPosts)
}

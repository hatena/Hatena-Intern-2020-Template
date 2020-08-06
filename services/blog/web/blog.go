package web

import (
	"net/http"
	"strconv"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	echo "github.com/labstack/echo/v4"
)

func (s *Server) BlogHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		path := c.Param("path")
		blog, err := s.app.FindBlogByPath(c.Request().Context(), path)
		if err != nil {
			if err == app.ErrNotFound {
				return c.String(http.StatusNotFound, "not found")
			}
			return err
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		entries, hasNextPage, err := s.app.ListEntriesByBlog(c.Request().Context(), blog, page, 5)
		return c.Render(http.StatusOK, "blog.html", map[string]interface{}{
			"IsAuthor":    user != nil && blog.UserID == user.ID,
			"Blog":        blog,
			"Entries":     entries,
			"Page":        page,
			"PrevPage":    page - 1,
			"NextPage":    page + 1,
			"HasPrevPage": page > 1,
			"HasNextPage": hasNextPage,
		})
	}
}

func (s *Server) EntryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		path := c.Param("path")
		blog, err := s.app.FindBlogByPath(c.Request().Context(), path)
		if err != nil {
			if err == app.ErrNotFound {
				return c.String(http.StatusNotFound, "not found")
			}
			return err
		}
		entryID, err := domain.ParseEntryID(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid id")
		}
		entry, err := s.app.FindEntryByID(c.Request().Context(), blog, entryID)
		if err != nil {
			if err == app.ErrNotFound {
				return c.String(http.StatusNotFound, "not found")
			}
			return err
		}
		return c.Render(http.StatusOK, "entry.html", map[string]interface{}{
			"IsAuthor": user != nil && blog.UserID == user.ID,
			"Blog":     blog,
			"Entry":    entry,
		})
	}
}

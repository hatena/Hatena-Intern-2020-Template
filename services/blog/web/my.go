package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	echo "github.com/labstack/echo/v4"
)

func (s *Server) MyBlogsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		blogs, hasNextPage, err := s.app.ListBlogsByUser(c.Request().Context(), user, page, 10)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "my-blogs.html", map[string]interface{}{
			"User":        user,
			"Blogs":       blogs,
			"Page":        page,
			"PrevPage":    page - 1,
			"NextPage":    page + 1,
			"HasPrevPage": page > 1,
			"HasNextPage": hasNextPage,
		})
	}
}

func (s *Server) WillCreateBlogHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "my-blogs-create.html", nil)
	}
}

func (s *Server) CreateBlogHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		params := new(struct {
			Path        string `form:"path"`
			Title       string `form:"title"`
			Description string `form:"description"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		blog, err := s.app.CreateBlog(c.Request().Context(), user, params.Path, params.Title, params.Description)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/my/blogs/%s", blog.Path))
	}
}

func (s *Server) MyBlogHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		entries, hasNextPage, err := s.app.ListEntriesByBlog(c.Request().Context(), blog, page, 10)
		return c.Render(http.StatusOK, "my-blog.html", map[string]interface{}{
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

func (s *Server) WillEditBlogHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		return c.Render(http.StatusOK, "my-blog-edit.html", map[string]interface{}{
			"Blog": blog,
		})
	}
}

func (s *Server) EditBlogHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		params := new(struct {
			Title       string `form:"title"`
			Description string `form:"description"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		blog, err = s.app.EditBlog(c.Request().Context(), user, blog, params.Title, params.Description)
		if err != nil {
			if err == app.ErrPermissionDenied {
				return c.String(http.StatusForbidden, "permission denied")
			}
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/my/blogs/%s", blog.Path))
	}
}

func (s *Server) DeleteBlogHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		err = s.app.DeleteBlog(c.Request().Context(), user, blog)
		if err != nil {
			if err == app.ErrPermissionDenied {
				return c.String(http.StatusForbidden, "permission denied")
			}
			return err
		}
		return c.Redirect(http.StatusSeeOther, "/my/blogs")
	}
}

func (s *Server) WillPublishEntryHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		return c.Render(http.StatusOK, "my-entries-publish.html", map[string]interface{}{
			"Blog": blog,
		})
	}
}

func (s *Server) PublishEntryHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
		}
		params := new(struct {
			Title string `form:"title"`
			Body  string `form:"body"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = s.app.PublishEntry(c.Request().Context(), user, blog, params.Title, params.Body)
		if err != nil {
			if err == app.ErrPermissionDenied {
				return c.String(http.StatusForbidden, "permission denied")
			}
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/my/blogs/%s", blog.Path))
	}
}

func (s *Server) MyEntryHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
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
		return c.Render(http.StatusOK, "my-entry.html", map[string]interface{}{
			"Blog":  blog,
			"Entry": entry,
		})
	}
}

func (s *Server) EditEntryHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
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
		params := new(struct {
			Title string `form:"title"`
			Body  string `form:"body"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = s.app.EditEntry(c.Request().Context(), user, blog, entry, params.Title, params.Body)
		if err != nil {
			if err == app.ErrPermissionDenied {
				return c.String(http.StatusForbidden, "permission denied")
			}
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/my/blogs/%s", blog.Path))
	}
}

func (s *Server) UnpublishEntryHandler() echo.HandlerFunc {
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
		if blog.UserID != user.ID {
			return c.String(http.StatusForbidden, "permission denied")
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
		err = s.app.UnpublishEntry(c.Request().Context(), user, blog, entry)
		if err != nil {
			if err == app.ErrPermissionDenied {
				return c.String(http.StatusForbidden, "permission denied")
			}
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/my/blogs/%s", blog.Path))
	}
}

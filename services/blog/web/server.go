package web

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server は Web サーバーを表す構造体
type Server struct {
	e   *echo.Echo
	app *app.App
}

// NewServer は Web サーバーを作成する
func NewServer(app *app.App) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	renderer, err := newRenderer()
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	return &Server{e, app}, nil
}

// Echo はサーバーが使用する Echo のインスタンスを取得する
func (s *Server) Echo() *echo.Echo {
	return s.e
}

// Start は Web サーバーを開始する
func (s *Server) Start(addr string) error {
	s.attachMiddlewares()
	s.attachHandlers()
	return s.e.Start(addr)
}

// Shutdown は Web サーバーを停止する
func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) attachMiddlewares() {
	s.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(ctx echo.Context) bool {
			// ヘルスチェックのログを無効化
			return strings.HasPrefix(ctx.Path(), "/server/health")
		},
	}))
	s.e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf_token",
		Skipper: func(c echo.Context) bool {
			r := c.Request()
			if h := r.Header.Get("X-Requested-With"); h != "" {
				return true
			}
			return false
		},
	}))
	s.e.Use(middleware.Secure())
	s.e.Use(s.CustomContextMiddleware())
}

func (s *Server) attachHandlers() {
	s.e.GET("/server/health", s.healthcheckHandler())

	s.e.GET("/", s.IndexHandler())

	s.e.GET("/signup", s.WillSignupHandler())
	s.e.POST("/signup", s.SignupHandler())
	s.e.GET("/signin", s.WillSigninHandler())
	s.e.POST("/signin", s.SigninHandler())
	s.e.GET("/signout", s.WillSignoutHandler(), requireSessionMiddleware())
	s.e.POST("/signout", s.SignoutHandler(), requireSessionMiddleware())

	s.e.GET("/my/blogs", s.MyBlogsHandler(), requireSessionMiddleware())
	s.e.GET("/my/blogs/-/create", s.WillCreateBlogHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs", s.CreateBlogHandler(), requireSessionMiddleware())

	s.e.GET("/my/blogs/:path", s.MyBlogHandler(), requireSessionMiddleware())
	s.e.GET("/my/blogs/:path/edit", s.WillEditBlogHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs/:path/edit", s.EditBlogHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs/:path/delete", s.DeleteBlogHandler(), requireSessionMiddleware())

	s.e.GET("/my/blogs/:path/entries/-/publish", s.WillPublishEntryHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs/:path/entries/-/publish", s.PublishEntryHandler(), requireSessionMiddleware())

	s.e.GET("/my/blogs/:path/entries/:id", s.MyEntryHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs/:path/entries/:id/edit", s.EditEntryHandler(), requireSessionMiddleware())
	s.e.POST("/my/blogs/:path/entries/:id/unpublish", s.UnpublishEntryHandler(), requireSessionMiddleware())

	s.e.GET("/blogs/:path", s.BlogHandler())
	s.e.GET("/blogs/:path/entries/:id", s.EntryHandler())
}

// CustomContext はカスタマイズされたコンテキスト
type CustomContext struct {
	echo.Context
	User *domain.User
}

func (s *Server) CustomContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := s.userFromSessionCookie(c)
			cc := &CustomContext{Context: c, User: user}
			return next(cc)
		}
	}
}

func (s *Server) healthcheckHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm alive.")
	}
}

func (s *Server) IndexHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		blogs, hasNextPage, err := s.app.ListBlogs(c.Request().Context(), page, 10)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
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

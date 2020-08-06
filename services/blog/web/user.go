package web

import (
	"net/http"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	echo "github.com/labstack/echo/v4"
)

const sessionCookieKey = "INTERN_2020_BLOG_SESSION"
const sessionLifespan = 24 * time.Hour

func (s *Server) userFromSessionCookie(c echo.Context) *domain.User {
	key, err := c.Cookie(sessionCookieKey)
	if err != nil || key.Value == "" {
		return nil
	}
	user, _, err := s.app.FindUserBySessionKey(c.Request().Context(), key.Value)
	if err != nil {
		return nil
	}
	return user
}

func requireSessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*CustomContext)
			if cc.User == nil {
				return c.Redirect(http.StatusSeeOther, "/signin")
			}
			return next(c)
		}
	}
}

func getSessionUser(c echo.Context) *domain.User {
	cc := c.(*CustomContext)
	return cc.User
}

func (s *Server) WillSignupHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		// サインイン済みの場合は index へ
		if user != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		return c.Render(http.StatusOK, "signup.html", nil)
	}
}

func (s *Server) SignupHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(struct {
			Name     string `form:"name"`
			Password string `form:"password"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		sessionExpiresAt := time.Now().Add(sessionLifespan)
		_, sess, err := s.app.Signup(c.Request().Context(), params.Name, params.Password, sessionExpiresAt)
		if err != nil {
			if err == app.ErrInvalidArgument {
				return c.String(http.StatusBadRequest, "invalid argument")
			}
			if err == app.ErrAlreadyRegistered {
				return c.String(http.StatusBadRequest, "already registered")
			}
			return err
		}
		c.SetCookie(&http.Cookie{
			Name:     sessionCookieKey,
			Value:    sess.Key,
			Expires:  sess.ExpiresAt,
			HttpOnly: true,
		})
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) WillSigninHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := getSessionUser(c)
		// サインイン済みの場合は index へ
		if user != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		return c.Render(http.StatusOK, "signin.html", nil)
	}
}

func (s *Server) SigninHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(struct {
			Name     string `form:"name"`
			Password string `form:"password"`
		})
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		sessionExpiresAt := time.Now().Add(sessionLifespan)
		_, sess, err := s.app.Signin(c.Request().Context(), params.Name, params.Password, sessionExpiresAt)
		if err != nil {
			if err == app.ErrInvalidArgument {
				return c.String(http.StatusBadRequest, "invalid argument")
			}
			if err == app.ErrAuthenticationFailed {
				return c.String(http.StatusUnauthorized, "authentication failed")
			}
			return err
		}
		c.SetCookie(&http.Cookie{
			Name:     sessionCookieKey,
			Value:    sess.Key,
			Expires:  sess.ExpiresAt,
			HttpOnly: true,
		})
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) WillSignoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "signout.html", nil)
	}
}

func (s *Server) SignoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     sessionCookieKey,
			Value:    "",
			HttpOnly: true,
		})
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

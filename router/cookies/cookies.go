package cookies

import (
	"andurel-site/models"
	"github.com/google/uuid"

	"github.com/labstack/echo/v5"
)

const (
	isAuthenticated = "is_authenticated"
	isAdmin         = "is_admin"
	userID          = "user_id"
)

// Session carries the runtime session cookie names derived from app identity
// so the cookies package holds no global state. It is constructed once and
// injected where session/flash cookies are needed.
type Session struct {
	appSessionName   string
	flashSessionName string
}

func NewSession(appSessionName, projectName, environment string) *Session {
	return &Session{
		appSessionName:   appSessionName,
		flashSessionName: buildFlashSessionName(projectName, environment),
	}
}

type App struct {
	CurrentPath     string
	UserID          uuid.UUID
	IsAdmin         bool
	IsAuthenticated bool
}

func (s *Session) CreateAppSession(c *echo.Context, user models.User) error {
	sess, err := getSession(s.appSessionName, c)
	if err != nil {
		return err
	}

	sess.Values[isAuthenticated] = true
	sess.Values[isAdmin] = user.IsAdmin
	sess.Values[userID] = user.ID.String()

	return sess.Save(c.Request(), c.Response())
}

func (s *Session) DestroyAppSession(c *echo.Context) error {
	sess, err := getSession(s.appSessionName, c)
	if err != nil {
		return err
	}

	sess.Options.MaxAge = -1
	return sess.Save(c.Request(), c.Response())
}

func (s *Session) ExtractFromCookieApp(c *echo.Context) App {
	sess, err := getSession(s.appSessionName, c)
	if err != nil {
		return App{}
	}

	app := App{}

	if v, ok := sess.Values[isAuthenticated].(bool); ok {
		app.IsAuthenticated = v
	}
	if v, ok := sess.Values[isAdmin].(bool); ok {
		app.IsAdmin = v
	}
	if v, ok := sess.Values[userID].(string); ok {
		app.UserID, _ = uuid.Parse(v)
	}

	return app
}

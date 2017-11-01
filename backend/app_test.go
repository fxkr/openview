package backend

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pkg/errors"
	. "gopkg.in/check.v1"

	"github.com/fxkr/openview/backend/util/safe"
)

func TestApp(t *testing.T) {
	_ = Suite(&AppSuite{})
	TestingT(t)
}

type AppSuite struct {
	app *Application

	tempDir  safe.Path
	cacheDir safe.Path
	imageDir safe.Path
}

func (s *AppSuite) SetUpTest(c *C) {
	tempDir, err := ioutil.TempDir("", "openview-test")
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}

	s.tempDir = safe.UnsafeNewPath(tempDir)
	s.cacheDir = s.tempDir.JoinUnsafe("cache")
	s.imageDir = s.tempDir.JoinUnsafe("images")

	err = os.Mkdir(s.cacheDir.String(), 0700)
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}
	err = os.Mkdir(s.imageDir.String(), 0700)
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}

	s.app, err = NewApplication(&Config{
		ResourceDir: safe.UnsafeNewPath("../dist"),
		CacheDir:    s.cacheDir,
		ImageDir:    s.imageDir,
	})
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}
}

func (s *AppSuite) TearDownSuite(c *C) {
	if !s.tempDir.IsEmpty() {
		os.RemoveAll(s.tempDir.String())
		s.tempDir = safe.Path{}
		s.cacheDir = safe.Path{}
		s.imageDir = safe.Path{}
	}
}

func (s *AppSuite) TestFavicon(c *C) {
	expectedBytes, err := ioutil.ReadFile("../dist/favicon.ico")
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}

	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		c.Fatalf("Error: %+v", errors.WithStack(err))
	}
	rr := httptest.NewRecorder()
	s.app.router.ServeHTTP(rr, req)

	c.Assert(rr.Code, Equals, http.StatusOK)
	c.Assert(rr.Body.Bytes(), DeepEquals, expectedBytes)
}

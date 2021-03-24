package ftp

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/sno6/gate-god/camera"

	"github.com/pkg/errors"
	_ "github.com/shurcooL/vfsgen"
	"goftp.io/server/v2"
)

type Server struct {
	logger  *log.Logger
	handler camera.MotionDetector
}

type Driver struct {
	handler camera.MotionDetector
}

func New(handler camera.MotionDetector) *Server {
	if handler == nil {
		panic("handler should not be nil")
	}

	return &Server{
		logger:  log.New(os.Stdout, "[FTP]: ", log.LstdFlags),
		handler: handler,
	}
}

func (d *Driver) Stat(ctx *server.Context, s string) (os.FileInfo, error) {
	return nil, nil
}

func (d *Driver) ListDir(ctx *server.Context, s string, fn func(os.FileInfo) error) error {
	return nil
}

func (d *Driver) DeleteDir(ctx *server.Context, s string) error {
	return nil
}

func (d *Driver) DeleteFile(ctx *server.Context, s string) error {
	return nil
}

func (d *Driver) Rename(ctx *server.Context, s1 string, s2 string) error {
	return nil
}

func (d *Driver) MakeDir(ctx *server.Context, s string) error {
	return nil
}

func (d *Driver) GetFile(ctx *server.Context, s string, i int64) (int64, io.ReadCloser, error) {
	return 0, ioutil.NopCloser(nil), nil
}

func (d *Driver) PutFile(ctx *server.Context, s string, r io.Reader, i int64) (int64, error) {
	err := d.handler.OnMotionDetection(r)
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (s *Server) Serve() error {
	s.logger.Println("Starting server on port 2121")

	opt := &server.Options{
		Name:   "gate-god",
		Driver: &Driver{handler: s.handler},
		Port:   2121,
		Auth:   &server.SimpleAuth{Name: "admin", Password: "password"},
		Perm:   server.NewSimplePerm("root", "root"),
		Logger: &server.StdLogger{},
	}

	ftp, err := server.NewServer(opt)
	if err != nil {
		return errors.Wrap(err, "ftp: error while serving")
	}

	return ftp.ListenAndServe()
}

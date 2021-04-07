package ftp

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	_ "github.com/shurcooL/vfsgen"
	"github.com/sno6/gate-god/camera"
	"github.com/sno6/gate-god/camera/batch"
	"goftp.io/server/v2"
)

const port = 2121

// Server is an FTP server implementation.
//
// But why FTP? I hear you ask..
//
// Because I have an old Hikvision camera and if I want to get
// motion detection events sent to me then it needs to be over FTP..
// however, an HTTP implementation should be similar and easy to do.
type Server struct {
	cfg     *Config
	logger  *log.Logger
	batcher *batch.FrameBatcher
}

type Config struct {
	User, Password string
}

type Driver struct {
	batcher *batch.FrameBatcher
}

func New(cfg *Config, batcher *batch.FrameBatcher) *Server {
	return &Server{
		cfg:     cfg,
		batcher: batcher,
		logger:  log.New(os.Stdout, "[FTP]: ", log.LstdFlags),
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

// PutFile.. the only method we care about - essentially just forwards the frame to the batcher.
func (d *Driver) PutFile(ctx *server.Context, s string, r io.Reader, i int64) (int64, error) {
	d.batcher.SendFrame(&camera.Frame{R: r, Name: s})
	return 1, nil
}

func (s *Server) Serve() error {
	s.logger.Printf("Starting server on port: %d\n", port)

	ftp, err := server.NewServer(&server.Options{
		Port:   port,
		Auth:   &server.SimpleAuth{Name: s.cfg.User, Password: s.cfg.Password},
		Perm:   server.NewSimplePerm("root", "root"),
		Driver: &Driver{batcher: s.batcher},
		Logger: &server.DiscardLogger{},
	})
	if err != nil {
		return errors.Wrap(err, "ftp: error while serving")
	}

	return ftp.ListenAndServe()
}

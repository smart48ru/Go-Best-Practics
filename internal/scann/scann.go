package scann

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type FileInfo interface {
	os.FileInfo
	Path() string
	Dir() string
}

func (fi fileInfo) Path() string {
	return fi.path
}
func (fi fileInfo) Dir() string {
	return fi.dir
}

type fileInfo struct {
	os.FileInfo
	path string
	dir  string
}

type scanner struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup
	dir       string // start dir
	resChan   chan fileInfo
	errChan   chan error
	depth     int64
	ext       string
	cDir      string // current work dir
}

type Scanner interface {
	ListDirectory(dir string, depth int64)
	FindFiles()
	DeIncDepth()
	IncDepth()
	ErrChan() chan error
	ResChan() chan fileInfo
	CurDir() string
	Depth() int64
	WG() *sync.WaitGroup
	Ctx() context.Context
	CtxCancel() context.CancelFunc
}

func (s *scanner) CurDir() string {
	return s.cDir
}

func (s *scanner) Ctx() context.Context {
	return s.ctx
}
func (s *scanner) CtxCancel() context.CancelFunc {
	return s.ctxCancel
}

func (s *scanner) ErrChan() chan error {
	return s.errChan
}

func (s *scanner) ResChan() chan fileInfo {
	return s.resChan
}

func (s *scanner) DeIncDepth() {
	atomic.AddInt64(&s.depth, -1)
}

func (s *scanner) WG() *sync.WaitGroup {
	return &s.wg
}

func (s *scanner) IncDepth() {
	var delta int64 = 2
	atomic.AddInt64(&s.depth, delta)
}

func (s *scanner) Depth() int64 {
	return s.depth
}

func (s *scanner) ListDirectory(dir string, depth int64) {
	defer s.WG().Done()
	if depth < 0 {
		return
	}
	select {
	case <-s.ctx.Done():
		return
	default:
		// time.Sleep(time.Second * 1)
		res, err := os.ReadDir(dir)
		if err != nil {
			s.errChan <- err
		}
		for _, entry := range res {
			path := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				s.cDir = dir
				log.Trace().Msgf("Recurse start ListDirectory in goroutine depth = %d", depth)
				s.WG().Add(1)
				go s.ListDirectory(path, depth-1)
			} else {
				info, err := entry.Info()
				if err != nil {
					s.errChan <- err
				}
				res := fileInfo{info, path, dir}
				if filepath.Ext(res.Name()) == s.ext {
					s.resChan <- res
				}
			}
		}
	}
}

func (s *scanner) FindFiles() {
	s.WG().Add(1)
	log.Trace().Msg("Starting ListDirectory in goroutine")
	go s.ListDirectory(s.dir, s.depth)
	s.WG().Wait()
	log.Trace().Msg("Done - Which WaitGroup")
	s.ctxCancel()
}

func New(timeout int, dir, ext string, depth int64) scanner {
	log.Trace().Msg("Make context WitchTimeOut")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)

	return scanner{ //nolint:exhaustivestruct
		ctx:       ctx,
		ctxCancel: cancel,
		wg:        sync.WaitGroup{},
		dir:       dir,
		resChan:   make(chan fileInfo),
		errChan:   make(chan error),
		depth:     depth,
		ext:       ext,
	}
}

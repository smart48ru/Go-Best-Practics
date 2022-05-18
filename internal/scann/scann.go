package scann

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
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
	ctx     context.Context
	dir     string
	resChan chan fileInfo
	errChan chan error
	depth   int64
	ext     string
	cDir    string
	check   map[string]struct{}
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
}

func (s *scanner) CurDir() string {
	return s.cDir
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

func (s *scanner) IncDepth() {
	atomic.AddInt64(&s.depth, 2)
}
func (s *scanner) Depth() int64 {
	return s.depth
}

func (s *scanner) ListDirectory(dir string, depth int64) {
	if depth < 0 {
		return
	}
	select {
	case <-s.ctx.Done():
		return
	default:
		time.Sleep(time.Second * 5)
		res, err := os.ReadDir(s.dir)
		if err != nil {
			s.errChan <- err
		}
		for _, entry := range res {
			path := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				s.cDir = dir
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
	go s.ListDirectory(s.dir, s.depth)
}

func New(ctx context.Context, dir, ext string, depth int64) scanner {
	return scanner{
		ctx:     ctx,
		dir:     dir,
		resChan: make(chan fileInfo),
		errChan: make(chan error),
		depth:   depth,
		ext:     ext,
		check:   make(map[string]struct{}),
	}
}

package osutil

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type Settings struct {
	initOnce sync.Once

	Folder string

	cmd *bufio.Reader
}

func (s *Settings) init() error {
	var err error
	s.initOnce.Do(func() {
		s.cmd = bufio.NewReader(os.Stdin)
	})
	return err
}

func (s *Settings) Open(name string) (io.ReadCloser, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

	rel := name
	if len(s.Folder) > 0 {
		rel = strings.Join([]string{s.Folder, name}, "/")
	}

	for {
		f, err := os.Open(rel)
		if os.IsNotExist(err) {
			println(`"` + rel + `" not found; drop file here:`)

			in, err := s.cmd.ReadString('\n')
			if err != nil {
				return nil, err
			}
			abs := strings.Replace(in[:len(in)-1], "\r", "", -1)

			if f, err = os.Open(abs); err != nil {
				println("Open err")
				return nil, err
			}
			f2, err := os.Create(rel)
			if err != nil {
				println("Create err")
				return nil, err
			}
			if _, err = io.Copy(f2, f); err != nil {
				f.Close()
				f2.Close()
				println("Copy err")
				return nil, err
			}
			f.Close()
			f2.Close()
			continue
		}
		return f, err
	}
}

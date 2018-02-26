package filedropper

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type Struct struct {
	initOnce sync.Once

	Folder string

	cmd *bufio.Reader
}

func (s *Struct) init() error {
	var err error
	s.initOnce.Do(func() {
		s.cmd = bufio.NewReader(os.Stdin)
	})
	return err
}

func (s *Struct) Open(name string) (io.ReadCloser, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

	for {
		f, err := os.Open(s.Folder + "/" + name)
		if os.IsNotExist(err) {
			println(`"` + s.Folder + `/` + name + `" not found; drop file here:`)

			dragNDrop, err := s.cmd.ReadString('\n')
			if err != nil {
				return nil, err
			}
			path := strings.Replace(dragNDrop, "\r", "", -1)

			if f, err = os.Open(path); err != nil {
				return nil, err
			}
			f2, err := os.Create(s.Folder + "/" + name)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(f2, f); err != nil {
				f.Close()
				f2.Close()
				return nil, err
			}
			f.Close()
			f2.Close()
			continue
		}

		return f, err
	}
}

package osutil

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	osutil Settings
)

// Var reads a variable stored in the standard settings.
// If not found, it asks for a value and stores it.
func Var(query string) (string, error) {
	return osutil.Var(query)
}

// Delete deletes the stored variable in the standard settings.
func Delete(query string) {
	osutil.Delete(query)
}

// Open opens a file with the standard settings.
// If not found, it asks for a file and copies it over.
func Open(name string) (io.ReadCloser, error) {
	return osutil.Open(name)
}

// Settings are the values osutil will use for functioning.
type Settings struct {
	initOnce sync.Once

	Folder string

	cmd cmdppt
	ppt map[string]string
}

func (s *Settings) init() error {
	var initErr error
	s.initOnce.Do(func() {
		s.cmd.Reader = bufio.NewReader(os.Stdin)
		s.ppt = make(map[string]string)
	})
	return initErr
}

// Var reads a variable stored in the settings.
// If not found, it asks for a value and stores it.
func (s *Settings) Var(query string) (string, error) {
	if err := s.init(); err != nil {
		return "", err
	}

	ans, found := s.ppt[query]
	if !found {
		print(query + ": ")
		in, err := s.cmd.ReadString('\n')
		if err != nil {
			return "", err
		}
		ans = in
		s.ppt[query] = ans
	}

	return ans, nil
}

// Delete deletes the stored variable in the settings.
func (s *Settings) Delete(query string) {
	delete(s.ppt, query)
}

// Open opens a file with the settings.
// If not found, it asks for a file and copies it over.
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

			abs, err := s.cmd.ReadString('\n')
			if err != nil {
				return nil, err
			}

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

type cmdppt struct {
	*bufio.Reader
}

func (cmd cmdppt) ReadString(delim byte) (string, error) {
	in, err := cmd.Reader.ReadString(delim)
	if err != nil {
		return "", err
	}
	return strings.Replace(in[:len(in)-1], "\r", "", -1), nil
}

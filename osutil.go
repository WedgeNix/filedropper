package osutil

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	osutil Settings
)

// Alert alters the user, waiting for input on the standard settings.
func Alert(query string) error {
	return osutil.Alert(query)
}

// Var reads a variable stored in the standard settings.
// If not found, it asks for a value and stores it.
func Var(query string, ptr interface{}) error {
	return osutil.Var(query, ptr)
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

	Folder   string
	Location *time.Location

	cmd cmdppt
	ppt map[string]string
}

func (s *Settings) init() error {
	var initErr error
	s.initOnce.Do(func() {
		if s.Location == nil {
			s.Location = time.Local
		}
		s.cmd.Reader = bufio.NewReader(os.Stdin)
		s.ppt = make(map[string]string)
	})
	return initErr
}

// Alert alters the user, waiting for input on the settings.
func (s *Settings) Alert(query string) error {
	if err := s.init(); err != nil {
		return err
	}

	print(query)
	_, err := s.cmd.Reader.ReadString('\n')
	return err
}

// Var reads a variable stored in the settings.
// If not found, it asks for a value and stores it.
func (s *Settings) Var(query string, ptr interface{}) error {
	if err := s.init(); err != nil {
		return err
	}

	ans, found := s.ppt[query]
	if !found {
		print(query + ": ")
		in, err := s.cmd.ReadString('\n')
		if err != nil {
			return err
		}
		ans = in
	}

Types:
	switch v := ptr.(type) {
	case *time.Time:
		for _, layout := range [...]string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
			time.RFC3339Nano,
			time.Kitchen,
			time.Stamp,
			time.StampMilli,
			time.StampMicro,
			time.StampNano,
			"1/06",
			"1/2006",
			"01/06",
			"01/2006",
			"1-06",
			"1-2006",
			"01-06",
			"01-2006",
			"1/_2/06",
			"1/_2/2006",
			"01/02/06",
			"01/02/2006",
			"1-_2-06",
			"1-_2-2006",
			"01-02-06",
			"01-02-2006",
		} {
			if t, err := time.ParseInLocation(layout, ans, s.Location); err == nil {
				*v = t
				break Types
			}
		}
		return errors.New("bad time format")
	case *url.URL:
		u, err := url.Parse(ans)
		if err != nil {
			return err
		}
		*v = *u
	case *float64:
		f, err := strconv.ParseFloat(ans, 64)
		if err != nil {
			return err
		}
		*v = f
	case *int:
		n, err := strconv.Atoi(ans)
		if err != nil {
			return err
		}
		*v = n
	case *string:
		*v = ans
	default:
		if err := json.Unmarshal([]byte(ans), ptr); err != nil {
			return err
		}
	}

	s.ppt[query] = ans
	return nil
}

// Delete deletes the stored variable in the settings.
func (s *Settings) Delete(query string) error {
	if err := s.init(); err != nil {
		return err
	}

	delete(s.ppt, query)
	return nil
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
		if !os.IsNotExist(err) {
			return f, err
		}
		println(`"` + rel + `" not found; drop file here:`)

		abs, err := s.cmd.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if f, err = os.Open(abs); err != nil {
			return nil, err
		}
		var f2 *os.File
		for {
			f2, err = os.Create(rel)
			if err == nil {
				break
			}
			if err := os.MkdirAll(filepath.Dir(rel), os.ModePerm); err != nil {
				return nil, err
			}
		}
		if _, err = io.Copy(f2, f); err != nil {
			f.Close()
			f2.Close()
			return nil, err
		}
		f.Close()
		f2.Close()
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

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
	Folder   string
	Location *time.Location

	osutil Settings
)

// Alert alters the user, waiting for input on the standard settings.
func Alert(query string) {
	osutil.Alert(query)
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
func Open(name string) *os.File {
	return osutil.Open(name)
}

// Create creates the named file and its respective directories.
func Create(name string) *os.File {
	return osutil.Create(name)
}

// Settings are the values osutil will use for functioning.
type Settings struct {
	initOnce sync.Once

	Folder   string
	Location *time.Location

	cmd cmdppt
	ppt map[string]string
}

func (s *Settings) init() {
	s.initOnce.Do(func() {
		if len(s.Folder) == 0 {
			s.Folder = Folder
		}
		if s.Location == nil {
			if Location == nil {
				Location = time.Local
			}
			s.Location = Location
		}
		s.cmd.Reader = bufio.NewReader(os.Stdin)
		s.ppt = make(map[string]string)
	})
}

// Alert alters the user, waiting for input on the settings.
func (s *Settings) Alert(query string) {
	s.init()
	print(query)
	s.cmd.Reader.ReadString('\n')
}

// Var reads a variable stored in the settings.
// If not found, it asks for a value and stores it.
func (s *Settings) Var(query string, ptr interface{}) error {
	s.init()

	ans, found := s.ppt[query]
	if !found {
		print(query + ": ")
		in := s.cmd.ReadString('\n')
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
func (s *Settings) Delete(query string) {
	s.init()
	delete(s.ppt, query)
}

// Open opens a file with the settings.
// If not found, it asks for a file and copies it over.
func (s *Settings) Open(name string) *os.File {
	s.init()

	path := s.Path(name)
	for {
		f, err := os.Open(path)
		if err == nil {
			return f
		}
		s.Alert(err.Error() + "; press enter to retry")
	}
}

// Path checks for a file with the settings.
// If not found, it asks for a file and copies it over.
// The actual relative path is returned.
func (s *Settings) Path(name string) string {
	s.init()

	rel := name
	if len(s.Folder) > 0 {
		rel = strings.Join([]string{s.Folder, name}, "/")
	}

	for {
		if _, err := os.Stat(rel); !os.IsNotExist(err) {
			return rel
		}
		println(`"` + rel + `" not found; drop file here:`)

		var f *os.File
		for {
			abs := s.cmd.ReadString('\n')
			if len(abs) > 2 && abs[0] == '"' && abs[len(abs)-1] == '"' {
				abs = abs[1 : len(abs)-1]
			}
			var err error
			if f, err = os.Open(abs); err == nil {
				break
			}
			println(err.Error() + "; drop file here:")
		}
		f2 := s.Create(name)
		if _, err := io.Copy(f2, f); err != nil {
			panic(err)
		}
		f.Close()
		f2.Close()
	}
}

// Create creates the named file and its respective directories.
func (s *Settings) Create(name string) *os.File {
	s.init()

	rel := name
	if len(s.Folder) > 0 {
		rel = strings.Join([]string{s.Folder, name}, "/")
	}

	for {
		f, err := os.Create(rel)
		if err == nil {
			return f
		}
		if err = os.MkdirAll(filepath.Dir(rel), os.ModePerm); err != nil {
			panic(err)
		}
	}
}

type cmdppt struct {
	*bufio.Reader
}

func (cmd cmdppt) ReadString(delim byte) string {
	in, err := cmd.Reader.ReadString(delim)
	if err != nil {
		panic(err)
	}
	return strings.Replace(in[:len(in)-1], "\r", "", -1)
}

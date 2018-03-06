package osutil

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
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

// Check checks for a file with the standard settings.
// If not found, it asks for a file and copies it over.
// The name is returned.
func Check(name string) string {
	return osutil.Check(name)
}

// CheckDir returns all file names from a folder.
func CheckDir(dir string) []string {
	return osutil.CheckDir(dir)
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

// Copy copies a file to a directory with the standard settings.
// The new path is returned.
func Copy(dir string) string {
	return osutil.Copy(dir)
}

// Settings are the values osutil will use for functioning.
type Settings struct {
	initOnce sync.Once

	Location *time.Location

	cmd cmdppt
	ppt map[string]string
}

func (s *Settings) init() {
	s.initOnce.Do(func() {
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
	_, err := s.cmd.Reader.ReadString('\n')
	if err != nil {
		panic("osutil: " + err.Error())
	}
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
		return errors.New("osutil: bad time format")
	case *url.URL:
		u, err := url.Parse(ans)
		if err != nil {
			return fmt.Errorf("osutil: %v", err)
		}
		*v = *u
	case *float64:
		f, err := strconv.ParseFloat(ans, 64)
		if err != nil {
			return fmt.Errorf("osutil: %v", err)
		}
		*v = f
	case *int:
		n, err := strconv.Atoi(ans)
		if err != nil {
			return fmt.Errorf("osutil: %v", err)
		}
		*v = n
	case *string:
		*v = ans
	default:
		if err := json.Unmarshal([]byte(ans), ptr); err != nil {
			return fmt.Errorf("osutil: %v", err)
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

	path := s.Check(name)
	for {
		f, err := os.Open(path)
		if err == nil {
			return f
		}
		s.Alert(err.Error() + "; press enter to retry")
	}
}

// CheckDir returns all file names from a folder.
func (s *Settings) CheckDir(dir string) []string {
	dir = fixDir(dir)
	for {
		infos, err := ioutil.ReadDir(dir)
		if err == nil {
			var files []string
			for _, info := range infos {
				if !info.IsDir() {
					files = append(files, dir+info.Name())
				}
			}
			return files
		}
		if os.IsNotExist(err) {
			s.MkDir(dir)
			s.Alert("dir '" + dir + "' not found; created, press enter to retry")
		} else {
			s.Alert(err.Error() + "; press enter to retry")
		}
	}
}

// Check checks for a file with the settings.
// If not found, it asks for a file and copies it over.
// The name is returned.
func (s *Settings) Check(name string) string {
	s.init()

	for {
		if _, err := os.Stat(name); !os.IsNotExist(err) {
			return name
		}
		println(`"` + name + `" not found; drop file here:`)

		var f *os.File
		for {
			path := s.cmd.ReadString('\n')
			if len(path) > 2 && path[0] == '"' && path[len(path)-1] == '"' {
				path = path[1 : len(path)-1]
			}
			var err error
			if f, err = os.Open(path); err == nil {
				break
			}
			println(err.Error() + "; drop file here:")
		}
		f2 := s.Create(name)
		if _, err := io.Copy(f2, f); err != nil {
			panic("osutil: " + err.Error())
		}
		f.Close()
		f2.Close()
	}
}

// Copy copies a file to a directory with the settings.
// The new path is returned.
func (s *Settings) Copy(dir string) string {
	s.init()

	println(`drop file here:`)

	var new string
	var f *os.File
	for {
		path := s.cmd.ReadString('\n')
		if len(path) > 2 && path[0] == '"' && path[len(path)-1] == '"' {
			path = path[1 : len(path)-1]
		}
		var err error
		if f, err = os.Open(path); err == nil {
			defer f.Close()
			new = fixDir(dir) + filepath.Base(path)
			break
		}
		println(err.Error() + "; drop file here:")
	}
	f2 := s.Create(new)
	defer f2.Close()
	if _, err := io.Copy(f2, f); err != nil {
		panic("osutil: " + err.Error())
	}

	return new
}

func fixDir(dir string) string {
	if len(dir) == 0 {
		return "./"
	}
	if dir[0] == '\\' || dir[0] == '/' {
		dir = dir[1:]
	}
	if dir[len(dir)-1] != '\\' && dir[len(dir)-1] != '/' {
		dir += `/`
	}
	return dir
}

// Create creates the named file and its respective directories.
func (s *Settings) Create(name string) *os.File {
	s.init()

	for {
		f, err := os.Create(name)
		switch {
		case err == nil:
			return f
		case os.IsNotExist(err):
			s.MkDir(filepath.Dir(name))
		default:
			panic("osutil: " + err.Error())
		}
	}
}

// MkDir makes all nonexistent directories.
func (s *Settings) MkDir(path string) {
	s.init()

	for {
		err := os.MkdirAll(path, os.ModePerm)
		if err == nil {
			break
		}
		s.Alert(err.Error() + "; press enter to retry")
	}
}

type cmdppt struct {
	*bufio.Reader
}

func (cmd cmdppt) ReadString(delim byte) string {
	in, err := cmd.Reader.ReadString(delim)
	if err != nil {
		panic("osutil: " + err.Error())
	}
	return strings.Replace(in[:len(in)-1], "\r", "", -1)
}

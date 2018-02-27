package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/WedgeNix/osutil"
)

func main() {
	os := osutil.Settings{Folder: "dump"}

	var name string
	if err := os.Var("file name", &name); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	var t time.Time
	err = os.Var("Date", &t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)

	os.Alert("Press enter to continue")
}

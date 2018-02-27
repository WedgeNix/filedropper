package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/WedgeNix/osutil"
)

func main() {
	os := osutil.Settings{}

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

	type Color int
	var x struct {
		Name  string
		Color Color
	}
	err = os.Var("Name & Color (JSON)", &x)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(x)
}

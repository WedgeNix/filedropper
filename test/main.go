package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/WedgeNix/osutil"
)

func main() {
	os := osutil.Settings{}

	in, err := os.Var("file name")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	in2, err := os.Var("file name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(in2)
}

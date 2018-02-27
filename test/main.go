package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/WedgeNix/osutil"
)

func main() {
	in, err := osutil.Prompt("file name")
	if err != nil {
		log.Fatal(err)
	}

	f, err := osutil.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	in2, err := osutil.Prompt("file name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(in2)
}

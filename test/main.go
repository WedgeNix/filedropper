package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/WedgeNix/osutil"
)

func main() {
	os := osutil.Settings{}
	f, err := os.Open("Something.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

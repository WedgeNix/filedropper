package main

import (
	"github.com/WedgeNix/osutil"
)

func main() {
	path := osutil.Copy("bin_x")
	osutil.Alert(path)
	// var name string
	// if err := osutil.Var("file name", &name); err != nil {
	// 	log.Fatal(err)
	// }

	// f := osutil.Open(name)
	// defer f.Close()

	// b, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	// var t time.Time
	// err = osutil.Var("Date", &t)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(t)

	// osutil.Alert("Press enter to continue")
}

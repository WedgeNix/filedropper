package main

import "path/filepath"

func main() {
	println(filepath.Dir("/x/Mango/"))
	// names := osutil.CheckDir("/x/mango/")
	// fmt.Println(names)
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

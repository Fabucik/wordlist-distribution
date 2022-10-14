package checkerr

import "log"

func CheckError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

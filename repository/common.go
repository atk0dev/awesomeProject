package repository

import "log"

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


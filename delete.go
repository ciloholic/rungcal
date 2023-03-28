package rungcal

import "log"

type DeleteOption struct {
	Option
}

func Delete(deleteOption DeleteOption) int {
	log.Println(deleteOption)

	return 0
}

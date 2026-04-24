package utils

import "log"

func PrintError(message string, err error) {
	log.Printf(message+": %v \n", err)
}

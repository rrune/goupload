package util

import (
	"errors"
	"log"
	"os"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"

func Check(err error) (r bool) {
	if err != nil {
		r = true
	}
	return
}

func CheckWLogs(err error) (r bool) {
	if err != nil {
		r = true
		log.Println(string(colorYellow), err, string(colorReset))
	}
	return
}

func CheckPanic(err error) {
	if err != nil {
		f, err2 := os.OpenFile("../data/err.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err2 != nil {
			log.Fatal(err, err2)
		}
		defer f.Close()
		_, err2 = f.Write([]byte(err.Error()))
		if err2 != nil {
			log.Fatal(err, err2)
		}
		panic(err)
	}
}

func EnsureUniqueFilenames(path string, filename string) (r string, err error) {
	r = filename
	for {
		_, err = os.Stat(path + r)

		if errors.Is(err, os.ErrNotExist) {
			err = nil
			return
		} else if Check(err) {
			return
		} else {
			r = "_" + r
		}
	}
}

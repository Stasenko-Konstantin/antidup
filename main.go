package main

import (
	"fmt"
	"github.com/Stasenko-Konstantin/phash"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type pic struct {
	name string
	hash string
}

func findHash(pic string) (string, error) {
	img, err := os.Open(pic)
	if err != nil {
		return "", err
	}
	defer img.Close()

	hash, err := phash.GetHash(img)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func findDuplicates(pics []pic) (string, error) {
	var (
		r    string
		dups []map[string]string
	)

	for n, p := range pics {
		var s []pic
		if len(pics) == 1 {
			s = pics[n:]
		} else {
			s = pics[n+1:]
		}
		for _, comp := range s {
			dup := make(map[string]string)
			distance := phash.GetDistance(p.hash, comp.hash)
			if distance < 3 {
				dup[p.name] = comp.name
			}
			if len(dup) > 0 {
				dups = append(dups, dup)
			}
		}
	}
	for _, dup := range dups {
		for k, e := range dup {
			r += k + ": " + e + "\n"
		}
	}
	return r, nil
}

func topErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	var pics []pic
	for _, file := range files {
		name := file.Name()
		if strings.Contains(name, ".") {
			format := strings.Split(name, ".")[1]
			if format == "png" || format == "jpg" || format == "jpeg" {
				hash, _ := findHash(name)
				topErr(err)
				pics = append(pics, pic{name, hash})
			}
		}
	}

	r, err := findDuplicates(pics)
	topErr(err)
	if len(pics) < 2 || r == "" {
		fmt.Println("no duplicates found")
	} else {
		fmt.Println(r[:len(r)-1])
	}
}

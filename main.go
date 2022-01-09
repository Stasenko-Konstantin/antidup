package main

import (
	"fmt"
	"github.com/Stasenko-Konstantin/phash"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func compare(pic string, compF string) (int, error) {
	comp, err := findHash(compF)
	if err != nil {
		return 0, err
	}
	return phash.GetDistance(pic, comp), nil
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

func findDuplicates(pics []string) (string, error) {
	var (
		r    string
		dups []map[string]string
	)

	for n, pic := range pics {
		hash, err := findHash(pic)
		if err != nil {
			return "", err
		}
		if len(pics) > n+1 {
			pics = append(pics[:n], pics[n+1:]...)
		} else if len(pics) == 1 {
			break
		}
		for _, comp := range pics {
			dup := make(map[string]string)
			distance, err := compare(hash, comp)
			if err != nil {
				return "", err
			}
			if distance < 3 {
				dup[pic] = comp
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

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	var pics []string
	for _, file := range files {
		name := file.Name()
		if strings.Contains(name, ".") {
			format := strings.Split(name, ".")[1]
			if format == "png" || format == "jpg" {
				pics = append(pics, name)
			}
		}
	}

	r, err := findDuplicates(pics)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(r)
}

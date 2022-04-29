package main

import (
	"antidup/help"
	"fmt"
	"github.com/Stasenko-Konstantin/phash"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
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

	fmt.Println("calculation...")
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
	fmt.Fprintf(os.Stdout, "\r\033[K")
	for _, dup := range dups {
		for k, e := range dup {
			r += "\t" + k + "  --  " + e + "\n"
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

func check(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var pics []pic
	for _, file := range files {
		name := file.Name()
		if strings.Contains(name, ".") {
			format := help.GetFormat(name)
			if format == "png" || format == "jpg" || format == "jpeg" {
				hash, err := findHash(name)
				topErr(err)
				size, err := help.GetSize(name)
				topErr(err)
				pics = append(pics, pic{name + ", " + size, hash})
			}
		}
	}

	r, err := findDuplicates(pics)
	topErr(err)
	if len(pics) < 2 || r == "" {
		fmt.Println(dir + ": no duplicates found")
	} else {
		fmt.Println(dir + ":\n" + r[:len(r)-1])
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "antidup",
		Usage: "to find duplicates of photos ",
		Commands: []*cli.Command{
			{
				Name:    "rec",
				Aliases: []string{"r"},
				Usage:   "shallow recursive check of the subfolders of the current directory",
				Action: func(c *cli.Context) error {
					files, err := ioutil.ReadDir("./")
					if err != nil {
						return err
					}
					err = check("./")
					if err != nil {
						return err
					}
					for _, file := range files {
						if file.IsDir() {
							err := check(file.Name())
							if err != nil {
								return err
							}
						}
					}
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			args := os.Args
			if len(args) > 1 {
				for _, dir := range args[1:] {
					err := check(dir)
					if err != nil {
						return err
					}
				}
				return nil
			} else {
				err := check("./")
				return err
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

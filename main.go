package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func compare(pic [16][4]int, compF string) (float32, error) {
	comp, err := findHistogram(compF)
	if err != nil {
		return 0, err
	}

	var percent float32
	fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	for i, x := range pic {
		fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
	}

	fmt.Printf("\n%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	for i, x := range comp {
		fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
	}

	return percent, nil
}

func findHistogram(pic string) ([16][4]int, error) {
	var histogram [16][4]int
	picReader, err := os.ReadFile(pic)
	if err != nil {
		return histogram, err
	}
	picm, _, err := image.Decode(bytes.NewBuffer(picReader)) // hmmmmm
	if err != nil {
		return histogram, err
	}
	bounds := picm.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := picm.At(x, y).RGBA()
			histogram[r>>12][0]++
			histogram[g>>12][1]++
			histogram[b>>12][2]++
			histogram[a>>12][3]++
		}
	}
	return histogram, nil
}

func findDuplicates(pics []string) (string, error) {
	var (
		r    string
		dups []map[string]string
	)

	for n, pic := range pics {
		histogram, err := findHistogram(pic)
		if err != nil {
			return "", err
		}
		for _, comp := range pics[n:] {
			dup := make(map[string]string)
			percent, err := compare(histogram, comp)
			if err != nil {
				return "", err
			}
			if percent < 0.1 {
				dup[pic] = comp
			}
		}
	}
	for _, dup := range dups {
		for k, e := range dup {
			r += k + ": " + e + "\n"
		}
		r += "\n"
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

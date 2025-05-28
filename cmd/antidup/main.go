package main

import (
	"context"
	"fmt"
	"github.com/Stasenko-Konstantin/phash"
	"github.com/urfave/cli/v3"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type recursive int

const (
	none recursive = iota
	segmented
	flat
)

var (
	quiet bool
	rm    bool
	rec   recursive
	deep  uint32
	path  string
)

func main() {
	(&cli.Command{
		Name:      "antidup",
		Usage:     "to find duplicates of photos",
		UsageText: "antidup [GLOBAL OPTIONS]",
		Action:    mainAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Value:   false,
				Usage:   "do not print anything",
				Action: func(ctx context.Context, command *cli.Command, q bool) error {
					quiet = q
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "rm",
				Value: false,
				Usage: "remove duplicates",
				Action: func(ctx context.Context, command *cli.Command, r bool) error {
					rm = r
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Value:   "none",
				Usage:   "recursively find duplicates (none/flat/segmented)",
				Action: func(ctx context.Context, command *cli.Command, r string) error {
					r = strings.TrimSpace(r)
					r = strings.ToLower(r)
					switch r {
					case "none":
						rec = none
					case "segmented":
						rec = segmented
					case "flat":
						rec = flat
					default:
						return fmt.Errorf("invalid recursive option: %s", r)
					}
					return nil
				},
			},
			&cli.UintFlag{
				Name:    "deep",
				Aliases: []string{"d"},
				Value:   0,
				Usage:   "deep finding duplicates",
				Action: func(ctx context.Context, command *cli.Command, d uint) error {
					deep = uint32(d)
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   ".",
				Usage:   "path to find duplicates",
				Action: func(ctx context.Context, command *cli.Command, p string) error {
					if path == "" {
						path = "."
					}
					path = p
					return nil
				},
			},
		},
	}).Run(context.Background(), os.Args)
}

func mainAction(_ context.Context, _ *cli.Command) error {
	if rec == none {
		deep = 1
	}
	if path == "" {
		path = "."
	}
	return check(path)
}

type pic struct {
	name string
	hash string
	size int
}

type t struct {
	p1 *pic
	p2 *pic
}

func (p *pic) equal(that *pic) bool {
	return p.hash == that.hash
}

func (p *pic) sizestr() string {
	const kb = 1024
	var res string
	switch {
	case p.size < kb:
		res = fmt.Sprintf("%db", p.size)
	case p.size < int(math.Pow(kb, 2)):
		res = fmt.Sprintf("%dkb", p.size/kb)
	case p.size < int(math.Pow(kb, 3)):
		res = fmt.Sprintf("%dmb", p.size/kb/kb)
	case p.size < int(math.Pow(kb, 4)):
		res = fmt.Sprintf("%dgb", p.size/kb/kb/kb)
	default:
		res = fmt.Sprintf("%db", p.size)
	}
	return res
}

func check(path string) error {
	isflat := rec == flat
	files := make(map[string][]os.DirEntry)
	err := mkfileidx(files, path, path, isflat, deep)
	if err != nil {
		return err
	}

	fmt.Println("calculation...")

	pics := make(map[string][]*pic)
	for k, v := range files {
		ps := make([]*pic, len(v))
		for i, f := range v {
			name := f.Name()
			if !quiet {
				fmt.Println(name)
			}
			hash, err := phash.Hash(filepath.Join(k, name))
			if err != nil {
				return err
			}
			info, err := f.Info()
			if err != nil {
				return err
			}
			ps[i] = &pic{name, string(hash), int(info.Size())}
		}
		pics[k] = ps
	}

	for path, ps := range pics {
		if err := processpics(path, ps); err != nil {
			return err
		}
	}

	return nil
}

func mkfileidx(res map[string][]os.DirEntry, root, path string, isflat bool, deep uint32) error {
	files, err := readdir(path)
	if err != nil {
		return err
	}
	if len(files) < 2 {
		fmt.Printf("no duplicates found in: %s\n", root)
		return nil
	}

	if deep == 0 {
		res[path] = files
	} else {
		if !isflat {
			res[path] = files
		} else if _, ok := res[path]; !ok && isflat {
			res[root] = files
		} else {
			res[path] = append(res[path], files...)
		}
		d, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, f := range d {
			if !f.IsDir() {
				continue
			}
			err := mkfileidx(res, path, filepath.Join(path, f.Name()), isflat, deep-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readdir(path string) ([]os.DirEntry, error) {
	formats := []string{".png", ".jpg", ".jpeg"}
	es, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := make([]os.DirEntry, 0, len(es)/2)
	for _, e := range es {
		if e.IsDir() {
			continue
		}
		if slices.Contains(formats, filepath.Ext(e.Name())) {
			result = append(result, e)
		}
	}
	return result, nil
}

func processpics(path string, pics []*pic) error {
	res := duplicates(pics)
	if len(res) == 0 {
		fmt.Printf("no duplicates found in: %s\n", path)
		return nil
	}
	var s strings.Builder
	for _, t := range res {
		if rm {
			var d string
			if t.p1.size > t.p2.size {
				d = t.p2.name
			} else {
				d = t.p1.name
			}
			if err := os.Remove(d); err != nil {
				return err
			}
		}
		s.WriteString(
			fmt.Sprintf("\t%s, %s -- %s. %s\n",
				t.p1.name, t.p1.sizestr(), t.p2.name, t.p2.sizestr(),
			),
		)
	}
	fmt.Println(s.String())
	return nil
}

func duplicates(pics []*pic) []*t {
	var (
		res  []*t
		dups []map[*pic]*pic
	)
	i := 0
	for _, p := range pics {
		if p == nil {
			continue
		}
		var s []*pic
		if len(pics) == 0 {
			s = pics[i:]
		} else {
			s = pics[i+1:]
		}
		for _, comp := range s {
			if comp == nil {
				continue
			}
			dup := make(map[*pic]*pic)
			distance := phash.Distance([]byte(p.hash), []byte(comp.hash))
			if distance < phash.MinDistance {
				dup[p] = comp
			}
			if len(dup) > 0 {
				dups = append(dups, dup)
			}
		}
		i += 1
	}
	for _, d := range dups {
		for k, v := range d {
			res = append(res, &t{k, v})
		}
	}
	return res
}

# antidup
to find duplicates of photos (.png, .jpg and .jpeg)                                
based on [phash algorithm](https://www.phash.org/docs/pubs/thesis_zauner.pdf)

![](example.png)

### old todo
- [x] loading ~~animation~~ message
- [x] display of image size in mb/KB/etc
- [x] analysis of the selected directory
- [x] delete duplicates
- [ ] fix image reading

### current todo
- [ ] rewrite it in golang

### get

requires [golang](https://go.dev/doc/install) for building executable file

```bash
git clone git@github.com:Stasenko-Konstantin/antidup.git 
cd antidup-
./install.sh   # requires sudo for cp executable file to /bin
               # reopen terminal
antidup -h                  
```

### usage

```bash
NAME:
   antidup - to find duplicates of photos

USAGE:
   antidup [GLOBAL OPTIONS]

GLOBAL OPTIONS:
   --quiet, -q                    do not print anything (default: false)
   --rm                           remove duplicates (default: false)
   --recursive string, -r string  recursively find duplicates (default: "none")
   --deep uint, -d uint           deep finding duplicates (default: 0)
   --path string, -p string       path to find duplicates (default: ".")
   --help, -h                     show help

```

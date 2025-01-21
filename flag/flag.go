package flag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Address *string
	Dir     *string
)

var restrictedDirs = []string{"go.mod", "flag", "handlers", "utils", "triple-s"}

func MyFlags() error {
	Address = flag.String("port", "8080", "HTTP network address")
	Dir = flag.String("dir", "data", "base dir")
	flag.Usage = usage
	flag.Parse()

	cleanedDir := filepath.Clean(*Dir)

	for _, res := range restrictedDirs {
		if strings.Contains(cleanedDir, res) {
			return fmt.Errorf("'%s' contains a restricted path and cannot be used as the base", *Dir)
		}
	}

	if strings.HasPrefix(cleanedDir, ".") || strings.Contains(cleanedDir, "..") || cleanedDir == "/" {
		return fmt.Errorf("'%s' is an invalid or restricted path and cannot be used", *Dir)
	}

	if cleanedDir == filepath.FromSlash("triple-s") {
		return fmt.Errorf("'%s' is a restricted directory and cannot be used as the base", *Dir)
	}

	dirInfo, err := os.Stat(cleanedDir)
	if err == nil && !dirInfo.IsDir() {
		return fmt.Errorf("'%s' exists as a file and cannot be used as a directory", *Dir)
	}

	return nil
}

func usage() {
	fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
}

package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateSrcAndDstPaths(src, dst string) (string, string, string) {
	if src == "" {
		return "", "", "Provide a file to encrypt or decrypt: \n $ znox <src_file> <dest_dir>"
	}

	// source
	srcStat, err := os.Stat(src)
	if err != nil {
		return "", "", fmt.Sprintf("source error: %s", err)
	}

	if srcStat.IsDir() {
		return "", "", "source is not a valid file"
	}

	// destination
	if dst == "" {
		dst = filepath.Dir(src) + filepath.Base(src) + ".enc"
	} else {
		dstStat, err := os.Stat(src)
		if err != nil {
			return "", "", fmt.Sprintf("destination error: %s", err)
		}
		if !dstStat.IsDir() {
			return "", "", "destination is not a valid directory"
		}
	}

	fmt.Printf("src=%s dst=%s\n", src, dst)
	return src, dst, ""
}

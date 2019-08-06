package utils


import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func Md5sum3(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}

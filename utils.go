package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/kr/s3/s3util"
	"io"
	"os"
	"os/user"
	"strings"
)

func envDefined(name string) bool {
	result := os.Getenv(name)
	return len(result) > 0
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func expandPath(path string) string {
	if path[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir + "/"
		new_path := strings.Replace(path, "~/", dir, 1)

		return new_path
	} else {
		return path
	}
}

func sha1Checksum(buffer string) string {
	h := sha1.New()
	io.WriteString(h, buffer)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func fatal(message string) {
	terminate(message, 1)
}

func terminate(message string, exit_code int) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(exit_code)
}

func terminateWithError(err error, exit_code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exit_code)
}

func isUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func s3Url(bucket string, filename string) string {
	format := "https://s3.amazonaws.com/%s/%s"
	url := fmt.Sprintf(format, bucket, filename)

	return url
}

func open(s string) (io.ReadCloser, error) {
	if isUrl(s) {
		return s3util.Open(s, nil)
	}
	return os.Open(s)
}

func create(s string) (io.WriteCloser, error) {
	if isUrl(s) {
		return s3util.Create(s, nil, nil)
	}
	return os.Create(s)
}

func transfer(file string, url string, fail_status int) {
	r, err := open(file)
	if err != nil {
		terminateWithError(err, fail_status)
	}

	w, err := create(url)
	if err != nil {
		terminateWithError(err, fail_status)
	}

	_, err = io.Copy(w, r)
	if err != nil {
		terminateWithError(err, fail_status)
	}

	err = w.Close()
	if err != nil {
		terminateWithError(err, fail_status)
	}
}

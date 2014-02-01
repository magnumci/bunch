package main

import(
  "testing"
  "os"
  "github.com/stretchr/testify/assert"
)

func Test_envDefined(t *testing.T) {
  os.Setenv("FOO", "1")
  os.Setenv("BAR", "")

  assert.Equal(t, envDefined("FOO"), true)
  assert.Equal(t, envDefined("BAR"), false)
}

func Test_fileExists(t *testing.T) {
  assert.Equal(t, fileExists("/tmp"), true)
  assert.Equal(t, fileExists("/tmp/foobar"), false)
}

func Test_expandPath(t *testing.T) {
  assert.Equal(t, expandPath("/path"), "/path")
  assert.Equal(t, expandPath("./path"), "./path")
  assert.Equal(t, expandPath("~/path"), os.Getenv("HOME") + "/path")
}

func Test_sha1Checksum(t *testing.T) {
  assert.Equal(t, sha1Checksum("foobar"), "8843d7f92416211de9ebb963ff4ce28125932878")
}

func Test_isUrl(t *testing.T) {
  assert.Equal(t, isUrl("foobar"), false)
  assert.Equal(t, isUrl("/path"), false)
  assert.Equal(t, isUrl("foobar.com"), false)
  assert.Equal(t, isUrl("ftp://foobar.com"), false)
  assert.Equal(t, isUrl("http://foobar.com"), true)
  assert.Equal(t, isUrl("https://foobar.com"), true)
}

func Test_s3Url(t *testing.T) {
  assert.Equal(t, s3Url("foo", "bar"), "https://s3.amazonaws.com/foo/bar")
}
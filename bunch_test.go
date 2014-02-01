package main

import(
  "testing"
  "os"
  "os/exec"
  "github.com/kr/s3/s3util"
  "github.com/stretchr/testify/assert"
)

func Test_loadS3CredentialsNotSet(t *testing.T) {
  loadS3Credentials()

  assert.Equal(t, options.S3Key, "")
  assert.Equal(t, options.S3Secret, "")
  assert.Equal(t, options.S3Bucket, "")
}

func Test_loadS3CredentialsSet(t *testing.T) {
  os.Setenv("S3_KEY", "key")
  os.Setenv("S3_SECRET", "secret")
  os.Setenv("S3_BUCKET", "bucket")

  loadS3Credentials()

  assert.Equal(t, options.S3Key, "key")
  assert.Equal(t, options.S3Secret, "secret")
  assert.Equal(t, options.S3Bucket, "bucket")
}

func Test_setS3Credentials(t *testing.T) {
  os.Setenv("S3_KEY", "key")
  os.Setenv("S3_SECRET", "secret")
  os.Setenv("S3_BUCKET", "bucket")

  loadS3Credentials()
  setS3Credentials()

  assert.Equal(t, s3util.DefaultConfig.AccessKey, "key")
  assert.Equal(t, s3util.DefaultConfig.SecretKey, "secret")
}

func Test_expandPathArguments(t *testing.T) {
  options.Path = "~/path"
  options.Manifest = "~/manifest"

  expandPathArguments()

  assert.Equal(t, options.Path, os.Getenv("HOME") + "/path")
  assert.Equal(t, options.Manifest, os.Getenv("HOME") + "/manifest")
}

func Test_extractDirectory(t *testing.T) {
  assert.Equal(t, extract("archive", "/tmp"), false)

  exec.Command("rm", "-rf", "/tmp/out").Start()
  assert.Equal(t, extract("archive", "/tmp/out"), false)

  exec.Command("rm", "-rf", "/tmp/out").Start()
  assert.Equal(t, extract("./test/archive.tar.gz", "/tmp/out"), true)
}

func Test_writeCacheFile(t *testing.T) {
  os.Remove("/tmp/.bunch")

  assert.Equal(t, writeCacheFile("/tmp"), true)
  assert.Equal(t, writeCacheFile("/foo"), false)
}
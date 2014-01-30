package main

import(
  "fmt"
  "os"
  "io"
  "strings"
  "crypto/sha1"
  "github.com/kr/s3/s3util"
  "github.com/jessevdk/go-flags"
)

const VERSION = "0.1.0"

var options struct {
  Path     string `long:"path"      description:"Path to package"`
  S3Key    string `long:"s3-key"    description:"Amazon S3 access key"`
  S3Secret string `long:"s3-secret" description:"Amazon S3 secret key"`
  S3Bucket string `long:"s3-bucket" description:"Amazon S3 bucket name"`
}

func envDefined(name string) bool {
  result := os.Getenv(name)
  return len(result) > 0
}

func fileExists(path string) bool {
  _, err := os.Stat(path)
  return err == nil
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

func printUsage() {
  fmt.Printf("Bunch v%s\n", VERSION)
  fmt.Println("Usage: bunch [rubygems|npm] [upload|download]")
}

func loadS3Credentials() {
  if envDefined("S3_KEY") {
    options.S3Key = os.Getenv("S3_KEY")
  }

  if envDefined("S3_SECRET") {
    options.S3Secret = os.Getenv("S3_SECRET")
  }

  if envDefined("S3_BUCKET") {
    options.S3Bucket = os.Getenv("S3_BUCKET")
  }
}

func checkS3Credentials() {
  if len(options.S3Key) == 0 { 
    fatal("S3 access key is not set.")
  }
  
  if len(options.S3Secret) == 0 { 
    fatal("S3 secret key is not set.")
  }
  
  if len(options.S3Bucket) == 0 { 
    fatal("S3 bucket name is not set.")
  }
}

func isUrl(s string) bool {
  return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func s3Url(filename string) string {
  format := "https://s3.amazonaws.com/%s/%s"
  url := fmt.Sprintf(format, options.S3Bucket, filename)

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

func transferArchive(file string, url string, fail_status int) {
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

func handleRubygems(action string) {
  /* NOOP */
}

func handleNpm(action string) {
  /* NOOP */
}

func main() {
  args := os.Args

  if (len(args) < 3) {
    printUsage()
    os.Exit(1)
  }

  new_args, err := flags.ParseArgs(&options, os.Args)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  fmt.Println(new_args)

  service := args[1]
  action  := args[2]

  loadS3Credentials()
  checkS3Credentials()

  if len(options.Path) == 0 {
    options.Path, _ = os.Getwd()
  }

  if len(options.Prefix) == 0 {
    options.Prefix = filepath.Base(options.Path)
  }

  fmt.Println(options)

  switch service {
  case "rubygems":
    handleRubygems(action)
    return
  case "npm":
    handleNpm(action)
    return
  default:
    fatal("Invalid service")
  }
}
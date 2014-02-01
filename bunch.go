package main

import(
  "fmt"
  "os"
  "os/exec"
  "github.com/jessevdk/go-flags"
  "github.com/kr/s3/s3util"
  "path/filepath"
  "io/ioutil"
  "runtime"
)

const VERSION = "0.1.0"

var options struct {
  Prefix   string `long:"prefix"    description:"Archive prefix"`
  Path     string `long:"path"      description:"Path to cache"`
  Manifest string `long:"manifest"  description:"Path to Gemfile.lock or package.json"`
  S3Key    string `long:"s3-key"    description:"Amazon S3 access key"`
  S3Secret string `long:"s3-secret" description:"Amazon S3 secret key"`
  S3Bucket string `long:"s3-bucket" description:"Amazon S3 bucket name"`
}

func printUsage() {
  fmt.Printf("Bunch v%s\n", VERSION)
  fmt.Println("Usage: bunch [upload|download]")
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

func setS3Credentials() {
  s3util.DefaultConfig.AccessKey = options.S3Key
  s3util.DefaultConfig.SecretKey = options.S3Secret
}

func checkOptions() {
  if options.Prefix == ""   { fatal("Please specify --prefix")    }
  if options.Manifest == "" { fatal("Please specify --manifest")  }
  if options.Path == ""     { fatal("Please specify --path")      }
  if options.S3Key == ""    { fatal("Please specify --s3-key")    }
  if options.S3Secret == "" { fatal("Please specify --s3-secret") }
  if options.S3Bucket == "" { fatal("Please specify --s3-bucket") }
}

// Expand path arguments to support tildas. Example: "~/path" 
//
func expandPathArguments() {
  options.Path = expandPath(options.Path)
  options.Manifest = expandPath(options.Manifest)
}

// Extract an archive to a path. Returs true on success
//
// extract("/tmp/archive.tar.gz", "/path/to/extract")
//
func extract(filename string, path string) bool {
  // Create a directory to extract archive to
  if _, err := exec.Command("mkdir", path).Output() ; err != nil {
    fmt.Println("Unable to create extaction path")
    return false
  }

  // Extract tarball to created directory
  if _, err := exec.Command("tar", "-xzf", filename, "-C", path).Output() ; err != nil {
    fmt.Println("Unable to extract archive")
    return false
  }

  return true
}

// Create a new gzip-compressed tarball
//
// archive("/path/to/archive", "/tmp/archive.tar.gz")
//
func archive(path string, archive_path string) bool {
  cmd := fmt.Sprintf("cd %s && tar -czf %s .", path, archive_path)

  if _, err := exec.Command("bash", "-c", cmd).Output() ; err != nil {
    fmt.Println("Unable to create archive")
    return false
  }

  return true
}

// Creaate .bunch file to indicate that dependencies were cached
//
// writeCacheFile("/path/to/bundle")
//
func writeCacheFile(path string) bool {
  cache_file := fmt.Sprintf("%s/.bunch", path)

  if _, err := exec.Command("touch", cache_file).Output(); err != nil {
    fmt.Println("Unable to create .bunch cache file")
    return false
  }

  return true
}

// Download an archive from Amazon S3 bucket and extract it to a directory
//
// download("http://example.com/archive.tar.gz", "/path/to/extract")
//
func download(url string, extract_path string) {
  // Save to temporary file named after archive
  save_path := fmt.Sprintf("/tmp/%s", filepath.Base(url))

  fmt.Println("Downloading archive from Amazon S3...")
  transfer(url, save_path, 0)

  fmt.Println("Extracting archive...")
  if !extract(save_path, extract_path) {
    os.Remove(save_path)
    os.Exit(1)
  }

  // Mark extracted archive as cached
  writeCacheFile(extract_path)

  fmt.Println("Done")
  os.Exit(0)
}

func upload(name string, path string, manifest_path string) {
  manifest, err := ioutil.ReadFile(manifest_path)

  if err != nil {
    fatal("Unable to read manifest file")
  }

  if string(manifest) == "" {
    fatal("Manifest is empty")
  }

  // Generate a SHA1 hash for manifest file
  checksum := sha1Checksum(string(manifest))

  // Create a new archive
  archive_path := fmt.Sprintf("/tmp/%s_%s_%s.tar.gz", name, checksum, runtime.GOARCH)

  fmt.Println("Archiving bundle...")

  if !archive(path, archive_path) {
    fatal("Unable to create archive")
  }

  // Upload archive to S3
  fmt.Println("Uploading bundle to Amazon S3...")
  transfer(archive_path, s3Url(options.S3Bucket, filepath.Base(archive_path)), 0)

  fmt.Println("Done")
  os.Exit(0)
}

func handleUpload() {
  if fileExists(fmt.Sprintf("%s/.bunch", options.Path)) {
    fatal("Already cached")
  }

  if !fileExists(options.Path) {
    fatal("Path does not exist")
  }

  if !fileExists(options.Manifest) {
    fatal("Manifest file does not exist")
  }

  upload(options.Prefix, options.Path, options.Manifest)
}

func handleDownload() {
  // Check if extraction directory already exists
  if fileExists(options.Path) {
    fatal("Path already exists")
  }

  // Load manifest file
  manifest, err := ioutil.ReadFile(options.Manifest)

  if err != nil {
    fatal("Unable to read manifest file")
  }

  if string(manifest) == "" {
    fatal("Manifest is empty")
  }

  // Generate SHA1 hash from manifest file contents
  checksum := sha1Checksum(string(manifest))

  // Build download url
  filename := fmt.Sprintf("%s_%s_%s.tar.gz", options.Prefix, checksum, runtime.GOARCH)
  url := s3Url(options.S3Bucket, filename)

  // Download and extract archive
  download(url, options.Path)
}

func handleCommand(command string) {
  switch command {
  default:
    fmt.Println("Invalid command:", command)
    printUsage()
  case "upload":
    handleUpload()
  case "download":
    handleDownload()
  }
}

func getCommand() string {
  if len(os.Args) < 2 {
    printUsage();
    os.Exit(1)
  }

  args, err := flags.ParseArgs(&options, os.Args)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  return args[1]
}

func main() {
  command := getCommand()

  loadS3Credentials()
  checkOptions()
  setS3Credentials()
  expandPathArguments()
  handleCommand(command)
}
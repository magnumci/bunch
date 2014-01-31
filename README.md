# bunch

Utility to cache project dependencies (Ruby, NPM) to Amazon S3 bucket.

## Overview

Bunch is a tool to upload/download contents of any directory from Amazon S3. Initially
created to cache Rubygems/NPM bundles as compressed tarballs (tar.gz) to speed up 
CI build times.

## Usage

Run bunch upload/download command with the following options:

``` bash
bunch [upload|download] \
  --prefix=my-project \
  --path=path/to/dir \
  --manifest=path/to/manifest
  --s3-key=key \
  --s3-secret=secret \
  --s3-bucket=bucket
```

If you don't want to specify Amazon S3 credentials in terminal, you can always
export variables into your environment and they will be loaded automatically.

Example:

``` bash
export S3_KEY=key
export S3_SECRET=secret
export S3_BUCKET=bucket
```

## Build

Project requires Go 1.2.

To build a binary run, install dependencies and execute build command:

```
go get
go build
```

To build for multiple operating systems and architectures, use [gox](https://github.com/mitchellh/gox).

## License

The MIT License (MIT)

Copyright (c) 2014 Dan Sosedoff, <dan.sosedoff@gmail.com>
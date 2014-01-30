# bunch

Utility to cache project dependencies (Ruby, NPM) to Amazon S3 bucket.

## Usage

Available options:

```
Usage:
  bunch [OPTIONS]

Application Options:
      --prefix=    Archive prefix
      --path=      Path to cache
      --manifest=  Path to Gemfile.lock or package.json
      --s3-key=    Amazon S3 access key
      --s3-secret= Amazon S3 secret key
      --s3-bucket= Amazon S3 bucket name

Help Options:
  -h, --help       Show this help message
```

If you dont want provide Amazon S3 credentials in terminal, you can always
export variables into your environment and they'll going to be automatically
loaded.

Example:

``` bash
export S3_KEY=key
export S3_SECRET=secret
export S3_BUCKET=bucket
```

### Examples

To upload bundle, execute command as follows:

``` bash
bunch upload \
  --path ~/my-project/.bundle \
  --manifest ~/my-project/Gemfile.lock \
  --prefix my-project
```

To download bundle cache:

``` bash
bunch download \
  --path ~/my-project/.bundle \
  --manifest ~/my-project/Gemfile.lock \
  --prefix my-project
```
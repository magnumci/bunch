# bunch

Utility to cache project dependencies (Ruby, NPM) to Amazon S3 bucket.

## Usage

```
bunch rubygems [upload|download]
bunch npm [upload|download]
```

Required environment variables:

- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`
- `S3_BUCKET`
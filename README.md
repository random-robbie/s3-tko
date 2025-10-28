# s3-tko

AWS S3 Bucket Takeover Scanner - A concurrent URL scanner that detects potentially vulnerable S3 buckets by checking for "NoSuchBucket" errors.

### Shout out to TomNomNom for 99.9% of his code....

## Requirements

- Go 1.24.4 or later

## Build

```bash
go mod tidy
go build -o s3-tko
```

## Usage

From stdin:
```bash
cat urls.txt | ./s3-tko
```

From file:
```bash
./s3-tko urls.txt
```

## Output

- Vulnerable URLs are logged to `text.log`
- Scan results are printed to stdout
- Errors are printed to stderr

## Features

- Concurrent scanning with 12 workers
- DNS resolution checking
- TLS support with InsecureSkipVerify
- 5-second request timeout
- Thread-safe logging
- Proper resource cleanup

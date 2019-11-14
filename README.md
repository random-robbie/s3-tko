# s3-tko
AWS S3 Bucket Finder it will take a list of urls and check for the nobucket text.

### Shout out to TomNomNom for 99.9% of his code....

### Build

```
go get github.com/fatih/color
go build
```

### Usage

```
cat list.txt | ./s3-tko
```

All vuln urls are logged in text.log and then its simple up to you to do a TKO.

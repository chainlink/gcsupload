
# GCS File Uploader

This utility takes a json google credential and a file argument and uploads it to Google Cloud Storage. Allow for a simple small binary to be used in CI and other environments without installing deps required for `gsutil` (Python, etc)

```
$ gcsupload -h
Usage of gcsupload:
  -b string
        Bucket name
  -f string
        File to upload, will be named the same in the bucket
  -k string
        GCP key.json file
```

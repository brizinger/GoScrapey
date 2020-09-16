# GoScrapey ![CodeScan](https://github.com/brizinger/GoScrapey/workflows/Go/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/brizinger/GoScrapey)](https://goreportcard.com/report/github.com/brizinger/GoScrapey) [![codebeat badge](https://codebeat.co/badges/49f2e42d-e78a-4fee-939e-ecf13feb2b7b)](https://codebeat.co/projects/github-com-brizinger-goscrapey-master)
Go tool that scrapes images from websites and downloads them.

# Install
```go build .```

Or use the Makefile provided:

```make build``` will build the tool and place the binary in the folder bin as scrapey.

# Usage
```scrapey https://www.google.com/``` will download the images from google.com and place them in a default location (home directory).

You can also use the -d (--directory) flag to place the images in another location.

-u will upload the images on Imgur 

Note: To use the tool from any directory, you need to export the current path or place the binary to a path you have already exported.

# Contribution
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)

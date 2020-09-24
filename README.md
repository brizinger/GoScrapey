# GoScrapey ![CodeScan](https://github.com/brizinger/GoScrapey/workflows/Go/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/brizinger/GoScrapey)](https://goreportcard.com/report/github.com/brizinger/GoScrapey) [![codebeat badge](https://codebeat.co/badges/49f2e42d-e78a-4fee-939e-ecf13feb2b7b)](https://codebeat.co/projects/github-com-brizinger-goscrapey-master)

Go tool that scrapes images from websites and downloads them.

# Install

`go build .`

Or use the Makefile provided:

`make build` will build the tool and place the binary in the folder bin as scrapey.

# Usage

`scrapey https://www.google.com/` will download the images from google.com and place them in a default location (home directory).

You can also use the -d (--directory) flag to place the images in another location. The -u flag will upload all images to an imgur album.

Note: The tool needs a Client-ID, which it reads from a file (ID.txt) placed in the root of the project. You need to supply your own Client-ID, which could be retrieved from the official Imgur api if you are going to make changes or do not want to use the already build file.

# Contribution

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)

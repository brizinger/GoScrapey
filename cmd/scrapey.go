package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var rootCmd = &cobra.Command{
	Use:   "scrapey",
	Short: "GoScrapey is a website image scraper",
	Long: `An image scraper build with love by Brizinger in Go
	More information can be found on http://github.com/brizinger/GoScrapey`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a web URL to read from")
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		scrapeWeb(args[0])
	},
}

var directory string // directory flag var
var upload bool      // upload flag var
var hashes []string  // deletehashes of all images
var originalDir string

func init() {
	path, err := homedir.Dir()
	checkError(err)
	info := "Directory where the files will be written. Default is " + path
	rootCmd.Flags().StringVarP(&directory, "directory", "d", "", info)
	rootCmd.Flags().BoolVarP(&upload, "upload", "u", false, "Upload to Imgur")

}

// Execute - executes the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Creates a file for the image and then places the image there
func createImg(webURL string, imgURL string, client *http.Client) (f *os.File, err error) {

	file, err := os.Create(createFileName(imgURL)) // Create file

	if getImageExtension(file.Name()) == "svg" {
		file.Close()
		os.Remove(file.Name())
		errText := "File " + file.Name() + " is of type SVG, which is currently not supported.\n"
		err = errors.New(errText)
		return nil, err
	}

	checkError(err)
	log.Println("Createing file for image")
	webURL = getHostName(webURL)
	fullURL := "http://www." + webURL + imgURL
	log.Println(fullURL)
	resp, err := client.Get(fullURL) // Open image address
	checkError(err)
	log.Println("Opening image ...")
	defer resp.Body.Close()

	io.Copy(file, resp.Body) // Copy img to file
	log.Printf("Writing image to file %s...", createFileName(imgURL))
	defer file.Close()

	fsize, err := file.Stat()
	checkError(err)
	size := fsize.Size() / 1024
	log.Printf("Image size in KB: %v\n", size)

	checkError(err)

	if upload {
		hashes = append(hashes, uploadImage(file))
	}
	f = file
	return f, nil

}

func createAlbum(imageHashes []string, web string) {
	website := "Images from: " + web

	url := "https://api.imgur.com/3/album"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for i := 0; i < len(imageHashes); i++ {
		_ = writer.WriteField("deletehashes[]", imageHashes[i])
	}
	_ = writer.WriteField("title", website)
	_ = writer.WriteField("description", "Created with love by GoScrapey")
	_ = writer.WriteField("privacy", "hidden")
	err := writer.Close()
	checkError(err)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	checkError(err)
	auth := GETAPIID()
	req.Header.Add("Authorization", auth)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	checkError(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkError(err)
	result := string(body)
	results := strings.Split(result, ",")

	if strings.Contains((string(body)), "error") {
		fmt.Println(string(body))
		os.Exit(1)
	}

	for i := 0; i < len(results); i++ {
		if strings.Contains(results[i], "id") {
			id := results[i]
			id = strings.Replace(id, `"`, "", -1)          // remove all "
			id = strings.Replace(id, `}`, "", -1)          // remove } at the end
			id = strings.Replace(id, `\`, "", -1)          // remove \
			id = strings.Replace(id, `{data:{id:`, "", -1) // remove unnecessary stuff
			link := `https://imgur.com/a/` + id
			log.Println("*** Link to image album: " + link + " ***")
		}
	}
}

func uploadImage(image *os.File) string { // Upload single image to imgur
	fstat, err := image.Stat()
	checkError(err)
	input, err := os.Open(fstat.Name())
	checkError(err)
	reader := bufio.NewReader(input)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encodedImg := base64.StdEncoding.EncodeToString(content)

	url := "https://api.imgur.com/3/image"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", encodedImg)
	err = writer.Close()
	checkError(err)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	checkError(err)
	auth := GETAPIID() // Gets the ID from api-key.go
	req.Header.Add("Authorization", auth)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	checkError(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	result := string(body)
	results := strings.Split(result, ",")

	if strings.Contains((string(body)), "error") {
		fmt.Println(string(body))
		os.Exit(1)
	}

	var hash string

	for i := 0; i < len(results); i++ {
		if strings.Contains(results[i], "link") {
			link := results[i]
			link = strings.Replace(link, `"`, "", -1) // remove all "
			link = strings.Replace(link, `}`, "", -1) // remove } at the end
			link = strings.Replace(link, `\`, "", -1) // remove \
			// log.Println("!!! " + link + " !!!")
		}

		if strings.Contains(results[i], "deletehash") {
			deletehash := results[i]
			deletehash = strings.Replace(deletehash, `"`, "", -1)           // remove all "
			deletehash = strings.Replace(deletehash, `}`, "", -1)           // remove } at the end
			deletehash = strings.Replace(deletehash, `\`, "", -1)           // remove \
			deletehash = strings.Replace(deletehash, `deletehash:`, "", -1) // remove \
			// log.Println("!!! " + deletehash + " !!!")
			hash = deletehash
		}
	}
	getDirectory()
	return hash
}

func getOriginalDir() string {
	dir, err := os.Getwd()
	checkError(err)
	return dir
}

// Returns the file name of the image without the relative image path in the website
func createFileName(imgURL string) string {
	slice := strings.Split(imgURL, "/")

	return slice[len(slice)-1]
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Checks if directory is specified
func getDirectory() {

	if directory != "" {
		os.Chdir(directory)
	} else {
		path, err := homedir.Dir()
		checkError(err)
		os.Chdir(path)
		log.Println("No directory specified, defaulting to home")
	}
}

func getImageExtension(file string) string {
	slice := strings.Split(file, ".")

	return slice[len(slice)-1]
}

func getHostName(URL string) string {
	u, err := url.Parse(URL)
	checkError(err)
	return u.Hostname()
}

// Opens the specified page and downloads the images
func scrapeWeb(web string) {

	var files []*os.File

	originalDir = getOriginalDir()
	getDirectory()

	url := web

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		log.Println("Load page complete")

		if resp != nil {
			log.Println("Page response is NOT nil")
			// --------------
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			hdata := strings.Replace(string(data), "<noscript>", "", -1)
			hdata = strings.Replace(hdata, "</noscript>", "", -1)
			// --------------

			if document, err := html.Parse(strings.NewReader(hdata)); err == nil {
				var parser func(*html.Node)
				parser = func(n *html.Node) {
					if n.Type == html.ElementNode && n.Data == "img" {

						var imgSrcURL, imgDataOriginal string

						for _, element := range n.Attr {
							if element.Key == "src" {
								imgSrcURL = element.Val
								file, err := createImg(web, imgSrcURL, httpClient())
								if err != nil {
									log.Printf("ERROR: %s \n", err.Error())
									continue
								}
								files = append(files, file)
								log.Printf("File downloaded: %v\n", file.Name())
							}
							if element.Key == "data-original" {
								imgDataOriginal = element.Val
							}
						}

						log.Println("Found from: ", imgSrcURL, imgDataOriginal)

					}

					for c := n.FirstChild; c != nil; c = c.NextSibling {
						parser(c)
					}

				}
				parser(document)
			} else {
				log.Panicln("Parse html error", err)
			}

		} else {
			log.Println("Page response IS nil")
		}
	}
	if upload {
		createAlbum(hashes, web)
	}
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

var directory string

func init() {
	rootCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory where the files will be written. Default is home")
}

// Execute - executes the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createImg(webURL string, imgURL string, client *http.Client) {

	file, err := os.Create(createFileName(imgURL)) // Create file
	checkError(err)
	log.Println("Createing file for image")
	fullURL := webURL + imgURL

	resp, err := client.Get(fullURL) // Open image address
	checkError(err)
	log.Println("Opening image ...")
	defer resp.Body.Close()

	io.Copy(file, resp.Body) // Copy img to file
	log.Printf("Writing image to file %s...", createFileName(imgURL))
	defer file.Close()

	checkError(err)

}

func createFileName(imgURL string) string {
	slice := strings.Split(imgURL, "/")

	return slice[len(slice)-1]
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

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

func scrapeWeb(web string) {

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
								createImg(web, imgSrcURL, httpClient())
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

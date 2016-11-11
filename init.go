/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"io/ioutil"
	"net/http"
)

type Configuration struct { // The Configuration type holds configuration data
	Host        string // Host string for the webserver to listen on
	Port        string // Port string for the webserver to listen on
	PublicDir   string // Path string to the directory to serve static files from
	CacheStatic bool   // Boolean to enable or disable file caching
	DisableCORS bool   // Boolean to strip CORS headers
	ExternalURL string // External URL string for formatting proxied HTML
	EnableTLS   bool   // Boolean to serve with TLS
}

type reqHandler func(http.ResponseWriter, *http.Request) *reqError

var Config Configuration        // Configuration for the entire program
var FileCache map[string][]byte // Files cached in the memory, stored as byte slices in a map that takes strings for the file names

func init() { // Init function
	folderPath, err := osext.ExecutableFolder() // Figure out where we are in the filesystem to make specifying the location of the public directory easier
	if err != nil {
		folderPath = ""          // If this doesn't work it's not a huge deal and we can just set the folder path to an empty string and print an error message
		fmt.Println(err.Error()) // Print an error message but don't do anything else
	}
	// Configuration flags
	flag.StringVar(&Config.Host, "host", "localhost", "host to listen on for the webserver")
	flag.StringVar(&Config.Port, "port", "8000", "port to listen on for the webserver")
	flag.StringVar(&Config.PublicDir, "pubdir", folderPath+"/pub", "path to the static files the webserver should serve")
	flag.BoolVar(&Config.CacheStatic, "cachestatic", true, "cache specific heavily used static files")
	flag.BoolVar(&Config.DisableCORS, "cors", true, "strip Cross Origin Resource Policy headers")
	flag.StringVar(&Config.ExternalURL, "exturl", Config.Host+":"+Config.Port, "external URL for formatting proxied HTML files to link back to the webproxy")
	flag.BoolVar(&Config.EnableTLS, "tls", false, "enable serving with TLS (https), certificate is cert.pem and key is key.pem, place both in the directory your terminal instance is in")
}

func main() { // Main function

	var err error

	FileCache = make(map[string][]byte) // Make the map for caching files
	if Config.CacheStatic == true {     // Cache certain static files if they exist and if Config.CacheStatic is set to true
		FileCache["index"], err = ioutil.ReadFile(Config.PublicDir + "/index.html")
		if err != nil {
			FileCache["index"] = nil
		}
		FileCache["404"], err = ioutil.ReadFile(Config.PublicDir + "/404.html")
		if err != nil {
			FileCache["404"] = nil
		}
	}
	// Create a HTTP Server, and handle requests and errors
	http.Handle("/", reqHandler(static))
	http.Handle("/p/", reqHandler(proxy))
	bind := fmt.Sprintf("%s:%s", Config.Host, Config.Port)
	fmt.Printf("Bypass listening on %s...\n", bind)
	if !Config.EnableTLS {
		err = http.ListenAndServe(bind, nil)
		if err != nil {
			panic(err)
		}
	} else if Config.EnableTLS {
		err = http.ListenAndServeTLS(bind, "cert.pem", "key.pem", nil)
		if err != nil {
			panic(err)
		}
	}
}

func (fn reqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { // Allows us to pass errors back through our http handling functions
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		fmt.Println(e.Error.Error(), e.Message) // Print the error message
		if e.Code == 404 {                      // Serve a pretty (potentially cached) file for 404 errors, if it exists
			w.WriteHeader(404)
			if FileCache["404"] != nil { // Serve the cached file if one exists
				io.WriteString(w, string(FileCache["404"]))
			} else { // Read a non-cached file from disk and serve it because there isn't a cached one
				file, err := ioutil.ReadFile(Config.PublicDir + "/404.html")
				if err != nil {
					http.Error(w, e.Message+"\n"+e.Error.Error(), e.Code) // Serve a generic error message if the file isn't cahced and doesn't exist
					return
				}
				io.WriteString(w, string(file))
			}
		} else { // If it's not a 404 error just serve a generic message
			http.Error(w, e.Message+"\n"+e.Error.Error(), e.Code)
		}

	}
}

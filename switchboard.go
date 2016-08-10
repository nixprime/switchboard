// Switchboard: link redirect server
// Copyright (C) 2013 the Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	DefaultConfigFilename string = "/etc/switchboard.conf"
	DefaultPort           uint16 = 80
)

type Switchboard map[string]interface{}

func (s Switchboard) handlePage(target string, w http.ResponseWriter,
	r *http.Request, path string) {
	var url string
	if strings.Contains(target, "%s") {
		url = fmt.Sprintf(target, path)
	} else {
		url = target
	}
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusFound)
}

func (s Switchboard) handlePath(pathMap map[string]interface{},
	w http.ResponseWriter, r *http.Request, path string) {
	pathParts := strings.SplitN(path, "/", 2)
	var pathHead, pathTail string
	if len(pathParts) >= 1 {
		pathHead = pathParts[0]
	}
	if len(pathParts) >= 2 {
		pathTail = pathParts[1]
	}
	i, ok := pathMap[pathHead]
	if ok {
		s.handleBranch(i, w, r, pathTail)
	} else {
		i, ok = pathMap["*"]
		if ok {
			s.handleBranch(i, w, r, path)
		} else {
			s.handleDefault(w)
		}
	}
}

func (s Switchboard) handleBranch(i interface{}, w http.ResponseWriter,
	r *http.Request, path string) {
	switch v := i.(type) {
	case string:
		s.handlePage(v, w, r, path)
	case map[string]interface{}:
		s.handlePath(v, w, r, path)
	default:
		log.Fatalf("Unexpected type in path map: %s",
			reflect.TypeOf(i).String())
	}
}

func (s Switchboard) HandleHost(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s http://%s%s", r.Method, r.Host, r.URL.Path)
	i, ok := s[r.Host]
	if ok {
		path := r.URL.Path
		var p string
		if len(path) > 0 && path[0] == '/' {
			p = path[1:]
		} else {
			p = path
		}
		s.handleBranch(i, w, r, p)
	} else {
		s.handleDefault(w)
	}
}

func (s Switchboard) handleDefault(w http.ResponseWriter) {
	hostMapString, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalf("Unexpected error marshaling host map: %s", err.Error())
	}
	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml/" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="Switchboard" content="Switchboard summary" />
  <title>Switchboard</title>
</head>
<body>
  <h1>Welcome to Switchboard</h1>
  <p>Here's what we've got:</p>
  <pre>%s</pre>
</body>
</html>`, hostMapString)))
}

func main() {
	// Read the configuration file
	var configFilename string
	flag.StringVar(&configFilename, "f", DefaultConfigFilename,
		"configuration file")
	flag.Parse()
	configFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("Error reading configuration file %s: %s", configFilename,
			err.Error())
	}
	var configBlob interface{}
	err = json.Unmarshal(configFile, &configBlob)
	if err != nil {
		log.Fatalf("Error parsing configuration file %s: %s", configFilename,
			err.Error())
	}
	config := configBlob.(map[string]interface{})

	// Select the TCP address (interface and port) to listen on
	var listenOn string
	var port uint16

	listenOnConfigBlob, ok := config["listen_on"]
	if ok {
		listenOn = listenOnConfigBlob.(string)
	} else {
		listenOn = ""
	}
	portConfigBlob, ok := config["port"]
	if ok {
		port = uint16(portConfigBlob.(float64))
	} else {
		port = DefaultPort
	}
	listenAddr := listenOn + ":" + strconv.FormatUint(uint64(port), 10)

	// Read the host map
	hostBlob, ok := config["hosts"]
	if !ok {
		log.Fatalf("Error parsing configuration file %s: Missing hosts",
			configFilename)
	}
	hostMap := Switchboard(hostBlob.(map[string]interface{}))

	// Start the server
	http.HandleFunc("/", hostMap.HandleHost)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err.Error())
	}
}

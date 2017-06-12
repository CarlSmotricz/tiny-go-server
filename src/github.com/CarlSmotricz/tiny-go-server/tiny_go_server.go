// Copyright 2017, Carl Smotricz.
// Rights: What rights?
// License: MIT.

// Based on the "Writing Web Applicatins" tutorial page
// by the Go authors: https://golang.org/doc/articles/wiki/ .

package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	_ "regexp"
	_ "sort"
	"time"
)

type WebOp struct {
	Title string
	Files []os.FileInfo
}

const (
	NAME       = "tiny-go-server"
	PORT       = 8080
	MAIN_DIR   = "/home/ubuntu/workspace/server"
	SUBDIR_LIB = "lib"
	SUBDIR_UPL = "upl"
	TEMPLATES  = `
		{{define "header"}}
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="UTF-8">
					<title>{{.}} - BotLibrarian</title>
					<style type="text/css">
						body { margin: 1cm; }
						h1,th { background-color: #CCC; }
					</style>
				</head>
				<body>
					<h1>BotLibrarian</h1>
		{{end}}
		
		{{define "footer"}}
			{{with .}}
				<br/>
				<a href="/">Home</a>
			{{end}}
			</body>
		</html>
		{{end}}
		
		{{define "FileTab"}}
			<h2>{{.Title}}</h2>
			<table>
			<tr><th>Timestamp</th><th>Size</th><th>File name</th></tr>
			<tr><td colspan="3"><hr/></td/></tr>
			{{range .Files}}
				<tr>
					<td>{{TS .ModTime}}</td>
					<td align="right">{{.Size}}</td>
					<td><a href="/raw-lib/{{.Name}}">{{.Name}}</a></td>
				</tr>
			{{end}}
			<tr><td colspan="3"><hr/></td/></tr>
			</table>
		{{end}}
		
		{{define "T_Lib"}}
			{{template "header" "Library"}}
			{{template "FileTab" .}}
			{{template "footer" true}}
		{{end}}
		
		{{define "T_Ups"}}
			{{template "header" "Uploads"}}
			{{template "FileTab" .}}
			{{template "footer" true}}
		{{end}}
		
		{{define "T_Upl"}}
			{{template "header" "Upload"}}
			<h2>Upload a file manually</h2>
			<form action="/upload/file" method="post" enctype="multipart/form-data">
				<p>Select a local file for upload...</p>
				<input type="file" name="file" id="fileToUpload"><br/>
				<p>...then hit "Upload" to upload it!</p>
				<input type="submit" value="Upload">
			</form>
			{{template "footer" true}}
		{{end}}
		
		{{define "T_Wel"}}
			{{template "header" "Welcome"}}
			<h2>Welcome</h2>
			<p>Available actions:</p>
			<table>
				<tr><td><a href="/library">Library</td><td> - </td><td>View the library</td></tr>
				<tr><td><a href="/uploaded">Uploads</td><td> - </td><td>View unprocessed uploads</td></tr>
				<tr><td><a href="/upload">Upload</td><td> - </td><td>Upload a file manually</td></tr>
			</table>
			{{template "footer"}}
		{{end}}`
)

var (
	templates = template.Must(template.New("XXX").Funcs(template.FuncMap{"TS": MyTimestamp}).Parse(TEMPLATES))
	EMPTY     = []os.FileInfo{}
)

func main() {
	http.HandleFunc("/library/", handleLib)
	http.HandleFunc("/raw-lib/", handleRawLib)
	http.HandleFunc("/uploaded/", handleUploaded)
	http.HandleFunc("/raw-upl/", handleRawUpl)
	http.HandleFunc("/upload/file", handleUploadFile)
	http.HandleFunc("/upload/", handleUpload)
	http.HandleFunc("/log", handleLog)
	http.HandleFunc("/", handleDefault)
	ownAddr := fmt.Sprintf(":%d", PORT)
	log.Printf("%s running.", NAME)
	http.ListenAndServe(ownAddr, nil)
}

// Display list of files in the library.
func handleLib(w http.ResponseWriter, r *http.Request) {
	dirPath := filepath.Join(MAIN_DIR, SUBDIR_LIB)
	files, err := listFiles(dirPath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
	} else {
		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "T_Lib", WebOp{"Library", files})
		if err != nil {
			log.Println(err)
		}
	}
}

// Retrieve a file as plain text from the library.
func handleRawLib(w http.ResponseWriter, r *http.Request) {
	retrieveRaw(w, r, "/raw-lib/", SUBDIR_LIB, r.URL)
}

// Display list of files in "uploads".
func handleUploaded(w http.ResponseWriter, r *http.Request) {
	dirPath := filepath.Join(MAIN_DIR, SUBDIR_UPL)
	files, err := listFiles(dirPath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
	} else {
		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "T_Ups", WebOp{"Uploads", files})
		if err != nil {
			log.Println(err)
		}
	}
}

// Retrieve a file as plain text from the uploads.
func handleRawUpl(w http.ResponseWriter, r *http.Request) {
	retrieveRaw(w, r, "/raw-upl/", SUBDIR_UPL, r.URL)
}

// Display the manual file upload form.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "T_Upl", WebOp{"Upload a file", EMPTY})
	if err != nil {
		log.Println(err)
	}
}

// Process a file upload.
func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(filepath.Join(SUBDIR_UPL, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		fmt.Fprintln(w, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	log.Printf("Uploaded file: '%s'.", handler.Filename)
	handleUploaded(w, r)
}

// Process a "log" request: writes the request string to log.
func handleLog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	values := r.Form
	if values["msg"] == nil {
		log.Printf("No 'msg' in log request")
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request.")
	} else {
		log.Printf("#BOTLOG# %s\n", values["msg"][0])
		fmt.Fprintln(w, "\nOK.")
	}
}

// Handle any URLs not explicitly handled.
func handleDefault(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "T_Wel", WebOp{"Welcome", EMPTY})
	if err != nil {
		log.Println(err)
	}
}

// Custom directory listing.
func listFiles(path string) (files []os.FileInfo, err error) {
	dirFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dirFile.Close()
	fileInfos, err := dirFile.Readdir(0)
	// TODO: Sort file list.
	// sort.Slice(fileInfos, func(i, j int) bool { return fileInfos[i].Name() < fileInfos[j].Name() })
	return fileInfos, nil
}

// Retrieve a file as raw text (= download)
func retrieveRaw(w http.ResponseWriter, r *http.Request, expectedDir string, physicalDir string, url *url.URL) {
	log.Printf("Request: %s\n", url.Path)
	content, err := retrieveFile(expectedDir, physicalDir, url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		fmt.Fprintln(w, "Bad request.")
	} else {
		defer content.Close()
		http.ServeContent(w, r, "content", time.Now(), content)
	}
}

// Retrieve the text of a file for raw download.
func retrieveFile(expectedDir, physicalDir string, url *url.URL) (content *os.File, err error) {
	reqDir, reqFileName := filepath.Split(url.Path)
	if reqDir != expectedDir {
		log.Printf("Bad directory name: '%s', should be '%s'", reqDir, expectedDir)
		return nil, errors.New("Bad path.")
	}
	path := filepath.Join(physicalDir, reqFileName)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return f, nil
}

// Timestamp "to the minute" for directory listings
func MyTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

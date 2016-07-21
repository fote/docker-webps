package main

import (
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"bufio"
	"fmt"
)

var logger *log.Logger
var header string = `<html>
<head>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
<title>Docker ps</title>
<body>
<nav class="navbar navbar-default navbar-static-top">
      <div class="container">
        <div class="navbaras-header">
          <a class="navbar-brand" href="#">Docker webps</a>
        </div>
        </div><!--/.nav-collapse -->
      </div>
    </nav>
    <div class="container">
      <div class="starter-template">
        `
var footer string = `</div></div></body></html>`


func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}



func main() {
	logger.Println("Starting http server on :7777")

	r := httprouter.New()
	r.GET("/", HomeHandler)

	if err := http.ListenAndServe(":7777", r); err != nil {
		logger.Println(err)
		os.Exit(1)
	}

}

func ReturnErrorInBrowser(rw http.ResponseWriter) {
		if rec := recover(); rec != nil {
			logger.Println(rec)
			http.Error(rw, "Oops..something went wrong", http.StatusInternalServerError)
		}
	}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer ReturnErrorInBrowser(rw)
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, nil)
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		logger.Println("Error while loading list of containers")
		logger.Println("Can't connect to unix-sock found")
		logger.Println("Maybe i'm inside container?")
		logger.Println("Try to connect docker-api to 4243")
		f, err := os.Open("/proc/net/route")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			fmt.Println(s.Text())
		}

		panic(err)
	}
	body := `<table class="table table-striped">
	<tr> <th>ID</th> <th>Name</th> <th>Image</th> <th>Command</th> <th>State</th> </tr>`
	for _, c := range containers {
		body += `<tr><th>` + c.ID[:7] + `</th><th>` + c.Names[0][1:] + `</th><th>` + c.Image + `</th><th>`
		body += c.Command + `</th><th>`+ c.State + `</tr>`
	}
	body += `</table>`
	rw.Write([]byte(header + body + footer))
}

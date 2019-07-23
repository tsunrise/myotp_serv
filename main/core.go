package main

import (
	"flag"
	"log"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"net/http"
)

func main() {
	// parse flags
	installMode := flag.Bool("install", false, "install mode")
	flag.Parse()

	if *installMode {
		mydb.Install()
		return
	}

	log.Println("Started: MyOTP Backend Server Development Edition")
	http.Handle("/", httpServer{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type httpServer struct {}


func (s httpServer) ServeHTTP(response http.ResponseWriter,request *http.Request) {
	// before the url matching
	//...

	switch request.URL.Path{
	case "/":
		logRequest(request)
		hello := struct {
			Hello string
		}{"world"}
		resp := shell.NewResponseStructure(hello)
		resp.Json(response)
	default:
		logRequest(request)
		shell.ErrorNotFound(response, request)
	}
}

func logRequest(request *http.Request)  {
	log.Printf("IP: %v Route: / Agent: %v",request.RemoteAddr,request.Header.Get("User-Agent"))
}




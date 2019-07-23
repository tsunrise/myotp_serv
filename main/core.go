package main

import (
	"database/sql"
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

	// install
	if *installMode {
		mydb.Install()
		return
	}

	// initialize database
	db, _, err := mydb.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	if db == nil {
		log.Fatal("Fatal Error: Database instance is nil.")
	}

	log.Println("Started: MyOTP Backend Server Development Edition")
	http.Handle("/", httpServer{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type httpServer struct {
	Database sql.DB
}

func (s httpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// before the url matching
	logRequest(request)

	switch request.URL.Path {
	case "/":
		hello := struct {
			Working bool
		}{true}
		resp := shell.NewResponseStructure(hello)
		resp.Json(response)
	default:
		logRequest(request)
		shell.ErrorNotFound(response, request)
	}
}

func logRequest(request *http.Request) {
	log.Printf("IP: %v Route: / Agent: %v", request.RemoteAddr, request.Header.Get("User-Agent"))
}

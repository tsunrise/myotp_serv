package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"myotp_serv/mydb"
	"net/http"
)

func main() {
	// parse flags
	installMode := flag.Bool("install", false, "Configure this app. ")
	port := flag.Int("port", 8080, "Set the portal. (Default: 8080)")
	flag.Parse()

	// handle install
	if *installMode {
		mydb.Install()
		return
	}

	// initialize database
	db, stmt, err := mydb.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	if db == nil {
		log.Fatal("Fatal Error: Database instance is nil.")
	}

	log.Println("Started: MyOTP Backend Server Development Edition")
	http.Handle("/", httpServer{db, stmt})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

type httpServer struct {
	Database     *sql.DB
	DBStatements *mydb.StatementsSet
}

func (s httpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// before the url matching
	logRequest(request)

	// go to url matching
	indexRouter(s, response, request)
}

func logRequest(request *http.Request) {
	log.Printf("IP: %v Route: / Agent: %v", request.RemoteAddr, request.Header.Get("User-Agent"))
}

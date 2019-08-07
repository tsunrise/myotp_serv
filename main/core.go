package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
)

func main() {
	// parse flags
	installMode := flag.Bool("install", false, "Configure this app. ")
	port := flag.Int("port", 8080, "Set the portal.")
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
	http.Handle("/", httpServer{db, stmt, tokenLib.NewStoreSet()})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

type httpServer struct {
	Database     *sql.DB
	DBStatements *mydb.StatementsSet
	StoreSet     *tokenLib.StoreSet
}

func (s httpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// URL protection
	if len(request.URL.Path) > 1024 || len(request.URL.RequestURI()) > 65535 {
		shell.NewMyError("Request Rejected", "Request URI is too long.", http.StatusBadRequest).Json(response)
		log.Printf("IP: %v %v: %v Agent: %v", request.RemoteAddr, request.Method, "...", request.Header.Get("User-Agent"))
		return
	}

	// before the url matching
	logRequest(request)

	// go to url matching
	indexRouter(s, response, request)
}

func logRequest(request *http.Request) {
	log.Printf("IP: %v %v: %v Agent: %v", request.RemoteAddr, request.Method, request.URL.Path, request.Header.Get("User-Agent"))
}

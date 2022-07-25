package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	r                 *mux.Router
	configurationFile *string
	config            Configuration

	externalToken string
	jitsiKey      string
	jitsiUrl      string
	jitsiIssuer   string
	noVncPassword string
)

type View struct {
	WorkplaceUrl string
	JitsiUrl     string
}

func init() {
	configurationFile = flag.String("f", "assigner.json", "configuration file")
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jwt := r.URL.Query().Get("token")
	authenticated, claims := decodeWebsiteToken(jwt)

	if !authenticated {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Invalid authentication token")
		return
	}

	variables := mux.Vars(r)
	vm := variables["vm"]

	vmUrl, exists := config.Machines[vm]

	if !exists || len(vmUrl) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid machine")
		return
	}

	data := View{
		WorkplaceUrl: vmUrl + "&password=" + noVncPassword,
		JitsiUrl:     "https://" + jitsiUrl + "/" + vm + "?jwt=" + generateJitsiToken(claims),
	}

	template := template.New("view.html")
	template, err := template.ParseFiles("view.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error while parsing template: %s\n", err)
		return
	}

	err = template.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error while executing template: %s\n", err)
		return
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	var exists bool

	config = readConfigurationFile(configurationFile)

	externalToken, exists = os.LookupEnv("EXTERNAL_TOKEN")
	if !exists {
		log.Fatalln("No external token set")
	}

	jitsiKey, exists = os.LookupEnv("SECRET_JITSI_KEY")
	if !exists {
		log.Fatalln("No Jitsi key set")
	}

	jitsiUrl, exists = os.LookupEnv("JITSI_URL")
	if !exists {
		log.Fatalln("No Jitsi URL set")
	}

	jitsiIssuer, exists = os.LookupEnv("JITSI_ISS")
	if !exists {
		log.Fatalln("No Jitsi issuer set")
	}

	noVncPassword, exists = os.LookupEnv("NOVNC_PASSWORD")
	if !exists {
		log.Fatalln("No NoVNC password set")
	}

	r = mux.NewRouter()
	r.HandleFunc("/{vm}", handle)

	log.Println("Assigner is listening on port 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("Assigner failed: %s", err)
	}
}

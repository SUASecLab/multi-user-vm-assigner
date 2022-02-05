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

	websiteKey    string
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
		fmt.Fprintf(w, "Invalid authentication token")
		return
	}

	variables := mux.Vars(r)
	vm := variables["vm"]

	vmUrl, exists := config.Machines[vm]

	if !exists || len(vmUrl) < 1 {
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
		log.Printf("Error while parsing template: %s\n", err)
		return
	}

	err = template.Execute(w, data)
	if err != nil {
		log.Printf("Error while executing template: %s\n", err)
		return
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	config = readConfigurationFile(configurationFile)
	websiteKey = os.Getenv("SECRET_WEBSITE_KEY")
	jitsiKey = os.Getenv("SECRET_JITSI_KEY")
	jitsiUrl = os.Getenv("JITSI_URL")
	jitsiIssuer = os.Getenv("JITSI_ISS")
	noVncPassword = os.Getenv("NOVNC_PASSWORD")

	r = mux.NewRouter()
	r.HandleFunc("/{vm}", handle)

	log.Println("Assigner is listening on port 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("Assigner failed: %s", err)
	}
}

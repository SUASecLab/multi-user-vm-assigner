// multi-user-vm-assigner
//
// This is the API documentation for the multi-user-vm-assigner. It is used in the SUASecLab to assign multiple people to one virtual machine. All people connected to a virtual machine are also connected to a Jitsi Meet room in which the users can talk.
//
// Version: 0.0.1
//
// License: GPL-3.0 https://www.gnu.org/licenses/gpl-3.0.en.html
//
// Contact: Tobias Tefke <t.tefke@stud.fh-sm.de>
//
// Schemes: http
//
// Consumes:
// - text/plain
// Produces:
// - text/html
//
// swagger:meta
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	r                 *mux.Router
	configurationFile *string
	config            Configuration

	domain           string
	extensionsSubdir string
	sidecarUrl       string
	noVncPassword    string
)

type View struct {
	WorkplaceUrl  string
	Workplace2Url string
	JitsiUrl      string
}

func init() {
	configurationFile = flag.String("f", "assigner.json", "configuration file")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	var exists bool

	config = readConfigurationFile(configurationFile)

	sidecarUrl, exists = os.LookupEnv("SIDECAR_URL")
	if !exists {
		log.Fatalln("No sidecar URL set")
	}

	domain, exists = os.LookupEnv("DOMAIN")
	if !exists {
		log.Fatalln("No domain set")
	}

	extensionsSubdir, exists = os.LookupEnv("EXTENSIONS_SUBDIR")
	if !exists {
		log.Fatalln("No extensions subdir set")
	}

	noVncPassword, exists = os.LookupEnv("NOVNC_PASSWORD")
	if !exists {
		log.Fatalln("No NoVNC password set")
	}

	r = mux.NewRouter()
	r.HandleFunc("/{vm}", virtualMachine)

	log.Println("Assigner is listening on port 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("Assigner failed: %s", err)
	}
}

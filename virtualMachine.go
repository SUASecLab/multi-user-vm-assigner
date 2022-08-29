package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	"github.com/gorilla/mux"
)

func virtualMachine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	token := r.URL.Query().Get("token")
	name := url.QueryEscape(r.URL.Query().Get("name"))

	// Validate token
	validationResult, err := extensions.GetValidationResult("http://" + sidecarUrl +
		"/validate?token=" + token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, validationResult.Error)
		log.Println(validationResult.Error, err)
		return
	}

	if !validationResult.Valid {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Invalid authentication token")
		return
	}

	// Get virtual machine
	variables := mux.Vars(r)
	vm := variables["vm"]

	vmUrl, exists := config.Machines[vm]

	if !exists || len(vmUrl) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid machine")
		return
	}

	// Generate Jitsi token
	jitsiIssuance, err := extensions.IssueToken("http://" + sidecarUrl +
		"/issuance?token=" + token + "&name=" + name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jitsiIssuance.Error)
		log.Println(jitsiIssuance.Error)
		return
	}
	data := View{
		WorkplaceUrl: vmUrl + "&password=" + noVncPassword,
		JitsiUrl:     "https://" + jitsiUrl + "/" + vm + "?jwt=" + jitsiIssuance.Token,
	}

	template := template.New("view.html")
	template, err = template.ParseFiles("view.html")
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

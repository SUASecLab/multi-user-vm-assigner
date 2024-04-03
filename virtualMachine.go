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
	// swagger:operation GET /{vm} getVm
	//
	// Access a virtual machine.
	//
	// Request access to the virtual machine having the handed over name.
	// ---
	// produces:
	// - text/html
	// parameters:
	// - name: vm
	//   in: path
	//   description: The name of the virtual machine which should be accessed.
	//   required: true
	//   type: string
	// - name: vm2
	//   in: query
	//   description: The name of a second virtual machine which should be accessed.
	//   required: false
	//   type: string
	// - name: token
	//   in: query
	//   description: The JWT token of the user handed out by the Sidecar, used for authentication.
	//   required: true
	//   type: string
	//   format: JWT
	// - name: name
	//   in: query
	//   description: The name of the user to be displayed in the Jitsi Meet room.
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: The provided data is correct, access to the virtual machine can be granted. The user directly receives the HTML code for a website embedding the virtual machine and a Jitsi Meet conference, in which all people connected to the same virtual machine meet.
	//     schema:
	//       type: file
	//       description: Website in which the virtual machine and a Jitsi Meet conference with all users connected to the specific virtual machine are embedded.
	//   "400":
	//     description: The request could not be processed because there is no virtual machine having the handed over name.
	//   "403":
	//     description: The request could not be processed because the user's access token is not valid.
	//   "500":
	//     description: The request could not be processed due to an internal error. Maybe the Sidecar experienced an error.
	w.Header().Set("Content-Type", "text/html")

	token := r.URL.Query().Get("token")
	name := url.QueryEscape(r.URL.Query().Get("name"))
	vm2 := r.URL.Query().Get("vm2")

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

	// check if second vm exists (if handed over)
	var vm2Url string
	if len(vm2) > 0 {
		vm2Url, exists = config.Machines[vm2]
		if !exists || len(vm2Url) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid second machine")
			return
		} else {
			vm2Url = vm2Url + "&password=" + noVncPassword
		}
	}

	// Open Jitsi room
	data := View{
		WorkplaceUrl:  vmUrl + "&password=" + noVncPassword,
		Workplace2Url: vm2Url,
		JitsiUrl: "https://" + domain + extensionsSubdir +
			"/jitsi/?roomName=" + vm + "&userName=" + name + "&token=" + token,
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

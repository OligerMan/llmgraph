package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setInstanceCapabilities(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "setInstanceCapabilities\n")

	CheckHTTPRequest(w, req)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var instanceCap InstanceCapability
	if err := dec.Decode(&instanceCap); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	diskHandler := createDiskHandler(DefaultCapabilityPath)
	if diskHandler != nil {
		diskHandler.setInstanceCapability(instanceCap)
		diskHandler.deleteDiskHandler()
		status := ReturnStatus{"connected"}
		_, err := w.Write(JSONWrapper(status))
		if err != nil {
			return
		}
	} else {
		fmt.Printf("Disk Handler not created\n")
		status := ReturnStatus{"server_error"}
		_, err := w.Write(JSONWrapper(status))
		if err != nil {
			return
		}
	}
}

func deleteInstance(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "deleteInstance\n")

	CheckHTTPRequest(w, req)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	type deleteInstanceReq struct {
		InstanceId string `json:"instance_id"`
	}
	var instanceReq deleteInstanceReq
	if err := dec.Decode(&instanceReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	diskHandler := createDiskHandler(DefaultCapabilityPath)
	if diskHandler != nil {
		diskHandler.deleteInstanceCapability(instanceReq.InstanceId)
		diskHandler.deleteDiskHandler()
	}
}

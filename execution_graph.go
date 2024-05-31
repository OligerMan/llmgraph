package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setExecutionGraph(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "setExecutionGraph\n")

	CheckHTTPRequest(w, req)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var executionGraph ExecutionGraphStructure
	if err := dec.Decode(&executionGraph); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	diskHandler := createDiskHandler(DefaultExecutionGraphPath)
	if diskHandler != nil {
		if verifyExecutionGraph(executionGraph) {
			diskHandler.setExecutionGraph(executionGraph)
			diskHandler.deleteDiskHandler()
			status := ReturnStatus{"connected"}
			_, err := w.Write(JSONWrapper(status))
			if err != nil {
				return
			}
		} else {
			fmt.Printf("Graph verification failed\n")
			status := ReturnStatus{"config_error"}
			_, err := w.Write(JSONWrapper(status))
			if err != nil {
				return
			}
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

func deleteExecutionGraph(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "deleteExecutionGraph\n")

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

	diskHandler := createDiskHandler(DefaultExecutionGraphPath)
	if diskHandler != nil {
		diskHandler.deleteExecutionGraph(instanceReq.InstanceId)
		diskHandler.deleteDiskHandler()
	}
}

func verifyExecutionGraph(executionGraph ExecutionGraphStructure) bool {
	// for i, link := range executionGraph.Links {
	// }
	return true
}

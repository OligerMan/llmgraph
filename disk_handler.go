package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var DefaultCapabilityPath = "./capabilities/"
var DefaultExecutionGraphPath = "./execution_graph/"
var CapabilityExtension = ".json"

type DiskHandler struct {
	fileptr  *os.File
	mainPath string
}

func createDiskHandler(MainPath string) *DiskHandler {
	diskHandler := &DiskHandler{}
	err := os.Mkdir(MainPath, 0777)
	file, err := os.Open(MainPath)
	if err != nil {
		return nil
	}
	diskHandler.fileptr = file
	diskHandler.mainPath = MainPath
	return diskHandler
}

func (diskHandler *DiskHandler) deleteDiskHandler() {
	err := diskHandler.fileptr.Close()
	if err != nil {
		return
	}
}

func (diskHandler *DiskHandler) setInstanceCapability(instCap InstanceCapability) {
	var instPath = diskHandler.mainPath + instCap.InstanceID + CapabilityExtension
	var file *os.File
	if _, err := os.Stat(instPath); err != nil {
		file, err = os.Create(instPath)
		if err != nil {
			return
		}
	} else {
		file, _ = os.Open(instPath)
	}
	_, err := file.Write(JSONWrapper(instCap))
	if err != nil {
		return
	}
	err = file.Close()
	if err != nil {
		return
	}
}

func (diskHandler *DiskHandler) getInstanceCapability(instanceID string) InstanceCapability {
	var instCap InstanceCapability
	var instPath = diskHandler.mainPath + instanceID + CapabilityExtension
	file, err := os.Open(instPath)
	if err == nil {
		instCap.InstanceID = ""
		return instCap
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&instCap)
	if err != nil {
		return InstanceCapability{}
	}
	return instCap
}

func (diskHandler *DiskHandler) deleteInstanceCapability(instanceID string) {
	var instPath = diskHandler.mainPath + instanceID + CapabilityExtension
	if _, err := os.Stat(instPath); err == nil {
		err := os.Remove(instPath)
		if err != nil {
			fmt.Printf("deleteInstanceCapability error - %s\n", err)
		}
	} else {
		fmt.Printf("deleteInstanceCapability error - file do not exist")
	}
}

func (diskHandler *DiskHandler) setExecutionGraph(execGraph ExecutionGraphStructure) {
	var instPath = diskHandler.mainPath + execGraph.GraphID + CapabilityExtension
	var file *os.File
	if _, err := os.Stat(instPath); err != nil {
		file, err = os.Create(instPath)
		if err != nil {
			return
		}
	} else {
		file, _ = os.Open(instPath)
	}
	_, err := file.Write(JSONWrapper(execGraph))
	if err != nil {
		return
	}
	err = file.Close()
	if err != nil {
		return
	}
}

func (diskHandler *DiskHandler) getExecutionGraph(execGraphID string) ExecutionGraphStructure {
	var execGraph ExecutionGraphStructure
	var instPath = diskHandler.mainPath + execGraphID + CapabilityExtension
	file, err := os.Open(instPath)
	if err != nil {
		execGraph.GraphID = ""
		return execGraph
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&execGraph)
	if err != nil {
		return ExecutionGraphStructure{}
	}
	return execGraph
}

func (diskHandler *DiskHandler) deleteExecutionGraph(execGraphID string) {
	var instPath = diskHandler.mainPath + execGraphID + CapabilityExtension
	if _, err := os.Stat(instPath); err == nil {
		err := os.Remove(instPath)
		if err != nil {
			fmt.Printf("deleteInstanceCapability error - %s\n", err)
		}
	} else {
		fmt.Printf("deleteInstanceCapability error - file do not exist")
	}
}

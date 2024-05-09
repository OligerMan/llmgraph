package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExecState struct {
	ID            string
	currentInput  []ValueWrapper
	currentOutput []ValueWrapper
}

func ExecuteState(state ExecState) ([]ValueWrapper, error) {
	sendToRabbit(state.currentInput, state.ID, nil, nil)
	obj, err := getFromRabbit(state.ID+"_response", nil, nil)
	if err != nil {
		fmt.Printf("Executing state %s error\n", state.ID)
		return nil, err
	}
	var output []ValueWrapper
	err = json.Unmarshal(obj, &output)
	if err != nil {
		fmt.Printf("Executing state %s response decoding error\n", state.ID)
		return nil, err
	}
	return output, nil
}

func ExecutionCycle(executeConfig ExecuteConfig) []ValueWrapper {
	id := executeConfig.ExecutionGraphID

	disk := createDiskHandler(DefaultExecutionGraphPath)
	execGraph := disk.getExecutionGraph(id)
	var inputToLink map[string][]Link
	var outputToLink map[string][]Link
	for i := 0; i < len(execGraph.Links); i++ {
		inputToLink[execGraph.Links[i].InputID] = append(inputToLink[execGraph.Links[i].InputID], execGraph.Links[i])
		outputToLink[execGraph.Links[i].OutputID] = append(outputToLink[execGraph.Links[i].OutputID], execGraph.Links[i])
	}

	var execQueue []ExecState
	queued := make(map[string]struct{})
	completed := make(map[string]ExecState)
	execQueue = append(execQueue, ExecState{execGraph.StartID, executeConfig.InputData, []ValueWrapper{}})
	queued[execGraph.StartID] = struct{}{}

	for i := 0; i < len(execQueue); i++ {
		currentOutput, err := ExecuteState(execQueue[i])
		if err != nil {
			fmt.Printf("Executing cycle on state %s error\n", execQueue[i].ID)
			continue
		}
		execQueue[i].currentOutput = currentOutput
		completed[execQueue[i].ID] = execQueue[i]

		inputLinksData := inputToLink[execQueue[i].ID]
		for i := 0; i < len(inputLinksData); i++ {
			outputID := inputLinksData[i].OutputID
			outputLinksData := outputToLink[outputID]

			all_completed := true
			var totalInput []ValueWrapper
			for j := 0; j < len(outputLinksData); j++ {
				execStateChosen, inside := completed[outputLinksData[j].InputID]
				if !inside {
					all_completed = false
					break
				} else {
					totalInput = append(totalInput, execStateChosen.currentOutput...)
				}
			}
			if all_completed {
				queued[outputID] = struct{}{}
				execQueue = append(execQueue, ExecState{outputID, totalInput, []ValueWrapper{}})
			}
		}
	}
	return execQueue[len(execQueue)-1].currentOutput
}

func execute(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "execute\n")

	CheckHTTPRequest(w, req)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var execConfig ExecuteConfig
	if err := dec.Decode(&execConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	executionAnswer := ExecutionCycle(execConfig)
	for i := 0; i < len(executionAnswer); i++ {
		fmt.Fprintf(w, "%s\n", executionAnswer[i].ValueInfo)
	}
}

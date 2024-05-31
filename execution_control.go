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

func ExecuteState(state ExecState, handlerID string) ([]ValueWrapper, error) {
	curInput := state.currentInput
	for {
		sendToRabbit(curInput, handlerID, nil, nil)
		obj, err := getFromRabbit(handlerID+"_response", nil, nil)
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
		list_output := false
		list_handle := false
		for i := 0; i < len(output); i++ {
			wrapper_id := output[i].WrapperID
			if len(wrapper_id) > 5 && wrapper_id[len(wrapper_id)-5:] == "_list" {
				list_output = true
			}
			if wrapper_id == "list_process_mode" {
				list_handle = true
			}
		}
		if list_output && list_handle {
			curInput = output
			continue
		}
		return output, nil
	}
}

func ExecutionCycle(executeConfig ExecuteConfig) []ValueWrapper {
	id := executeConfig.ExecutionGraphID

	disk := createDiskHandler(DefaultExecutionGraphPath)
	execGraph := disk.getExecutionGraph(id)
	inputToLink := make(map[string][]Link)
	outputToLink := make(map[string][]Link)
	for i := 0; i < len(execGraph.Links); i++ {
		inputToLink[execGraph.Links[i].InputID] = append(inputToLink[execGraph.Links[i].InputID], execGraph.Links[i])
		outputToLink[execGraph.Links[i].OutputID] = append(outputToLink[execGraph.Links[i].OutputID], execGraph.Links[i])
	}

	nodeIDToHandlerID := make(map[string]string)
	for i := 0; i < len(execGraph.RequestConfigs); i++ {
		nodeIDToHandlerID[execGraph.RequestConfigs[i].ID] = execGraph.RequestConfigs[i].HandlerID
	}

	var execQueue []ExecState
	queued := make(map[string]struct{})
	completed := make(map[string]ExecState)
	for i := 0; i < len(execGraph.Constants); i++ {
		tmp_state := ExecState{execGraph.Constants[i].ID, []ValueWrapper{}, execGraph.Constants[i].Values}
		execQueue = append(execQueue, tmp_state)
		completed[execGraph.Constants[i].ID] = tmp_state
	}
	execQueue = append(execQueue, ExecState{execGraph.StartID, executeConfig.InputData, []ValueWrapper{}})
	queued[execGraph.StartID] = struct{}{}

	for i := len(execGraph.Constants); i < len(execQueue); i++ {
		currentOutput, err := ExecuteState(execQueue[i], nodeIDToHandlerID[execQueue[i].ID])
		if err != nil {
			fmt.Printf("Executing cycle on state %s error\n", nodeIDToHandlerID[execQueue[i].ID])
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
					reworkedOutput := execStateChosen.currentOutput
					for k := 0; k < len(reworkedOutput); k++ {
						for l := 0; l < len(outputLinksData[j].Replacements); l++ {
							if reworkedOutput[i].WrapperID == outputLinksData[j].Replacements[l].Input {
								reworkedOutput[i].WrapperID = outputLinksData[j].Replacements[l].Output
							}
						}
					}
					totalInput = append(totalInput, reworkedOutput...)
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
		fmt.Fprintf(w, "%s", string(JSONWrapper(executionAnswer)))
	}
}

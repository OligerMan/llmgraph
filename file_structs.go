package main

type IOType struct {
	DataType   string `json:"data_type"`
	ID         string `json:"data_type_id"`
	ListStatus string `json:"list_status"`
}

type InstanceCapability struct {
	InstanceID   string   `json:"instance_id"`
	InstanceName string   `json:"instance_name"`
	InputType    []IOType `json:"input_type"`
	OutputType   []IOType `json:"output_type"`
	HandlerID    string   `json:"handler_id"`
}

type ReturnStatus struct {
	Status string `json:"status"`
}

type RequestConfig struct {
	ID         string   `json:"id"`
	InputType  []IOType `json:"input_type"`
	OutputType []IOType `json:"output_type"`
	HandlerID  string   `json:"handler_id"`
}

type ValueWrapper struct {
	WrapperID string `json:"wrapper_id"`
	ValueType string `json:"value_type"`
	ValueInfo string `json:"value_info"`
}

type ConstantConfig struct {
	ID     string         `json:"id"`
	Values []ValueWrapper `json:"values"`
}

type FieldReplacement struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Link struct {
	InputID          string             `json:"input_id"`
	InputDataTypeID  string             `json:"input_data_type_id"`
	OutputID         string             `json:"output_id"`
	OutputDataTypeID string             `json:"output_data_type_id"`
	Replacements     []FieldReplacement `json:"replacements"`
}

type ExecutionGraphStructure struct {
	GraphID        string           `json:"graph_id"`
	GraphName      string           `json:"graph_name"`
	RequestConfigs []RequestConfig  `json:"request_configs"`
	Links          []Link           `json:"links"`
	Constants      []ConstantConfig `json:"constants"`
	StartID        string           `json:"start_id"`
}

type ExecuteConfig struct {
	ExecutionGraphID string         `json:"execution_graph_id"`
	InputData        []ValueWrapper `json:"input_data"`
}

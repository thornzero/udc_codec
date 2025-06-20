package api

type APITag struct {
	FullTag    string `json:"full_tag"`
	SystemCode string `json:"system_code"`
	FunctionCode string `json:"function_code"`
	EquipmentID string `json:"equipment_id"`
	Description string `json:"description"`
	UDCCode    string `json:"udc_code,omitempty"`
}

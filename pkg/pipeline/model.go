package pipeline

type BOMEntry struct {
	SystemCode  string `yaml:"system_code"`
	EquipmentID string `yaml:"equipment_id"`
	FunctionCode string `yaml:"function_code"`
	UDCCode     string `yaml:"udc_code,omitempty"`
	Description string `yaml:"description"`
}

type ProjectBOM struct {
	ProjectName string      `yaml:"project_name"`
	Entries     []BOMEntry  `yaml:"entries"`
}

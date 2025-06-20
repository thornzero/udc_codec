package aggregator

type AggregatedSystem struct {
	SystemCode string `yaml:"system_code"`
	SystemName string `yaml:"system_name"`
	UDCCode    string `yaml:"udc_code,omitempty"`
	IECCode    string `yaml:"iec_code,omitempty"`
	ISAFunction map[string]string `yaml:"isa_function,omitempty"`
}

type AggregatedDatabase struct {
	Systems []AggregatedSystem `yaml:"systems"`
}

package udc

type UDCNode struct {
	Code     string     `yaml:"code"`
	Title    string     `yaml:"title"`
	Children []*UDCNode `yaml:"children,omitempty"`
}

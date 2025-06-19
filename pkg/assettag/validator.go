package assettag

import "fmt"

func (r *Resolver) ValidateTag(tag *Tag) error {
	if _, ok := r.IEC81346[tag.SystemCode]; !ok {
		return fmt.Errorf("unknown system code: %s", tag.SystemCode)
	}
	if _, ok := r.ISA[tag.FunctionCode]; !ok {
		return fmt.Errorf("unknown ISA function code: %s", tag.FunctionCode)
	}
	if _, ok := r.UDC.Lookup(tag.UDCCode); !ok && tag.UDCCode != "" {
		return fmt.Errorf("unknown UDC code: %s", tag.UDCCode)
	}
	return nil
}

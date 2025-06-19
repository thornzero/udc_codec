package assettag

import "fmt"

func (r *Resolver) DescribeTag(tag *Tag) string {
	sys := r.IEC81346[tag.SystemCode]
	isa := r.ISA[tag.FunctionCode]
	udc := ""
	if tag.UDCCode != "" {
		if desc, ok := r.UDC.Lookup(tag.UDCCode); ok {
			udc = fmt.Sprintf(" (%s)", desc)
		}
	}
	return fmt.Sprintf("%s: %s %s %s%s", sys, isa, tag.EquipmentID, tag.InstrumentID, udc)
}

package assettag

import (
	"github.com/thornzero/udc_codec/pkg/udc"
)

type Tag struct {
	SystemCode   string // IEC 81346 system
	EquipmentID  string // Equipment unique ID
	InstrumentID string // Instrument unique ID
	FunctionCode string // ISA-5.1 function letters
	UDCCode      string // UDC functional context
}

type Resolver struct {
	UDC      *udc.Codec
	ISA      map[string]string
	IEC81346 map[string]string
}

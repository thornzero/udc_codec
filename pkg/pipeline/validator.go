package pipeline

import (
	"fmt"

	"github.com/thornzero/udc_codec/pkg/aggregator"
	"github.com/thornzero/udc_codec/pkg/udc"
)

type Validator struct {
	Aggregator *aggregator.AggregatedDatabase
	UDC        *udc.Codec
}

func (v *Validator) ValidateEntry(entry BOMEntry) error {
	sys := v.Aggregator.LookupSystem(entry.SystemCode)
	if sys == nil {
		return fmt.Errorf("unknown system code: %s", entry.SystemCode)
	}

	if _, ok := sys.ISAFunction[entry.FunctionCode]; !ok {
		return fmt.Errorf("invalid ISA function code %s for system %s", entry.FunctionCode, sys.SystemCode)
	}

	if entry.UDCCode != "" {
		if _, ok := v.UDC.Lookup(entry.UDCCode); !ok {
			return fmt.Errorf("invalid UDC code %s", entry.UDCCode)
		}
	}
	return nil
}

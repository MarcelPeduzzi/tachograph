package vu

import (
	"github.com/way-platform/tachograph-go/internal/dd"
)

// MarshalOptions configures the marshaling of VU files into binary format.
type MarshalOptions struct {
	// Embed dd.MarshalOptions to inherit marshaling configuration.
	dd.MarshalOptions
}

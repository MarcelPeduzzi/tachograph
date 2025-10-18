// Package certcache provides an embedded cache of certificates.
package certcache

import "embed"

//go:embed root/EC_PK.bin
var rootG1 []byte

//go:embed root/ERCA_Gen2_Root_Certificate.bin
var rootG2 []byte

//go:embed g1/*.bin
var g1 embed.FS

//go:embed g2/*.bin
var g2 embed.FS

// RootG1 returns the Gen1 ERCA root certificate (RSA).
func RootG1() []byte {
	return rootG1
}

// RootG2 returns the Gen2 ERCA root certificate (ECC).
func RootG2() []byte {
	return rootG2
}

// Root returns the Gen1 ERCA root certificate for backward compatibility.
// Deprecated: Use RootG1() or RootG2() explicitly.
func Root() []byte {
	return rootG1
}

// ReadG1 reads a cached Gen1 certificate by its CHR.
func ReadG1(chr string) ([]byte, bool) {
	data, err := g1.ReadFile("g1/" + chr + ".bin")
	if err != nil {
		return nil, false
	}
	return data, true
}

// ReadG2 reads a cached Gen2 certificate by its CHR.
func ReadG2(chr string) ([]byte, bool) {
	data, err := g2.ReadFile("g2/" + chr + ".bin")
	if err != nil {
		return nil, false
	}
	return data, true
}

package onetimeauth

import "testing"

func TestComplete(t *testing.T) {
	options := &Options{}
	options.complete()

	if options.Issuer != DefaultOptions.Issuer {
		t.Error("Issuer not defaulted")
	}
	if options.Lifetime != DefaultOptions.Lifetime {
		t.Error("Lifetime not defaulted")
	}
	if options.SigningKeyLength != DefaultOptions.SigningKeyLength {
		t.Error("SigningKeyLength not defaulted")
	}
	if options.TokenKeyLength != DefaultOptions.TokenKeyLength {
		t.Error("TokenKeyLength not defaulted")
	}
	if options.SigningMethod != DefaultOptions.SigningMethod {
		t.Error("SigningMethod not defaulted")
	}
}

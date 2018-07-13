package didi_json

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var DIDIJSON = jsoniter.ConfigCompatibleWithStandardLibrary

func FixJSION() {
	extra.RegisterFuzzyDecoders() // php json tolerant
}

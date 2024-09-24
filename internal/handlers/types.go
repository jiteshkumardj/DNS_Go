package handlers

import "errors"

var errNoOutputGiven = errors.New("no known output given")

type outputType string

const (
	outputSTDOUT outputType = "stdout"
	outputGRPC   outputType = "grpc"
)

var knownOUTS = map[outputType]struct{}{
	outputGRPC:   {},
	outputSTDOUT: {},
}
var listOUTS = []string{
	string(outputGRPC),
	string(outputSTDOUT),
}

// @todo: validate flags
func IsKnownOUT(s string) bool {
	_, ok := knownOUTS[outputType(s)]
	return ok
}

func ListKnownOutputs() []string {
	return listOUTS
}

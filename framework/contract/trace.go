package contract

import "context"

const TraceKey = "hade:trace"

type Trace interface {
	GetTrace(ctx context.Context) interface{}
	ToMap(interface{}) map[string]interface{}
}

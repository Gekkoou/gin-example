package drive

import (
	"strings"
)

type Drive struct {
	Name        string
	FailureName string
}

func getQueueNames(topic, prefix, failureSuffix string) (name, failureName string) {
	var builder strings.Builder
	builder.WriteString(prefix)
	builder.WriteString(topic)
	name = builder.String()
	builder.WriteString(failureSuffix)
	failureName = builder.String()
	return
}

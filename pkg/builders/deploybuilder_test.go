package builders

import (
	"fmt"
	"testing"
)

func TestDeployBuilder_Build(t *testing.T) {
	builder, err := NewDeployBuilder("defualt", "abc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(builder.Replicas(5).Build())
}

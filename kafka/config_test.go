package kafka

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTopic_String(t *testing.T) {
	assert.Equal(t, "hello", Topic("hello").String())
}

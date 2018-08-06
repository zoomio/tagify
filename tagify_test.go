package tagify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	tags, err := GetTags("http://stackoverflow.com", HTML, 10, false, false)
	assert.Nil(t, err)
	assert.Len(t, tags, 10)
}

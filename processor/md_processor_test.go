package processor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mdSmallText = `
	There was a boy
	Whose name was Jim.
`

	mdMediumText = `
	# A story about Jim
	
	There was a boy whose name was **Jim**.	His Friends were very good to him.
	They gave him Tea, and Cakes, and Jam, And slices of delicious Ham...
`
)

// table driven tests
var parseMDTests = []struct {
	name        string
	text        string
	noStopWords bool
	tags        []string
	title       string
	hash        string
}{
	{
		"small",
		mdSmallText,
		false,
		[]string{"was", "there", "a", "boy", "whose", "name", "jim"},
		"",
		"cd75ea8c309f25c6774d4d7cebbf22932de578ffdd0e48afe30f658613745440e48e91c671ab23a019dc6d0d54f9048c7920f94014061be09c6b10590c937b1b",
	},
	{
		"medium",
		mdMediumText,
		true,
		[]string{"boy", "friends", "tea", "ham", "story", "jim", "good", "cakes", "jam", "slices", "delicious"},
		"A story about Jim",
		"4b07663d525c70edf6a6519e05869c2e0d45c93ef4f397f8026c98408c2ff607ddec1cd0a6f104eff90317ec1dc34b7a9be8580f3ebc2d8892fc50ed755f478d",
	},
}

func Test_ParseMD(t *testing.T) {
	for _, tt := range parseMDTests {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseMD(&inputReadCloser{strings.NewReader(tt.text)}, NoStopWords(tt.noStopWords))
			assert.Equal(t, tt.title, out.DocTitle)
			assert.Equal(t, tt.hash, fmt.Sprintf("%x", out.DocHash))
			assert.ElementsMatch(t, tt.tags, ToStrings(out.Tags))
		})
	}
}

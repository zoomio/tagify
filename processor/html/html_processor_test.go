package html

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
	"github.com/zoomio/tagify/model"
)

var (
	simpleHTML, _       = os.ReadFile("../../_resources_test/html/simple.html")
	doubledTitleHTML, _ = os.ReadFile("../../_resources_test/html/doubled-title.html")
	complexHTML, _      = os.ReadFile("../../_resources_test/html/complex.html")
	complexTextHTML, _  = os.ReadFile("../../_resources_test/html/complex-text.html")
	cssyHTML, _         = os.ReadFile("../../_resources_test/html/cssy.html")
	theVergeHTML, _     = os.ReadFile("../../_resources_test/html/theverge.html")
)

type inputReadCloser struct {
	io.Reader
}

func (in *inputReadCloser) Close() error {
	return nil
}

// table driven tests
var processHTMLTests = []struct {
	name        string
	in          string
	expect      []string
	title       string
	hash        string
	noStopWords bool
	contentOnly bool
}{
	{
		"empty",
		"",
		[]string{},
		"",
		"",
		false,
		true,
	},
	{
		"simple",
		string(simpleHTML),
		[]string{"there", "was", "a", "boy", "whose", "name", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		false,
		true,
	},
	{
		"simple exclude stopWords",
		string(simpleHTML),
		[]string{"boy", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		true,
		true,
	},
	{
		"complex",
		string(complexHTML),
		[]string{"go", "golang", "html", "extract", "all", "certain", "parse", "content", "from", "tags", "theme", "blog", "help"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"04c09437103091df65d3c8d464017156bd181951adf614bccf30c5b40332641a7bd3d9a3a5042119d9e72312e2ce545c4522a546f9e869fe5d0c2dc6c988ab13",
		false,
		true,
	},
	{
		"complex exclude stopWords",
		string(complexHTML),
		[]string{"parse", "html", "extract", "content", "tags", "theme", "golang", "blog", "help"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"04c09437103091df65d3c8d464017156bd181951adf614bccf30c5b40332641a7bd3d9a3a5042119d9e72312e2ce545c4522a546f9e869fe5d0c2dc6c988ab13",
		true,
		true,
	},
	{
		"complex exclude stopWords tag everything",
		string(complexHTML),
		[]string{"tags", "help", "blog", "html", "content", "extract", "theme", "golang", "parse"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"04c09437103091df65d3c8d464017156bd181951adf614bccf30c5b40332641a7bd3d9a3a5042119d9e72312e2ce545c4522a546f9e869fe5d0c2dc6c988ab13",
		true,
		false,
	},
	{
		"css-y",
		string(cssyHTML),
		[]string{"stuff", "foo", "texty", "text", "people", "cool"},
		"People are looking for cool stuff",
		"09e63717d8ea919f68c3f8cc9403ebe5d119baf924e3bb0d7e7db7d317f6c3ba46f1319da2857f0fe965ff06a4bb5ee17e35bdd1c16d2402b8a5a6d3748b49e4",
		true,
		false,
	},
	{
		"meta description",
		string(theVergeHTML),
		[]string{"longer", "title", "verge"},
		"Hi This is Slightly Longer Title",
		"13ea1c679ec7d1678d60b614f595192c47907fed5ea0e2883de001e3e2bcfd4fea61dea4a9cfa5f9a8f91a575c181a582577987d7938e062b1820da80cfb64dd",
		true,
		false,
	},
	{
		"complex text",
		string(complexTextHTML),
		[]string{"document", "tags", "funny", "thing", "funtivity", "fun", "drag", "extra", "complex", "text", "testing", "hussle", "yup", "lotsa"},
		"Complex text line",
		"8200fbd4839ec87a58faf5eb889cbf1542b18645fba96230a171daa27175a282b0109215fdab04692b148e6908c1a72eb9f1ff969d975eaba16b6a394f6559bd",
		true,
		true,
	},
}

func Test_ProcessHTML(t *testing.T) {
	for _, tt := range processHTMLTests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New(config.NoStopWords(tt.noStopWords), config.ContentOnly(tt.contentOnly))
			out := ProcessHTML(cfg, &inputReadCloser{strings.NewReader(tt.in)})
			assert.Equal(t, tt.title, out.Meta.DocTitle)
			assert.Equal(t, tt.hash, out.Meta.DocHash)
			assert.ElementsMatch(t, tt.expect, model.ToStrings(out.Flatten()))
		})
	}
}

/* func Test_ProcessHTML_DedupeTitleAndHeading(t *testing.T) {
	out := ProcessHTML(config.New(config.NoStopWords(true)), &inputReadCloser{bytes.NewReader(doubledTitleHTML)})
	assert.Equal(t, "A story about a boy", out.Meta.DocTitle)
	assert.Equal(t,
		"4f652c47205d3b922115eef155c484cf81096351696413c86277fa0ed89ebfefe30f81ef6fc6a9d7d654a9292c3cb7aa6f3696052e53c113785a9b1b3be7d4a8",
		out.Meta.DocHash)
	assert.Contains(t, out.Flatten(), &model.Tag{Value: "story", Score: defaultTagWeights[atom.Title.String()], Count: 1, Docs: 1, DocsCount: 4})
} */

func Test_ProcessHTML_NoSpecificStopWords(t *testing.T) {
	out := ProcessHTML(config.New(config.NoStopWords(true)), &inputReadCloser{bytes.NewReader(doubledTitleHTML)})
	assert.Equal(t, "A story about a boy", out.Meta.DocTitle)
	assert.Equal(t,
		"4f652c47205d3b922115eef155c484cf81096351696413c86277fa0ed89ebfefe30f81ef6fc6a9d7d654a9292c3cb7aa6f3696052e53c113785a9b1b3be7d4a8",
		out.Meta.DocHash)
	assert.NotContains(t, out.Flatten(), &model.Tag{Value: "part", Score: 1.4, Count: 1})
}

func Test_ParseHTML(t *testing.T) {
	const htmlPage = `
	<html>
	<body>
	<p>There was a boy <b>whose</b> name was Jim.</p>
	</body>
	</html>
`
	contents := ParseHTML(
		&inputReadCloser{strings.NewReader(htmlPage)},
		&config.Config{TagWeights: defaultTagWeights},
		nil,
		nil,
	)
	assert.NotNil(t, contents)

	assert.Len(t, contents.lines, 1)

	line := contents.lines[0]
	assert.Len(t, line.parts, 3)

	assert.Equal(t, atom.P.String(), line.parts[0].tag)
	assert.Equal(t, "There was a boy ", string(line.pData(line.parts[0])))

	assert.Equal(t, atom.B.String(), line.parts[1].tag)
	assert.Equal(t, "whose", string(line.pData(line.parts[1])))

	assert.Equal(t, atom.P.String(), line.parts[2].tag)
	assert.Equal(t, " name was Jim.", string(line.pData(line.parts[2])))
}

func Test_ParseReaderHTML_visits_all_tags(t *testing.T) {
	counter := &testCountingExt{BaseExtension: extension.NewExtension("testCountingExt", "1")}
	contents := ParseHTML(
		io.NopCloser(bytes.NewReader(theVergeHTML)),
		&config.Config{Verbose: false, SkipLang: true, AllTagWeights: true},
		[]HTMLExt{counter},
		nil,
	)

	assert.NotNil(t, contents)
	assert.Len(t, contents.lines, 2)
	assert.Equal(t, 37, counter.count)
}

type testCountingExt struct {
	*extension.BaseExtension
	count int
}

func (ext *testCountingExt) ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error) {
	ext.count++
	return false, nil
}

// table driven tests
var isSameDomainTests = []struct {
	name     string
	href     string
	domain   string
	expected bool
}{
	{
		"same",
		"https://zoomio.org/tagify",
		"https://zoomio.org",
		true,
	},
	{
		"different scheme",
		"http://zoomio.org/tagify",
		"https://zoomio.org",
		true,
	},
	{
		"subdomain",
		"http://api.zoomio.org/api/tagify",
		"https://zoomio.org",
		true,
	},
	{
		"path",
		"/tagify",
		"https://zoomio.org",
		true,
	},
	{
		"different",
		"https://google.com",
		"https://zoomio.org",
		false,
	},
	{
		"one letter diff",
		"https://zoomioo.org",
		"https://zoomio.org",
		false,
	},
}

func Test_isSameDomain(t *testing.T) {
	for _, tt := range isSameDomainTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSameDomain(tt.href, tt.domain))
		})
	}
}

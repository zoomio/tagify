package tagify

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
	"github.com/zoomio/tagify/model"
	thtml "github.com/zoomio/tagify/processor/html"
)

var ctx = context.TODO()

// table driven tests
var runTests = []struct {
	name        string
	in          []Option
	expectTags  []string
	expectTitle string
	expectHash  string
}{
	{
		"run",
		[]Option{Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML), Limit(5), NoStopWords(true), ContentOnly(true)},
		[]string{"test", "boy", "cakes", "chocolate", "delicious"},
		"Test",
		"d91b531c92a14fe5556b9bf3e82ef6dac0da69914affca86795181d2ca9ca3046630a41eba734b146eec0c5e78c5780734dfb42bcef453cf1d9d20830b562dac",
	},
	{
		"run with query",
		[]Option{Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML), Limit(5), NoStopWords(true),
			Query("#box3 p"), ContentOnly(true)},
		[]string{"bang", "began", "boy", "day", "eat"},
		"",
		"e5e0aef65e77e87a3e23a3f157357444910f94f5dccd5d0fe185da73cb72a8b7bff6ac80d71cfca1da27e9d1b7a3e810a348ceeee52c2e4b68393c8ba5d92cc4",
	},
	{
		"run custom weights",
		[]Option{Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML), Limit(5), NoStopWords(true), TagWeightsString("title:3")},
		[]string{"test"},
		"Test",
		"20c62640489dbc272c51abfd1fbe7b5aa7280f814fbfdb2baf993fb1e8b4c860fb1f1c6964760144e2ef15849ef073f47cb89284481d17845565395d7574e2e7",
	},
}

func Test_Run_HTML(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port), indexHTML))
	for _, tt := range runTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Run(ctx, tt.in...)
			assert.Nil(t, err)
			assert.Equal(t, HTML, res.Meta.ContentType)
			assert.Equal(t, tt.expectTitle, res.Meta.DocTitle)
			assert.Equal(t, tt.expectHash, res.Meta.DocHash)
			assert.ElementsMatch(t, tt.expectTags, res.TagsStrings())
		})
	}
}

func Test_GetTagsFromString(t *testing.T) {
	res, err := Run(ctx,
		Content("Test input reader of type text"),
		TargetType(Text),
		Limit(3),
		NoStopWords(true),
	)
	assert.Nil(t, err)
	assert.Len(t, res.Tags, 3)
	assert.Equal(t,
		"7d95ed3e8436c978f3e7f19f1645f89091f9fdb0439c15547f0a6f82bc4a0babebd06ff6285d9dff8db77861edf2cc8e6919ea5613bec0f30dba24bace839dda",
		res.Meta.DocHash)
}

func Test_ToStrings(t *testing.T) {
	res, err := Run(ctx,
		Content("Test input reader of type text"),
		TargetType(Text),
		Limit(3),
		NoStopWords(true),
	)
	assert.Nil(t, err)

	strs := ToStrings(res.Tags)
	assert.Len(t, strs, 3)
}

func Test_CustomHTML(t *testing.T) {
	ytPage, _ := ioutil.ReadFile("yt_page.html")
	ext := &customHTML{}
	defer stopServer(startServer(fmt.Sprintf(":%d", port), string(ytPage)))
	res, err := Run(ctx,
		Source(fmt.Sprintf("http://localhost:%d", port)),
		Limit(2),
		TargetType(HTML),
		NoStopWords(true),
		ExtraTagWeightsString("link:0"),
		Extensions([]extension.Extension{ext}),
	)
	assert.Nil(t, err)
	assert.Len(t, res.Extensions, 1)
	assert.Equal(t, "Next Level Reynolds - YouTube", res.Meta.DocTitle)
	assert.Equal(t, "Ryan Reynolds", ext.text)

	var found int
	res.ForEach(func(i int, tag *model.Tag) {
		if tag.Value == "ryan" || tag.Value == "reynolds" {
			found++
		}
	})
	assert.Equal(t, 2, found)
}

// startServer is a simple HTTP server that displays the passed headers in the html.
func startServer(addr string, pageHTML string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(res, pageHTML)
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

func stopServer(srv *http.Server) {
	// close the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

const (
	port      = 8655
	indexHTML = `<!doctype html>
<html>
<head>
  <title>Test</title>
</head>
<body>
  <div id="box1">
    <div id="box2">
      <p>There was a Boy whose name was Jim;</p>
	  <p>His Friends were very good to him.
	  <p>They gave him Tea, and Cakes, and Jam,</p>
	  <p>And slices of delicious Ham,</p>
	  <p>And Chocolate with pink inside,</p>
	  <p>And little Tricycles to ride,</p>
	  <p>Andread him Stories through and through,</p>
	  <p>And even took him to the Zoo—</p>
	  <p>But there it was the dreadful Fate</p>
	  <p>Befell him, which I now relate.</p>
    </div>
  </div>
  <div id="box3" style="display:none">
	<p class="line">Now this was Jim’s especial Foible,</p>
	<p class="line">He ran away when he was able,</p>
	<p class="line">And on this inauspicious day</p>
	<p class="line">He slipped his hand and ran away!</p>
	<p class="line">He hadn’t gone a yard when—Bang!</p>
	<p class="line">With open Jaws, a Lion sprang,</p>
	<p class="line">And hungrily began to eat</p>
	<p class="line">The Boy: beginning at his feet.</p>
  </div>
  <foo-stuff>
  	<bar-stuff>
	  <a href="https://www.zoomio.org">Zoom IO is here</a>
	</bar-stuff>
  </foo-stuff>
  <script>
  	setTimeout(function() {
		document.querySelector('#box3').style.display = '';
	}, 3000);
  </script>
</body>
</html>`
)

type customHTML struct {
	text string
}

func (ext *customHTML) Name() string {
	return "custom-html"
}

func (ext *customHTML) Version() string {
	return "v0.0.1"
}

func (ext *customHTML) Result() *extension.Result {
	return extension.NewResult(ext, map[string]interface{}{"text": ext.text}, nil)
}

func (ext *customHTML) ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *thtml.HTMLContents) (bool, error) {
	tag := token.Data
	var appended bool
	if ext.text == "" && tag == "link" {
		var itemprop, content string
		for _, v := range token.Attr {
			if v.Key == "itemprop" {
				itemprop = v.Val
			}
			if v.Key == "content" {
				content = v.Val
			}
		}
		if itemprop == "name" && content != "" {
			// collect YouTube channel name
			ext.text = content
			// make it count as a tag too
			// 1st check if line is there and append it if it is not
			if lineIdx >= cnts.Len() {
				cnts.Append(lineIdx, tag, []byte(content))
				appended = true
			}
			// 2nd weight the line higher to boost its tags
			cnts.Weigh(lineIdx, 6)
		}
	}
	return appended, nil
}

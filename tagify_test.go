package tagify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		[]Option{Source(fmt.Sprintf("http://localhost:%d", port)), TargetType(HTML),
			Limit(5), NoStopWords(true), ContentOnly(true)},
		[]string{"test", "boy", "cakes", "chocolate", "delicious"},
		"Test",
		"bdb03356c79b2b1d9c69f4528ee398bbafc4a572629b713dcf4992bd43fd650ecedb4355ddd08fe1da748ac2c4babff71e3c425724793f0d4e636037121e123e",
	},
	{
		"run with query",
		[]Option{Source(fmt.Sprintf("http://localhost:%d", port)), TargetType(HTML),
			Limit(5), NoStopWords(true), Query("#box3 p"), ContentOnly(true)},
		[]string{"bang", "began", "boy", "day", "eat"},
		"",
		"e5e0aef65e77e87a3e23a3f157357444910f94f5dccd5d0fe185da73cb72a8b7bff6ac80d71cfca1da27e9d1b7a3e810a348ceeee52c2e4b68393c8ba5d92cc4",
	},
}

func Test_Run_HTML(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))
	for _, tt := range runTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Run(context.TODO(), tt.in...)
			assert.Nil(t, err)
			assert.Equal(t, HTML, res.Meta.ContentType)
			assert.Equal(t, tt.expectTitle, res.Meta.DocTitle)
			assert.Equal(t, tt.expectHash, res.Meta.DocHash)
			assert.ElementsMatch(t, tt.expectTags, res.TagsStrings())
		})
	}
}

func Test_GetTagsFromString(t *testing.T) {
	res, err := Run(context.TODO(), Content("Test input reader of type text"), TargetType(Text), Limit(3), NoStopWords(true))
	assert.Nil(t, err)
	assert.Len(t, res.Tags, 3)
	assert.Equal(t,
		"7d95ed3e8436c978f3e7f19f1645f89091f9fdb0439c15547f0a6f82bc4a0babebd06ff6285d9dff8db77861edf2cc8e6919ea5613bec0f30dba24bace839dda",
		res.Meta.DocHash)
}

func Test_ToStrings(t *testing.T) {
	res, err := Run(context.TODO(), Content("Test input reader of type text"), TargetType(Text), Limit(3), NoStopWords(true))
	assert.Nil(t, err)

	strs := ToStrings(res.Tags)
	assert.Len(t, strs, 3)
}

// startServer is a simple HTTP server that displays the passed headers in the html.
func startServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(res, indexHTML)
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
	if err := srv.Shutdown(context.TODO()); err != nil {
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
  <script>
  	setTimeout(function() {
		document.querySelector('#box3').style.display = '';
	}, 3000);
  </script>
</body>
</html>`
)

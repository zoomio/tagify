package tagify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTags(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))
	tags, err := GetTags(context.TODO(), fmt.Sprintf("http://localhost:%d", port), HTML, 5, false, true)
	assert.Nil(t, err)
	assert.Len(t, tags, 5)
	assert.Equal(t, []string{"test", "boy", "him", "able", "andread"}, ToStrings(tags))
}

func Test_GetTagsWithQuery(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))
	tags, err := GetTagsWithQuery(context.TODO(), fmt.Sprintf("http://localhost:%d", port), "#box3 p", HTML, 5, false, true)
	assert.Nil(t, err)
	assert.Len(t, tags, 5)
	assert.Equal(t, []string{"able", "away", "began", "boy", "day"}, ToStrings(tags))
}

func Test_Run(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))
	tags, err := Run(context.TODO(), Source(fmt.Sprintf("http://localhost:%d", port)), Query("#box3 p"),
		TargetType(HTML), Limit(5), Verbose(false), NoStopWords(true))
	assert.Nil(t, err)
	assert.Len(t, tags, 5)
	assert.Equal(t, []string{"able", "away", "began", "boy", "day"}, ToStrings(tags))
}

func Test_GetTagsFromString(t *testing.T) {
	tags, err := GetTagsFromString("Test input reader of type text", Text, 3, false, true)
	assert.Nil(t, err)
	assert.Len(t, tags, 3)
}

func Test_ToStrings(t *testing.T) {
	tags, _ := GetTagsFromString("Test input reader of type text", Text, 3, false, true)
	strs := ToStrings(tags)
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

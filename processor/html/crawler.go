package html

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zoomio/inout"

	"github.com/zoomio/tagify/config"
)

const (
	crwlBadThreshold         = 1.0
	crwlBadThresholdInterval = 1500
)

type parseFunc func(io.Reader, *config.Config, []HTMLExt, *webCrawler) *HTMLContents

type parseOut struct {
	cnt *HTMLContents
	err error
}

type crwlStats struct {
	start   time.Time
	total   float64
	bad     float64
	stopped bool
}

type webCrawler struct {
	cfg *config.Config
	parseFunc
	dataCh  chan *parseOut
	stopCh  chan struct{}
	stats   *atomic.Value
	wg      *sync.WaitGroup
	links   *sync.Map
	docs    *sync.Map
	domain  string
	verbose bool
	exts    []HTMLExt
}

func newWebCrawler(parse parseFunc, exts []HTMLExt, source string, verbose bool) (*webCrawler, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var links sync.Map
	var docs sync.Map
	var av atomic.Value
	return &webCrawler{
		parseFunc: parse,
		dataCh:    make(chan *parseOut, 5),
		stopCh:    make(chan struct{}),
		stats:     &av,
		wg:        &wg,
		links:     &links,
		docs:      &docs,
		domain:    toDomain(u),
		verbose:   verbose,
		exts:      exts,
	}, nil
}

func (c *webCrawler) getStats() *crwlStats {
	return c.stats.Load().(*crwlStats)
}

func (c *webCrawler) setStats(sts *crwlStats) {
	c.stats.Store(sts)
}

func (c *webCrawler) run(r io.Reader) *HTMLContents {
	defer close(c.dataCh)

	c.setStats(&crwlStats{
		start: time.Now(),
	})

	result := c.parseFunc(r, c.cfg, c.exts, c)

	// waiter
	go func(stopCh chan struct{}, wg *sync.WaitGroup) {
		wg.Wait()
		select {
		case <-c.stopCh:
			return
		default:
			close(stopCh)
		}
	}(c.stopCh, c.wg)

	wgReceiver := sync.WaitGroup{}
	wgReceiver.Add(1)

	// the receiver
	go func(crwl *webCrawler) {
		defer wgReceiver.Done()

		for {
			// try to exit the receiver goroutine
			// as early as possible.
			select {
			case <-crwl.stopCh:
				return
			default:
			}

			// even if stopCh is closed, the first
			// branch in this select block might be
			// still not selected for some loops.
			select {
			case <-crwl.stopCh:
				return
			case value := <-crwl.dataCh:
				if value.err == nil {
					result.lines = append(result.lines, value.cnt.lines...)
				}
				crwl.wg.Done()
			}
		}
	}(c)

	wgReceiver.Wait()

	return result
}

func (c *webCrawler) crawl(href string) {
	c.wg.Add(1)
	go c.schedule(href)
}

func (c *webCrawler) schedule(href string) {
	// add domain if relative link
	if strings.HasPrefix(href, "/") {
		href = c.domain + href
	}

	// skip wrong address
	u, err := url.Parse(href)
	if err != nil {
		c.trySend(&parseOut{err: err})
		return
	}

	// skip if different domain
	if c.domain != toDomain(u) {
		c.trySend(&parseOut{err: fmt.Errorf("different domain: %s", u.Host)})
		return
	}

	src := u.String()

	// skip visited links
	if _, ok := c.links.Load(src); ok {
		c.trySend(&parseOut{err: fmt.Errorf("already visisted: %s", src)})
		return
	}
	c.links.Store(src, true)

	r, err := inout.NewInOut(context.TODO(), inout.Source(src), inout.Verbose(c.verbose))
	if err != nil {
		c.trySend(&parseOut{err: err})
		return
	}

	cnt := c.parseFunc(&r, c.cfg, c.exts, c)
	h := fmt.Sprintf("%x", cnt.hash())

	// skip visited docs
	if _, ok := c.docs.Load(h); ok {
		c.trySend(&parseOut{err: fmt.Errorf("already visisted: %s", src)})
		return
	}
	c.docs.Store(h, true)

	if c.verbose {
		fmt.Printf("visit: %s\n", src)
	}

	c.trySend(&parseOut{cnt: cnt})
}

func (c *webCrawler) trySend(out *parseOut) {
	// try to exit the sender goroutine
	// as early as possible.
	select {
	case <-c.stopCh:
		c.wg.Done()
		return
	default:
		sts := c.getStats()

		if sts.stopped {
			c.wg.Done()
			return
		}

		total := sts.total + 1
		bad := sts.bad
		if out.err != nil {
			bad++
		}

		// check bad vs good rate every "crawlerBadThresholdInterval" and stop
		// if it is above the the "crawlerBadThreshold".
		if time.Since(sts.start).Milliseconds() >= crwlBadThresholdInterval {
			r := bad / total
			if r >= crwlBadThreshold {
				c.setStats(&crwlStats{stopped: true})
				c.wg.Done()
				return
			}
			c.setStats(&crwlStats{
				start: time.Now(),
				total: 0.0,
				bad:   0.0,
			})
		} else {
			c.setStats(&crwlStats{
				start: sts.start,
				total: total,
				bad:   bad,
			})
		}

		c.dataCh <- out
	}
}

func toDomain(u *url.URL) string {
	var sb strings.Builder
	sb.WriteString(u.Scheme)
	sb.WriteString("://")
	sb.WriteString(u.Host)
	return sb.String()
}

package tasks

import (
	"sync"

	"github.com/ZakkBob/AskDave/gocommon/url"
)

type Tasks struct {
	Robots   taskSlice `json:"robots"` // Bit clunky
	Sitemaps taskSlice `json:"sitemaps"`
	Pages    taskSlice `json:"pages"`
}

type taskSlice struct {
	Mu    sync.Mutex `json:"-"`
	Slice []url.URL  `json:"slice"`
}

// Returns next url in slice, returns nil if slice is empty
func (t *taskSlice) Next() (url.URL, bool) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if len(t.Slice) == 0 {
		return url.URL{}, false
	}
	u := t.Slice[0]
	t.Slice = t.Slice[1:len(t.Slice)]
	return u, true
}

package tasks

import (
	"ZakkBob/AskDave/crawler/url"
	"sync"
)

type Tasks struct {
	Robots   taskSlice `json:"robots"` // Bit clunky
	Sitemaps taskSlice `json:"sitemaps"`
	Pages    taskSlice `json:"pages"`
}

type taskSlice struct {
	Mu    sync.Mutex `json:"-"`
	Slice []url.Url  `json:"slice"`
}

// Returns next url in slice, returns nil if slice is empty
func (t *taskSlice) Next() (url.Url, bool) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if len(t.Slice) == 0 {
		return url.Url{}, false
	}
	var u url.Url
	t.Slice[0].Copy(&u)
	t.Slice = t.Slice[1:len(t.Slice)]
	return u, true
}

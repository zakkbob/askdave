package orm_test

import (
	"testing"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewLink(t *testing.T) {
	u1, err := url.ParseAbs("www.test.com/TestSaveNewLink")
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	u2, err := url.ParseAbs("www.test.com/TestSaveNewLink2")
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	p := page.Page{
		Url:           u1,
		Title:         "",
		OgTitle:       "",
		OgDescription: "",
		OgSiteName:    "",
		Links:         []url.Url{},
		Hash:          hash.Hashs(""),
	}

	ormP1, err := orm.SaveNewPage(p)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	p.Url = u2

	ormP2, err := orm.SaveNewPage(p)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	_, err = orm.SaveNewLink(ormP1, ormP2)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	links, err := orm.LinkDstsBySrc(u1.String())
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	assert.Equal(t, []url.Url{u2}, links)
}

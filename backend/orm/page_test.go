package orm_test

import (
	"testing"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"

	"github.com/ZakkBob/AskDave/backend/orm"
)

func TestSaveNewPage(t *testing.T) {
	u, err := url.ParseAbs("https://www.zakkdev.com")

	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	p := page.Page{
		Url:           u,
		Title:         "title",
		OgTitle:       "og titlke",
		OgDescription: "descriptiuon",
		OgSiteName:    "site_name",
		Links:         []url.Url{},
		Hash:          hash.Hashs(""),
	}

	orm.Connect("")
	defer orm.Close()

	_, err = orm.SaveNewSite("https://www.zakkdev.com")
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	_, err = orm.SaveNewPage(p)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}
}

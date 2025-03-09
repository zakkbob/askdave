package main

import (
	"regexp"
	"testing"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidatorByUrl(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	allowed := []string{"allowed1", "allowed2"}
	disallowed := []string{"disallowed1", "disallowed2", "disallowed3"}

	rows := pgxmock.NewRows([]string{"allowed_patterns", "disallowed_patterns"}).
		AddRow(allowed, disallowed)

	mock.ExpectQuery(`SELECT`).
		WithArgs("https://mateishome.page").
		WillReturnRows(rows)

	v, err := ValidatorByUrl(mock, "https://mateishome.page")

	assert.Equal(t, v.AllowedStrings(), allowed, "allowed patterns should match")
	assert.Equal(t, v.DisallowedStrings(), disallowed, "disallowed patterns should match")

	if err != nil {
		t.Errorf("error was not expected while getting validator: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveResultsRobots(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	u, _ := url.ParseAbs("test.com")
	u2, _ := url.ParseAbs("test2.com")
	u3, _ := url.ParseAbs("test3.com")

	allowedStrings := []string{
		"123",
		"123",
		"123",
	}

	disallowedStrings := []string{
		"123",
		"123",
		"123",
	}

	validator, _ := robots.FromStrings(allowedStrings, disallowedStrings)

	results := tasks.Results{
		Robots: map[string]*tasks.RobotsResult{
			u.String(): {
				Url:           &u,
				Success:       true,
				FailureReason: tasks.NoFailure,
				Hash:          hash.Hashs(""),
				Changed:       true,
				Validator:     validator,
			}, u2.String(): {
				Url:           &u2,
				Success:       false,
				FailureReason: tasks.NoFailure,
				Hash:          hash.Hashs(""),
				Changed:       true,
				Validator:     validator,
			}, u3.String(): {
				Url:           &u3,
				Success:       true,
				FailureReason: tasks.NoFailure,
				Hash:          hash.Hashs(""),
				Changed:       false,
				Validator:     validator,
			}},
		Pages: map[string]*tasks.PageResult{},
	}

	mock.ExpectExec("UPDATE").
		WithArgs(allowedStrings, disallowedStrings, u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	err = SaveResults(mock, &results)
	if err != nil {
		t.Errorf("error was not expected while saving results: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveResultsPageSuccessNoChange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	u, _ := url.ParseAbs("test.com")

	results := tasks.Results{
		Robots: map[string]*tasks.RobotsResult{},
		Pages: map[string]*tasks.PageResult{
			u.String(): {
				Url:           &u,
				Success:       true,
				FailureReason: tasks.NoFailure,
				Changed:       false,
				Page:          nil,
			},
		},
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO crawl (page_id, datetime, success, content_changed) SELECT page.id, CURRENT_TIMESTAMP, TRUE, FALSE FROM page WHERE url = $1;")).
		WithArgs(u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	mock.ExpectExec("UPDATE page").
		WithArgs(30, 1, u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	err = SaveResults(mock, &results)
	if err != nil {
		t.Errorf("error was not expected while saving results: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveResultsPageSuccessAndChange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	u, _ := url.ParseAbs("test.com")
	u2, _ := url.ParseAbs("test2.com")
	u3, _ := url.ParseAbs("test3.com")

	page := page.Page{
		Url:           u,
		Title:         "Title",
		OgTitle:       "OgTitle",
		OgDescription: "OgDescription",
		OgSiteName:    "OgSiteName",
		Links: []url.Url{
			u2, u3,
		},
		Hash: hash.Hashs("testhash"),
	}

	results := tasks.Results{
		Robots: map[string]*tasks.RobotsResult{},
		Pages: map[string]*tasks.PageResult{
			u.String(): {
				Url:           &u,
				Success:       true,
				FailureReason: tasks.NoFailure,
				Changed:       true,
				Page:          &page,
			},
		},
	}

	// Page success and change
	mock.ExpectExec("INSERT INTO crawl").
		WithArgs(page.Title, page.OgTitle, page.OgDescription, page.Hash, u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	mock.ExpectExec("DELETE FROM link").
		WithArgs(u.String()).
		WillReturnResult(pgxmock.NewResult("", 3))

	mock.ExpectExec("INSERT INTO link").
		WithArgs(u.String(), u2.String(), 1).
		WillReturnResult(pgxmock.NewResult("", 1))

	mock.ExpectExec("INSERT INTO link").
		WithArgs(u.String(), u3.String(), 1).
		WillReturnResult(pgxmock.NewResult("", 1))

	mock.ExpectExec("UPDATE page").
		WithArgs(30, 1, u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	err = SaveResults(mock, &results)
	if err != nil {
		t.Errorf("error was not expected while saving results: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveResultsPageNoSuccess(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	u, _ := url.ParseAbs("test.com")

	page := page.Page{
		Url:           u,
		Title:         "Title",
		OgTitle:       "OgTitle",
		OgDescription: "OgDescription",
		OgSiteName:    "OgSiteName",
		Links:         []url.Url{},
		Hash:          hash.Hashs("testhash"),
	}

	results := tasks.Results{
		Robots: map[string]*tasks.RobotsResult{},
		Pages: map[string]*tasks.PageResult{
			u.String(): {
				Url:           &u,
				Success:       false,
				FailureReason: tasks.NoFailure,
				Changed:       false,
				Page:          &page,
			},
		},
	}

	// Page success and change
	mock.ExpectExec("INSERT INTO crawl").
		WithArgs(u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	mock.ExpectExec("UPDATE page").
		WithArgs(30, 1, u.String()).
		WillReturnResult(pgxmock.NewResult("", 1))

	err = SaveResults(mock, &results)
	if err != nil {
		t.Errorf("error was not expected while saving results: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

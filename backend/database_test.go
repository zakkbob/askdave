package main

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestShouldGetValidatorByUrl(t *testing.T) {
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

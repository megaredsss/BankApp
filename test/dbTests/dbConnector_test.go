package dbTests

import (
	"BankApp/db"
	"testing"
)

func test(t *testing.T) {
	got := db.GetDB()
	if got != nil {
		t.Errorf("no error")
	}

}

package lib_test

import (
	"testing"

	"github.com/caffeines/sharehub/lib"
)

func TestHashPassword(t *testing.T) {
	t.Run("it should pass test for valid password and hash", func(t *testing.T) {
		password := "password"
		hash, _ := lib.HashPassword(password)
		comp := lib.CheckPasswordHash(password, hash)
		if comp != true {
			t.Errorf("expected 'true' but got '%v'", comp)
		}
	})
	t.Run("it should pass test for invalvalid password and hash", func(t *testing.T) {
		password := "password"
		fakePassword := "fakePassword"
		hash, _ := lib.HashPassword(password)
		comp := lib.CheckPasswordHash(fakePassword, hash)
		if comp != false {
			t.Errorf("expected 'false' but got '%v'", comp)
		}
	})
}

// +build !integration

package util

import (
	"testing"
)

// TestUUIDUniqueness - ensure two calls to UUID() generate unique values
func TestUUIDUniqueness(t *testing.T) {
  s1 := UUID()
  s2 := UUID()

  if s1 == s2 {
    t.Errorf("calling UUID() twice yields the same result: %v and %v, expected 2 different results", s1, s2)
  }
}

// TestRightPassword - ensure comparing the hash and an incorrect password gives an error
func TestWrongPassword(t *testing.T) {
  hash, err := HashAndSalt("terrapukka")

  err = CompareHashAndText(string(hash), "not_terrapukka")

  if err == nil {
    t.Errorf("comparing a hash (%v) and an incorrect (not_terrapukka) password yields no error, expected an error", string(hash))
  }
}

// TestRightPassword - ensure comparing the hash and the correct password gives no error
func TestRightPassword(t *testing.T) {
  hash, err := HashAndSalt("terrapukka")

  err = CompareHashAndText(string(hash), "terrapukka")

  if err != nil {
    t.Errorf("comparing a hash and the correct password yields an error: %v, expected no error", err.Error())
  }
}

// TestSalting - ensure the same password hashed and salted twice gives different results
func TestSalting(t *testing.T) {
  hash1, _ := HashAndSalt("terrapukka")
  hash2, _ := HashAndSalt("terrapukka")

  if hash1 == hash2 {
    t.Errorf("Hashing the same password twice yields the same result: %v and %v, expected 2 different results", hash1, hash2)
  }
}

// +build !integration

package db

import (
  "testing"
)

// TestAuthorizeEmailPasswordTrue - ensure email and password validation fails for incorrect inputs
func TestValidateEmailPasswordFalse(t *testing.T) {
  nopass := User {
    Email: "pukka@terraling.com",
  }

  err := validateEmailPassword(nopass)

  if err == nil {
    t.Errorf("validating a user with no password yields no error, expected an error")
  }

  noemail := User {
    Password: "password",
  }

  err = validateEmailPassword(noemail)

  if err == nil {
    t.Errorf("validating a user with no email yields no error, expected an error")
  }

  invalidemail := User {
    Email: "pukka@terraling",
    Password: "password",
  }

  err = validateEmailPassword(invalidemail)

  if err == nil {
    t.Errorf("validating a user with an invalid email yields no error, expected an error")
  }

  shortpass := User {
    Email: "pukka@terraling.com",
    Password: "pass",
  }

  err = validateEmailPassword(shortpass)

  if err == nil {
    t.Errorf("validating a user with a password less than 8 characters yields no error, expected an error")
  }
}

// TestAuthorizeEmailPasswordTrue - ensure email and password validation succeeds for correct inputs
func TestValidateEmailPasswordTrue(t *testing.T) {
  user := User {
    Email: "pukka@terraling.com",
    Password: "password",
  }

  err := validateEmailPassword(user)

  if err != nil {
    t.Errorf("validating a valid user yields an error: %v, expected no error", err.Error())
  }
}

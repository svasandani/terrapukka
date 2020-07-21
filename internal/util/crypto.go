package util

import (
  "fmt"
  "math/rand"
  "time"

  "golang.org/x/crypto/bcrypt"
)

var randInt int64

// HashAndSalt - hash a specified text with salt and return the salted hash as a string
func HashAndSalt(text string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost);

  return string(hash), err
}

// UUID - generate a unique UUID
func UUID() (string) {
  randInt = rand.Int63n(time.Now().UnixNano())
  rand.Seed(time.Now().UnixNano() + randInt)

  b := make([]byte, 16)
  _, err := rand.Read(b)

  CheckError("Error reading byte slice:", err)

  if err != nil {
    return ""
  }

  uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

  return uuid
}

// CompareHashAndText - compare hash and text passed in via bcrypt
func CompareHashAndText(hash string, text string) (error) {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}

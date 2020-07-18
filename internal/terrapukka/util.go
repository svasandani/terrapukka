package api

import (
  "log"
)

/**************************************************************

UTILITY

**************************************************************/

// Check for error; if not nil, print a message along with the error.
func checkError(msg string, err error) {
  if err != nil {
    log.Println(msg)
    log.Println(err.Error())
  }
}

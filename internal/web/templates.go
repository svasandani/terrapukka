package web

import (
  "html/template"

  "github.com/svasandani/terrapukka/internal/util"
)

var tmpl *template.Template
var err error

// Init - initialize common templates
func Init() {
  tmpl, err = template.ParseGlob("./web/templates/*.html")

  util.CheckError("Error parsing templates:", err)
}

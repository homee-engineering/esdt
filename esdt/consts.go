package esdt

import "regexp"

const DefaultConnUrl = "http://localhost:9200"
const DefaultTargetDir = "es/operations"

var JsonRegEx = regexp.MustCompile(".+\\.json")

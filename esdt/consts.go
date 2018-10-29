package esdt

import "regexp"

const DefaultConnUrl = "http://localhost:9200"
const DefaultTargetDir = "es/operations"
const DefaultConfigFile = "es/config.yml"

var JsonRegEx = regexp.MustCompile(".+\\.json")

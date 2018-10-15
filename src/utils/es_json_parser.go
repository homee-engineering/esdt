package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var JsonRegEx = regexp.MustCompile(".+\\.json")

type DataTemplate struct {
	Method   string                 `json:"method"`
	Uri      string                 `json:"uri"`
	Body     map[string]interface{} `json:"body"`
	Rollback RollbackTemplate       `json:"rollback"`
	Id       string
}

type RollbackTemplate struct {
	Method string                 `json:"method"`
	Uri    string                 `json:"uri"`
	Body   map[string]interface{} `json:"body"`
}

func RunFiles(config *config, dir string) error {
	fi, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not find directory %s", dir))
	}

	var dataTemplates []*DataTemplate
	for _, v := range fi {
		if JsonRegEx.MatchString(v.Name()) {
			dataTemplate, err := LoadDataTemplate(dir, v.Name())
			if err != nil {
				return err
			}
			dataTemplates = append(dataTemplates, dataTemplate)
		}
	}

	executeDataTemplates(config.Conn, dataTemplates)

	return nil
}

func LoadDataTemplate(dir string, dataTemplateFilename string) (*DataTemplate, error) {
	fp := filepath.Join(dir, dataTemplateFilename)
	out, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problems reading file %s", fp))
	}
	var dataTemplate DataTemplate
	err = json.Unmarshal(out, &dataTemplate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse file %s, double check your json", fp))
	}
	dataTemplate.Id = strings.TrimSuffix(dataTemplateFilename, filepath.Ext(dataTemplateFilename))
	return &dataTemplate, nil
}

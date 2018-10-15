package utils

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"time"
)

type operations struct {
	InsertedAt time.Time `json:"inserted_at"`
}

type documentExistsRes struct {
	Found bool `json:"found"`
}

type documentDeletedRes struct {
	Result string `json:"result"`
}

const NoRollbackFieldErrorMsg = "No rollback listed"

func CreateOperationsIndex(conn string) error {
	body := "{ \"mappings\": { \"_doc\": { \"properties\": { \"inserted_at\": { \"type\": \"date\" } } } } }"
	return runEsQueryAndValidate(conn, "operations", "put", body)
}

func OperationsIndexExists(conn string) (bool, error) {
	res, err := runEsQuery(conn, "operations", "head", nil)

	if err != nil {
		return false, err
	}

	return res.Response().StatusCode > 199 && res.Response().StatusCode < 300, nil
}

func DeleteOperationIndex(conn string, rollbackId string) error {
	res, err := runEsQuery(conn, "operations/_doc/"+rollbackId, "delete", nil)

	if err != nil {
		return err
	}

	var d documentDeletedRes
	json.NewDecoder(res.Response().Body).Decode(&d)

	bodyBytes, _ := ioutil.ReadAll(res.Response().Body)

	if d.Result != "deleted" {
		return errors.New(fmt.Sprintf("Failed to delete %s from the operations records. Got %s", rollbackId, string(bodyBytes)))
	}

	return nil
}

func RollbackDataTemplate(conn string, dt *DataTemplate) error {
	if dt.Rollback.Uri == "" || dt.Rollback.Method == "" {
		return errors.New(NoRollbackFieldErrorMsg)
	} else {
		return runEsQueryAndValidate(conn, dt.Rollback.Uri, dt.Rollback.Method, dt.Rollback.Body)
	}
}

func operationsDocumentExists(conn string, id string) bool {
	res, err := runEsQuery(conn, "operations/_doc/"+id, "get", nil)

	if err != nil {
		return false
	}

	var d documentExistsRes
	json.NewDecoder(res.Response().Body).Decode(&d)

	return d.Found
}

func executeDataTemplates(conn string, dataTemplates []*DataTemplate) ([]*DataTemplate, []error) {
	failed := make([]*DataTemplate, 0)
	errs := make([]error, 0)

	for _, v := range dataTemplates {
		err := executeDataTemplate(conn, v)
		if err != nil {
			failed = append(failed, v)
			errs = append(errs, err)
			if strings.Contains(err.Error(), "Already ran") {
				color.Yellow("%s has already run", v.Id)
			} else {
				color.Red("%s failed to run: %s", v.Id, err.Error())
			}
		} else {
			color.Green("%s ran successfully", v.Id)
		}
	}

	if len(failed) == 0 {
		failed = nil
	}

	return failed, errs
}

func executeDataTemplate(conn string, dataTemplate *DataTemplate) error {
	if !operationsDocumentExists(conn, dataTemplate.Id) {
		err := runEsQueryAndValidate(conn, dataTemplate.Uri, dataTemplate.Method, dataTemplate.Body)

		operations := operations{
			InsertedAt: time.Now(),
		}
		if err != nil {
			newError := errors.New(err.Error())
			err = RollbackDataTemplate(conn, dataTemplate)
			if err != nil {
				newError = errors.Wrap(newError, "Rollback failed")
			}
			return newError
		} else {
			err = runEsQueryAndValidate(conn, "/operations/_doc/"+dataTemplate.Id, "post", &operations)
			if err != nil {
				return errors.New("Failed to add data template to operations")
			}
		}
	} else {
		return errors.New("Already ran")
	}
	return nil
}

func runEsQuery(conn string, uri string, method string, bodyJson interface{}) (*req.Resp, error) {
	r := req.New()
	u, err := url.Parse(conn)
	if err != nil {
		return nil, errors.New("Invalid connection URL")
	}
	u.Path = path.Join(u.Path, uri)
	esUrl := u.String()
	body := req.BodyJSON(bodyJson)

	var res *req.Resp

	switch strings.ToLower(method) {
	case "get":
		res, err = r.Get(esUrl, body)
	case "post":
		res, err = r.Post(esUrl, body)
	case "put":
		res, err = r.Put(esUrl, body)
	case "head":
		res, err = r.Head(esUrl, body)
	case "delete":
		res, err = r.Delete(esUrl, body)
	default:
		return nil, errors.New("Invalid HTTP method")
	}

	return res, nil
}

func runEsQueryAndValidate(conn string, uri string, method string, bodyJson interface{}) error {
	res, err := runEsQuery(conn, uri, method, bodyJson)

	if err != nil {
		return err
	}

	if res.Response().StatusCode < 200 || res.Response().StatusCode > 299 {
		bodyBytes, _ := ioutil.ReadAll(res.Response().Body)
		return errors.New(fmt.Sprintf("status code was not 200: %s. Reason: %s", res.Response().Status, string(bodyBytes)))
	}

	return nil
}

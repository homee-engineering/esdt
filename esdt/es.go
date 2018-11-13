package esdt

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

const NoRollbackFieldErrorMsg = "No rollback listed"

type documentDeletedRes struct {
	Result string `json:"result"`
}

func (e *esdtImpl) deleteOperationIndex(rollbackId string) error {
	res, err := e.runEsQuery("operations/_doc/"+rollbackId, "delete", nil)

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

func (e *esdtImpl) createOperationsIndex() error {
	body := "{ \"mappings\": { \"_doc\": { \"properties\": { \"inserted_at\": { \"type\": \"date\" } } } } }"
	return e.runEsQueryAndValidate("operations", "put", body)
}

func (e *esdtImpl) operationsIndexExists() (bool, error) {
	res, err := e.runEsQuery("operations", "head", nil)

	if res == nil {
		return false, errors.New("no response received from Elasticsearch")
	}

	if err != nil {
		return false, err
	}

	return res.Response().StatusCode > 199 && res.Response().StatusCode < 300, nil
}

func (e *esdtImpl) rollbackDataTemplate(dt *Operation) error {
	if dt.Rollback.Uri == "" || dt.Rollback.Method == "" {
		return errors.New(NoRollbackFieldErrorMsg)
	} else {
		err := e.runEsQueryAndValidate(dt.Rollback.Uri, dt.Rollback.Method, dt.Rollback.Body)
		if err != nil {
			return err
		}
		return e.deleteOperationIndex(dt.Id)
	}
}

func (e *esdtImpl) operationsDocumentExists(id string) bool {
	res, err := e.runEsQuery("operations/_doc/"+id, "get", nil)

	if res == nil || err != nil {
		return false
	}

	var d documentExistsRes
	json.NewDecoder(res.Response().Body).Decode(&d)

	return d.Found
}

func (e *esdtImpl) executeDataTemplates(dataTemplates []*Operation) ([]*Operation, []error) {
	failed := make([]*Operation, 0)
	errs := make([]error, 0)

	for _, v := range dataTemplates {
		err := e.executeDataTemplate(v)
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

func (e *esdtImpl) executeDataTemplate(operation *Operation) error {
	if !e.operationsDocumentExists(operation.Id) {
		err := e.runEsQueryAndValidate(operation.Uri, operation.Method, operation.Body)

		operations := operations{
			InsertedAt: time.Now(),
		}
		if err != nil {
			newError := errors.Wrap(err, "")
			err = e.rollbackDataTemplate(operation)
			if err != nil {
				newError = errors.Wrap(newError, "RollbackFile failed")
			}
			return newError
		} else {
			err = e.runEsQueryAndValidate("/operations/_doc/"+operation.Id, "post", &operations)
			if err != nil {
				return errors.New("Failed to add data template to operations")
			}
		}
	} else {
		return errors.New("Already ran")
	}
	return nil
}

func (e *esdtImpl) runEsQuery(uri string, method string, bodyJson interface{}) (*req.Resp, error) {
	r := req.New()
	u, err := url.Parse(e.getConn())
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

func (e *esdtImpl) runEsQueryAndValidate(uri string, method string, bodyJson interface{}) error {
	res, err := e.runEsQuery(uri, method, bodyJson)

	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("did not receive a response from elasticsearch")
	}

	if res.Response().StatusCode < 200 || res.Response().StatusCode > 299 {
		bodyBytes, _ := ioutil.ReadAll(res.Response().Body)
		errors.New(fmt.Sprintf("status code was not 200: %s. Reason: %s", res.Response().Status, string(bodyBytes)))
		return errors.New(fmt.Sprintf("status code was not 200: %s. Reason: %s", res.Response().Status, string(bodyBytes)))
	}

	return nil
}

func (e *esdtImpl) getConn() (url string) {
	strs := strings.Split(e.Config.Conn, "://")
	if len(strs) == 1 {
		url = fmt.Sprintf("%s:%s@%s", e.Config.Username, e.Config.Password, strs[0])
	} else {
		url = fmt.Sprintf("%s://%s:%s@%s", strs[0], e.Config.Username, e.Config.Password, strings.Join(strs[1:], ""))
	}
	return
}

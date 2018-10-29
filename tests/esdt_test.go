package tests

import (
	"context"
	"encoding/json"
	"esdt/esdt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EsdtTestSuite struct {
	ElasticsearchTestSuite
}

var op = &esdt.Operation{
	Id:     "some_operation",
	Method: "PUT",
	Uri:    "test",
	Rollback: esdt.RollbackTemplate{
		Uri:    "test",
		Method: "DELETE",
	},
}

func (ets *EsdtTestSuite) TestRun() {
	e := esdt.New(&esdt.Config{
		Conn: ets.url,
	})

	err := e.Run(op)

	ets.NoError(err)

	ex, err := ets.client.IndexExists("test").Do(context.Background())
	ets.Nil(err)
	ets.True(ex)

	ex, err = ets.client.IndexExists("operations").Do(context.Background())
	ets.Nil(err)
	ets.True(ex)

	gr, err := ets.client.Get().Id("some_operation").Index("operations").Do(context.Background())
	ets.Nil(err)
	source := make(map[string]interface{})
	if err != nil {
		ets.FailNow(err.Error())
	}
	json.Unmarshal(*gr.Source, &source)
	ets.Equal("some_operation", gr.Id)
	ets.NotNil(source["inserted_at"])
}

func (ets *EsdtTestSuite) TestRollback() {
	e := esdt.New(&esdt.Config{
		Conn: ets.url,
	})

	operation := &esdt.Operation{
		Id:     "some_operation_1",
		Method: "PUT",
		Uri:    "test1",
		Rollback: esdt.RollbackTemplate{
			Uri:    "test1",
			Method: "DELETE",
		},
	}

	err := e.Run(operation)
	ets.Nil(err)

	err = e.Rollback(operation)
	ets.Nil(err)

	ex, err := ets.client.IndexExists(operation.Uri).Do(context.Background())
	ets.Nil(err)
	ets.False(ex)

	ex, err = ets.client.IndexExists("operations").Do(context.Background())
	ets.Nil(err)
	ets.True(ex)

	_, err = ets.client.Get().Id(operation.Id).Index("operations").Do(context.Background())
	ets.EqualError(err, "elastic: Error 404 (Not Found)")
}

func TestName(t *testing.T) {
	suite.Run(t, new(EsdtTestSuite))
}

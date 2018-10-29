package tests

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"log"
)

func createEsDb() (pool *dockertest.Pool, resource *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err = pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "elasticsearch",
			Tag:        "6.4.1",
		})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		sniff := false

		url := fmt.Sprintf("http://localhost:%s", resource.GetPort("9200/tcp"))

		elasticClient, err := elastic.NewClientFromConfig(&config.Config{
			URL:   url,
			Sniff: &sniff,
		})

		if err != nil {
			return err
		}

		_, _, err = elasticClient.Ping(url).Do(context.TODO())

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return
}

type ElasticsearchTestSuite struct {
	suite.Suite
	pool     *dockertest.Pool
	resource *dockertest.Resource
	client   *elastic.Client
	url      string
}

func (suite *ElasticsearchTestSuite) SetupSuite() {
	suite.pool, suite.resource = createEsDb()
	sniff := false
	suite.url = fmt.Sprintf("http://localhost:%s", suite.resource.GetPort("9200/tcp"))

	client, err := elastic.NewClientFromConfig(&config.Config{
		URL:   suite.url,
		Sniff: &sniff,
	})

	if err != nil {
		panic(err)
	}

	if _, _, err := client.Ping(suite.url).Do(context.Background()); err != nil {
		panic(err)
	}

	suite.client = client
}

func (suite *ElasticsearchTestSuite) TearDownSuite() {
	err := suite.pool.Purge(suite.resource)
	if err != nil {
		log.Fatalf("Could not close docker resource: %s", err)
	}
}

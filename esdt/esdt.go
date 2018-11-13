// esdt is a tool for seeding/migrating your Elasticsearch datastore. It can be used both as a CLI or a
// library.
package esdt

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// A tool for seeding/migrating your Elasticsearch datastore. It can be used both as a CLI or a
// library
type Esdt interface {
	// Runs all operations in the TargetDir. If any of the operations have been run previously,
	// it is skipped.
	//
	// If an error occurs while running the Operations, it is skipped, and
	// will be reattempted the next call to Run. The function will not
	// return an error if an Operation fails. Only if the error prevents
	// subsequent Operations from running
	RunAll() error

	// Runs a specified Operation. If the operation has been run previously, no action is taken.
	//
	// If the operations index has not yet been created on the Elasticsearch, it is created here.
	// This command also attempts to perform a rollback if an error occurs. There is no need to
	// call Rollback(Operation) if an error is returned.
	Run(operation *Operation) error

	// Same as Rollback but combines the steps of Load and Rollback
	//
	// The filename passed in must also contain the file extension (*.json)
	RollbackFile(filename string) error

	// Attempts to rollback any previously run Operation. If the operation
	// has not yet been run, an error is returned
	Rollback(operation *Operation) error

	// Load an operation from the TargetDir into an Operation struct.
	//
	// The filename passed in must also contain the file extension (*.json)
	Load(filename string) (*Operation, error)

	// Get the passed in Config struct that is tied to this instance
	// of esdt
	GetConfig() *Config
}

type esdtImpl struct {
	Config *Config
}

// A singular piece of instruction to be run against the Elasticsearch cluster
// No operation with the same Id will be run more than once against the cluster
type Operation struct {
	// The HTTP method for the Elasticsearch call.
	//
	// Valid methods are GET, PUT, POST, DELETE, HEAD, PATCH
	Method string `json:"method"`

	// The URI for the resource you're targeting
	//
	// e.g. users/_search
	Uri string `json:"uri"`

	// The body of the Elasticsearch request. Not required.
	Body map[string]interface{} `json:"body"`

	// The work that will be done if Rollback is called on this Operation
	Rollback RollbackTemplate `json:"rollback"`

	// The Id for the Operation. This identifies whether or not the Operation
	// has run previously. If two Operations have the same Id, only one can run.
	Id string
}

// Same as an Operation but is only run when Rollback is called on the operation
type RollbackTemplate struct {
	Method string                 `json:"method"`
	Uri    string                 `json:"uri"`
	Body   map[string]interface{} `json:"body"`
}

// The configuration used for all calls on the esdt struct
type Config struct {
	// The Elasticsearch URL e.g. http://my-elasticsearch.com or http://localhost:9200
	Conn string

	// The directory housing all of the files relevant to the esdt. By default,
	// this is es/operations within your current directory
	TargetDir string

	// The fullfile path location of your config file if it is not in the default es/config.yml
	ConfigFile string

	// The environment used for this esdt. This is the key used in your ConfigFile. Defaults to dev.
	Env string

	// The username used for the Elasticsearch cluster
	Username string

	// The password used for the Elasticsearch cluster
	Password string
}

func (e *esdtImpl) GetConfig() *Config {
	return e.Config
}

func (e *esdtImpl) RunAll() error {
	ex, err := e.operationsIndexExists()
	if err != nil {
		return err
	}

	if !ex {
		err = e.createOperationsIndex()
		if err != nil {
			return err
		}
	}

	fi, err := ioutil.ReadDir(e.Config.TargetDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not find directory %s", e.Config.TargetDir))
	}

	var operations []*Operation
	for _, v := range fi {
		operation, err := e.Load(v.Name())
		if err == nil {
			operations = append(operations, operation)
		}
	}

	e.executeDataTemplates(operations)

	return nil
}

func (e *esdtImpl) RollbackFile(filename string) error {
	dt, err := e.Load(filename)
	if err != nil {
		return err
	}
	return e.rollbackDataTemplate(dt)
}

func (e *esdtImpl) Load(filename string) (*Operation, error) {
	if !JsonRegEx.MatchString(filename) {
		return nil, errors.New(fmt.Sprintf("invalid elasticsearch operation %s", filename))
	}
	fp := filepath.Join(e.Config.TargetDir, filename)
	out, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problems reading file %s", fp))
	}
	var dataTemplate Operation
	err = json.Unmarshal(out, &dataTemplate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse file %s, double check your json", fp))
	}
	dataTemplate.Id = strings.TrimSuffix(filename, filepath.Ext(filename))
	return &dataTemplate, nil
}

// Attempts to rollback any previously run Operation. If the operation
// has not yet been run, an error is returned
func (e *esdtImpl) Rollback(operation *Operation) error {
	return e.rollbackDataTemplate(operation)
}

// Runs a specified Operation.
//
// If the operations index has not yet been created on the Elasticsearch, it is created here.
// This command also attempts to perform a rollback if an error occurs. There is no need to
// call Rollback(Operation) if an error is returned.
func (e *esdtImpl) Run(operation *Operation) error {
	operation.Id = strings.TrimSpace(operation.Id)
	if operation.Body == nil {
		operation.Body = make(map[string]interface{})
	}
	if operation.Rollback.Body == nil {
		operation.Rollback.Body = make(map[string]interface{})
	}
	ex, err := e.operationsIndexExists()
	if err != nil {
		return err
	}

	if !ex {
		err = e.createOperationsIndex()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("Could not find directory %s", e.Config.TargetDir))
	}

	return e.executeDataTemplate(operation)
}

// Create a new esdt object which can run operations against an Elasticsearch cluster.
//
// Config
//
// The config used on esdt will take precedence over an config found in the config.yml in
// the operations directory
func New(config *Config) Esdt {
	return &esdtImpl{
		Config: loadConfig(config),
	}
}

type yamlConfig map[string]*Config

func loadConfig(in *Config) (c *Config) {
	if in == nil {
		in = &Config{}
	}

	if in.ConfigFile == "" {
		in.ConfigFile = DefaultConfigFile
	}

	content, err := ioutil.ReadFile(in.ConfigFile)
	if err == nil {
		var yc yamlConfig
		err = yaml.Unmarshal(content, &yc)
		if err != nil {
			return nil
		}

		c = yc[in.Env]
		if c == nil {
			c = in
		}
	}
	c = defaultConfig(c)

	mergo.Merge(c, *in, mergo.WithOverride)

	return c
}

func defaultConfig(c *Config) *Config {
	if c == nil {
		c = &Config{}
	}
	if c.TargetDir == "" {
		c.TargetDir = DefaultTargetDir
	}
	if c.Conn == "" {
		c.Conn = DefaultConnUrl
	}
	return c
}

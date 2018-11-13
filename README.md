# Elasticsearch Data Tool
esdt is a tool for seeding/migrating your Elasticsearch datastore. It can be used both as a CLI or a 
library

## Library installation
```bash
go get github.com/homee-engineering/esdt
```

## CLI Installation
1. Go to the [releases](https://github.com/homee-engineering/esdt/releases)
1. Select your platform and download
1. Move the binary `mv ~/Downloads/*.esdt /usr/local/bin/esdt`:
1. Run `chmod +x /usr/local/bin/esdt`

## CLI
### Usage
Generate the `es/operations` directory in your current directory
```bash
esdt gen dir
```
Generate an `operation` - a json file containing the Elasticsearch query:
```bash
esdt gen temp -m put -uri my_index create_my_index
``` 
A json file will appear in the `es/operations` folder with a name of `<timestamp>_create_my_index.json`
that looks like
```json
{
  "method": "PUT",
  "uri": "my_index",
  "body": {

  },
  "rollback": {
    "method": "DELETE",
    "uri": "",
    "body": {

    }
  }
}
```
* The `method` field is the HTTP method
* `uri` is the resource to target within Elasticsearch
* `body` is the body of the Elasticsearch request
* `rollback` is run during the `esdt rollback` command. This query should undo the above operation

To run all of the `*.json` operations against your Elasticsearch store simply run
```bash
esdt run
``` 
This will run the operations against the Elasticsearch store located at `http://localhost:9200` by default

To undo the index creation, add `my_index` to the `rollback.uri` field and run
```bash
esdt rollback <timestamp>_create_my_index
```

### Config

All global flags can be configured via command line flag, environment variable, or `config.yml` in your target
directory (by default `es/operations`)

#### Global flags
| Flag       | Env Var             | Config.yml field | Description                                                                                    |
|------------|---------------------|------------------|------------------------------------------------------------------------------------------------|
| `conn`     | `ELASTICSEARCH_URL` | `conn`           | The Elasticsearch base URL to run all operations against. Default is `http://localhost:9200`   |
| `dir`      | `ESDT_TARGET_DIR`   | `dir`            | The directory of the data operations. Default is `es/operations`                               |
| `config`   | N/A                 | N/A              | The location of your config YAML. Default is ./es/config.tml                                   |
| `env`      | `ESDT_ENV`          | <env\>           | The environment to run the tool in. Default is dev                                             |
| `username` | `ESDT_USER`         | `user`           | The username for the Elasticsearch cluster. Default is ""                                      |
| `password` | `ESDT_PASSWORD`     | `pw`             | The password for the Elasticsearch cluster. Default is ""                                      |

#### Config.yml
The default config file looks like
```yaml
dev:
  conn: http://localhost:9200
  dir: es/operations
prod:
  conn: http://elasticsearch:9200
  dir: es/operations
```
The top level fields are the `env` global flag which defaults to `dev`. In order to use the `prod` config simply
run
```bash
esdt -env prod run
```

### Precedence
1. Command line options
1. Environment variables
1. `config.yml`

## Library
### Basic Usage
As stated previously, the pieces of work that the `esdt` deals with are called `operations`. Here's an
example of how to load, run, and rollback an operation
```go
package main

import (
   "github.com/homee-engineering/esdt/esdt"
   "os"
)

func main() {
    // Create a new esdt struct
    e := esdt.New(&esdt.Config{
       Env:        os.Getenv("ENV"),
       ConfigFile: "somewhere/else/weird-config.yml",
       Conn:       "http:my-elasticsearch.com",
       TargetDir:  "something/different",
    })
    
    // Load an operation from your something/different/* directory
    operation, err := e.Load("20181025164223_something.json") // Must be the full filename of the operation including the extension
    if err != nil {
       panic(err)
    }
    
    // Run the operation
    err := e.Run(operation)
    if err != nil {
       panic(err)
    }
    
    // Rollback the operation
    err := e.Rollback(operation)
    if err != nil {
        panic(err)
    }
}
```
To run all of the `operations`
```go
package main

import (
   "github.com/homee-engineering/esdt/esdt"
   "os"
)

func main() {
    // Create a new esdt struct
    e := esdt.New(&esdt.Config{
       Env:        os.Getenv("ENV"),
       ConfigFile: "somewhere/else/weird-config.yml",
       Conn:       "http:my-elasticsearch.com",
       TargetDir:  "something/different",
    })
    
    // Run all operations
    err := e.RunAll()
    if err != nil {
       panic(err)
    }
}
```
### Config
The SDK does not take into account environment variables, only the passed in config object or your config.yml.
Like the CLI, an order of precedence is used for configuration.

#### Precedence
1. Passed in `esdt.Config` struct
1. `config.yml`

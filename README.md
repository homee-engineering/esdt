# Elasticsearch Data Tool
esdt is a tool for seeding/migrating your Elasticsearch datastore

## Installation
1. Go to the [releases](https://github.com/homee-engineering/esdt/releases)
1. Select your platform and download
1. Move the binary `mv ~/Downloads/*.esdt. /usr/local/bin/esdt`:

## Usage

Generate the `es/operations` directory in your current directory
```bash
esdt gen dir
```
Generate a `Data Template` - a json file containing the Elasticsearch query:
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

## Config

All global flags can be configured via command line flag, environment variable, or `config.yml` in your target
directory (by default `es/operations`)

### Global flags
| Flag     | Env Var             | Config.yml field | Description                                                                                  |
|----------|---------------------|------------------|----------------------------------------------------------------------------------------------|
| `conn`   | `ELASTICSEARCH_URL` | `conn`           | The Elasticsearch base URL to run all operations against. Default is `http://localhost:9200` |
| `dir`    | `ESDT_TARGET_DIR`   | `dir`            | The directory of the data templates. Default is `es/operations`                              |
| `config` | N/A                 | N/A              | The location of your config YAML. Default is ./es/config.tml                                 |
| `env`    | `ESDT_ENV`          | <env\>           | The environment to run the tool in. Default is dev                                           |

### Config.yml
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
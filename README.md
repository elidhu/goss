![logo](https://user-images.githubusercontent.com/22187575/93597198-3f224f80-f9ed-11ea-8b6a-fa3aec24f133.png)

**Goss** is a tool for managing AWS SSM parameters from the CLI. It was mainly developed to managed batches of secrets for application deployment and infrastructure.

## Installation
Not sure yet? this is day 3 of golang!

## AWS auth
Authentication with AWS is pretty standard as this uses the AWS go SDK. Do a google search if you need more information. The places that the SDK looks for credentials are:
- Environment
- `~/.aws/config`
- `~/.aws/credentials`

It is advised to use **goss** in conjuction with **aws-vault** so that your credentials are stored encrypted locally and you just inject them each time you run **goss**, like so.

```
aws-vault exec prod -- goss
```

## Usage
```
For importing and exporting secrets to AWS SSM from a local file for
use in AWS applications

Inspired by another cli tool 'chamber' but with AWS SSM support only. This is to
facillitate syncing path based secrets from TOML files to AWS SSM.

Usage:
  goss [command]

Available Commands:
  delete      Delete parameters from SSM
  help        Help about any command
  import      Import a file into SSM at the given path
  list        List parameters in SSM by path
  put         Put a parameter (or a file of them) into SSM

Flags:
      --config string   config file (default is $HOME/.goss.toml)
  -h, --help            help for goss
      --json            output as json

Use "goss [command] --help" for more information about a command.
```

### List
```bash
goss list /path
```

### Put
```
goss put -n /test/param -v somevalue -t SecureString
```

### Import
```
goss import -f test.env -t SecureString
```

Import allows reading a file into SSM Parameter Store. Currently only .env key-values files are supported. **However** the parsers are already accessible in the code for the other 3 major formats - I just need create a flag to allow a choice of input format.

| File format | Currently supported |
| :---------: | ------------------- |
|   dotenv    | yes                 |
|    json     | soon!               |
|    toml     | soon!               |
|    yaml     | soon!               |

### Delete
```
goss delete -n /test/param
```

#### Obligatory fancy jq pipe
Just some fanciness showing interop with other Unix tools, such as the popular **jq**. This will use **goss** to list the parameters in the store, output as json, filter to the names and pass them to **goss** again to delete.

```bash
goss list / --json | jq '.[].Name' | xargs -n1 goss delete -n
```

## Why?
I made this tool because although **chamber** is an excellent tool - it uses **viper** underneath and the problem with **viper** is that the keys are **CASE INSENSITIVE** which for me was unacceptable. So I decided to roll-my-own using the wonderful [**koanf**](https://github.com/knadh/koanf) library to manage the deserialisation of various config files.
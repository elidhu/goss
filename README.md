![Logo](https://user-images.githubusercontent.com/22187575/93653012-a968dd80-fa49-11ea-947a-537a91282ab2.png)

**goss** is a tool for managing AWS SSM parameters from the CLI. It was mainly developed to manage batches of secrets / parameters stored in local env files for application and infrastructure deployment.

## Installation
### Using go get
```
go get -u github.com/kevinglasson/goss
```
### Pre-built binaries
Download the appropriate binary for your system from the releases page.

## AWS auth
Authentication with AWS is pretty standard as this uses the AWS go SDK. Do a google search if you need more information. The places that the SDK looks for credentials are:
- Environment
- `~/.aws/config`
- `~/.aws/credentials`

It is advised to use **goss** in conjuction with **aws-vault** so that your credentials are stored encrypted locally and you just inject them each time you run **goss**, like so.

```bash
aws-vault exec prod -- goss
```

### Tip
if you are going to run multiple goss commands in a sessions you can start a shell that holds your credentials with.

```bash
# This will put your AWS credentials into the environment
aws-vault exec prod -- bash

# Now proceed to use goss without the aws-vault... prefix
goss list -p /
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
List all parameters at a given path, by default the output is a table with a subset of all of the fields AWS returns (the important ones).

- Parameters can be output as a JSON using the `--json` flag which facilitates interaction with other CLI tools such as **jq**.
- Parameters are returned encrypted by default, use the `-d` flag to have them decrypted.
- Parameters in sub-paths of the specified path are not returned by default, use the `-r` parameter to recursively list the parameters.

#### Default
```bash
goss list -p /dev/test-env -r
```
```md
+------------------------+--------------------------------------+---------+----------------------+
|          NAME          |                VALUE                 | VERSION |       LAST MOD       |
+------------------------+--------------------------------------+---------+----------------------+
| /dev/test-env/COMMENT  | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:10Z |
| /dev/test-env/MORE     | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:10Z |
| /dev/test-env/MiXeD    | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:09Z |
| /dev/test-env/UPPER    | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:09Z |
| /dev/test-env/lower    | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:09Z |
| /dev/test-env/oddChars | AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...¨ |       1 | 2020-09-19T03:35:10Z |
+------------------------+--------------------------------------+---------+----------------------+
```
#### JSON
```
```bash
goss list -p /dev/test-env -r --json
```
```
[
  {
    "ARN": "arn:aws:ssm:ap-southeast-2:XXXXXXXXXXXX:parameter/dev/test-env/COMMENT",
    "DataType": "text",
    "LastModifiedDate": "2020-09-19T03:35:10.111Z",
    "Name": "/dev/test-env/COMMENT",
    "Selector": null,
    "SourceResult": null,
    "Type": "SecureString",
    "Value": "AQICAHhEgSOjHIIiYIkJp/zSBm7c5cy7...",
    "Version": 1
  },
  ...
]

```

### Put
Put a single named parameter into the store. Note that the name, `-n` is the full path to the parameter.

```
goss put -n /test/param -v somevalue -t SecureString
```

### Delete
Delete a single named parameter from the store. Note that the name, `-n` is the full path to the parameter.
```
goss delete -n /test/param
```

#### Obligatory fancy jq pipe
Just some fanciness showing interop with other Unix tools, such as the popular **jq**. This will use **goss** to list the parameters in the store, output as json, filter to the names and pass them to **goss** again to delete.

```bash
goss list -p / --json | jq '.[].Name' | xargs -n1 goss delete -n
```
### Import
Import allows reading a file into the parameter store.

- All parameters from the file must be stored as the same type i.e. String or SecretString etc.
- Currently only .env key-values files are supported. **However** the parsers are already accessible in the code for the other 3 major formats - I just need create a flag to allow a choice of input format.

```
goss import -f test.env -t SecureString
```

#### File format support
| File format | Currently supported |
| :---------: | ------------------- |
|   dotenv    | yes                 |
|    json     | soon!               |
|    toml     | soon!               |
|    yaml     | soon!               |

## Why?
I made this tool because although **chamber** is an excellent tool - it uses **viper** underneath and the problem with **viper** is that the keys are **CASE INSENSITIVE** which for me was unacceptable. So I decided to *roll-my-own* using the wonderful [**koanf**](https://github.com/knadh/koanf) library to manage the deserialisation of various config files.

## Acknowledgements
- [koanf](https://github.com/knadh/koanf): For being the solution to the problems I had.
- [chamber](https://github.com/segmentio/chamber): For being the inspiration.
- [cobra](https://github.com/spf13/cobra): For being an awesome CLI tool library.

[Buy me a ☕!](https://www.paypal.me/kevinglasson)
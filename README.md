# CUE 🤝 Terraform

## Example

```shell
$ cd ./example
$ go run ../cmd/cueform-gen -d
$ cue eval -c -e "infra" -f -o ./cueform.tf.json
$ terraform init
$ terraform plan
```

## Generate provider schema

1. Install `cueform-gen`:
    ```shell
    $ go install github.com/cueformproject/cueform/cmd/cueform-gen@latest
    ```
2. Create a `cueform.gen.cue` file:
    ```cue
    providers: {providers: {
    	"hashicorp/aws": version:       "5.58.0"
    	"hashicorp/random": version:    "3.6.2"
    	"integrations/github": version: "6.2.3"
    }
    ```
3. Execute `cueform-gen`:
    ```shell
    $ cueform-gen
    ```

By default, the provider schema will be generated in the `providers` directory.
Each provider will have its own directory with the schema files.

Check `cueform-gen` configuration schema [here](/cmd/cueform-gen/internal/config/config.cue).

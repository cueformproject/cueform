# CUE 🤝 Terraform

## Generate provider schema

1. Install `cueform-gen`:
    ```shell
    go get -u github.com/cueformproject/cueform/cmd/cueform-gen
    ```
2. Create a `cueform.gen.cue` file:
    ```cue
    providers: {
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

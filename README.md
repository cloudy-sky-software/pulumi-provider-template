# Pulumi Native Provider Template

This repository is a boilerplate showing how to create a native Pulumi provider using:

1. [`pulschema`](https://github.com/cloudy-sky-software/pulschema) to convert an OpenAPI spec to Pulumi schema
2. [`pulumi-provider-framework`](https://github.com/cloudy-sky-software/pulumi-provider-framework) (closed-source) to make HTTP requests against the cloud provider API. It uses the metadata returned by `pulschema` as one of the outputs from
   converting an OpenAPI spec.

## `cloudy-sky-software/pulumi-provider-template` vs. `cloudy-sky-software/pulumi-provider-boilerplate`

This template serves as a base for providers that are built using `pulschema` _and_ `pulumi-provider-framework`.
[`cloudy-sky-software/pulumi-provider-boilerplate`](https://github.com/cloudy-sky-software/pulumi-provider-boilerplate) still serves as a base for a provider that uses `pulschema`
to derive a Pulumi schema from an OpenAPI spec, but does not have any dependency on `pulumi-provider-framework`.
Without `pulumi-provider-framework` it is up to as to how you want to make the HTTP requests to the respective
cloud provider API you are working with. You have to implement Pulumi's `ResourceProviderServer` interface
yourself and handle all the concerns of communicating with the cloud provider API, including performing
diffs, checks etc.

## Background

Follow this link to see what a [Pulumi native provider is comprised of](https://github.com/cloudy-sky-software/cloud-provider-api-conformance/blob/main/pulumi.md).

A Pulumi Resource Provider:

- is a gRPC server which allows for the Pulumi engine to create resources in a specific cloud
- holds the lifecycle logic for these cloud resources
- holds a pulumi JSON schema that describes the provider
- provides language-specific SDKs so resources can be created in whichever language you prefer

When we speak of a "native" provider, we mean that all implementation is native to Pulumi, as opposed
to [Terraform-based providers](https://github.com/pulumi/pulumi-tf-provider-boilerplate).

## Authoring a Pulumi Native Provider

This boilerplate creates a working Pulumi-owned provider named `xyz`.
It implements a random number generator that you can [build and test out for yourself](#test-against-the-example) and then replace the Random code with code specific to your provider.

### Prerequisites

Ensure the following tools are installed and present in your `$PATH`:

- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [Go 1.21](https://golang.org/dl/) or 1.latest
- [NodeJS](https://nodejs.org/en/) 18.x. We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
- [Yarn](https://yarnpkg.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [Python](https://www.python.org/downloads/) (called as `python3`)
- [.NET](https://dotnet.microsoft.com/download)

### Creating and Initializing the Repository

Pulumi offers this repository as a [GitHub template repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) for convenience. From this repository:

1. Click "Use this template".
1. Set the following options:
   - Owner: pulumi
   - Repository name: pulumi-xyz-native (replace "xyz" with the name of your provider)
   - Description: Pulumi provider for xyz
   - Repository type: Public
1. Clone the generated repository.

From the templated repository:

1. Search-replace `xyz` with the name of your desired provider.

#### Build the provider and install the plugin

```bash
$ make build install
```

This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### Test against the example

```bash
$ cd examples/simple
$ yarn link @pulumi/xyz
$ yarn install
$ pulumi stack init test
$ pulumi up
```

Now that you have completed all of the above steps, you have a working provider that generates a random string for you.

#### A brief repository overview

You now have:

1. A `provider/` folder containing the building and implementation logic
   1. `cmd/`
      1. `pulumi-gen-xyz/` - generates language SDKs from the schema
      2. `pulumi-resource-xyz/` - holds the package schema, injects the package version, and starts the gRPC server
   2. `pkg`
      1. `provider` - holds the gRPC methods (and for now, the sample implementation logic) required by the Pulumi engine
      2. `version` - semver package to be consumed by build processes
2. `deployment-templates` - a set of files to help you around deployment and publication
3. `sdk` - holds the generated code libraries created by `pulumi-gen-xyz/main.go`
4. `examples` a folder of Pulumi programs to try locally and/or use in CI.
5. A `Makefile` and this `README`.

### Implementing the provider callback methods

You will find a mostly blank implementation of these in `pkg/provider/provider.go`.
Note that these methods do not link 1:1 to the Pulumi resource provider interface
because `pulumi-provider-framework` provides a convenient callback mechanism
and handles all other responsibilities. You should use the callback methods
to alter the HTTP request (in Pre* methods) or the response (in Post*) as
needed. If you don't need any customization, then you don't need to do
anything at all. Yep! The framework handles it all.

### Build Examples

Create an example program using the resources defined in your provider, and place it in the `examples/` folder.

You can now repeat the steps for [build, install, and test](#test-against-the-example).

## Documentation

Please [follow this guide to add documentation to your provider](https://www.pulumi.com/docs/guides/pulumi-packages/how-to-author/#write-documentation).

## Configuring CI and releases

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

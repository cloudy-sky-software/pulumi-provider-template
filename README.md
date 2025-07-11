# Pulumi Provider Template For OpenAPI-based Providers

This repository is a boilerplate showing how to create a native Pulumi provider using:

1. [`pulschema`](https://github.com/cloudy-sky-software/pulschema) to convert an OpenAPI spec to Pulumi schema
2. [`pulumi-provider-framework`](https://github.com/cloudy-sky-software/pulumi-provider-framework) to make HTTP requests against the cloud provider API. It uses the metadata returned by `pulschema` as one of the outputs from
   converting an OpenAPI spec.

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
   - Owner: `<your GH organization>`
   - Repository name: pulumi-xyz (replace "xyz" with the name of your provider)
     - Providers built from Cloudy Sky Software's templates are _always_ native providers, by default.
     - However, if there is already a TF-bridged provider with that name, you should add the suffix `-native` so that the package name in some package registries do not conflict with the other providers.
   - Description: Pulumi provider for xyz
   - Repository type: Public
1. Clone the generated repository.

From the templated repository:

Search-replace the following occurrences with the corresponding names.

| Search | Replace With                                                                   |
| ------ | ------------------------------------------------------------------------------ |
| xyz    | Lower-cased name of the provider                                               |
| Xyz    | Pascal-case name of the provider                                               |
| XYZ    | Upper-cased name of the provider                                               |
| x_y_z  | Lower snake-cased name of the provider if the provider name has multiple words |

#### Embed the OpenAPI spec

The OpenAPI spec file for the provider you are building must be placed in the `provider/cmd/pulumi-gen-*` folder as `openapi.yml`.
Unlike some of Pulumi's own native providers which download the OpenAPI spec from an upstream repo, this template does not do that
as it does not know where to download the OpenAPI spec from.

- You can, of course, add a Make target similar to what Pulumi does with some of its native providers and have it download the latest
OpenAPI spec from an upstream repo.
- You can also rename the file to something other than `openapi.yml` if you wish. Be sure to change the name of the file that Go
should embed in `provider/cmd/pulumi-gen-*/main.go`.

#### Generate Pulumi schema

Now that you have an OpenAPI spec downloaded and all the placholders renamed to the appropriate provider/package name,
i.e. replaced `xyz` with the appropriate name (see the section above), you can generate a Pulumi schema by running
`make gen generate_schema`. You must have a Pulumi schema generated successfully in order to generate the language
SDKs.

The larger the spec the more likely there are errors in the spec itself. In most cases, it has nothing to do with `pulschema`.
If it's a genuine bug in `pulschema`, please open an [issue](https://github.com/cloudy-sky-software/pulschema/issues).
But it's more than likely you'll need to patch the OpenAPI spec itself. You can do that using Go instead of manually editing the spec file
which can be quite cumbersome, especially if you are dealing with a very large spec.
Anyway, here's where you can write Go code to modify the spec: https://github.com/cloudy-sky-software/pulumi-provider-template/blob/main/provider/pkg/gen/openapi_fixes.go.
Here's an example of an OpenAPI spec that needed to be modified: https://github.com/cloudy-sky-software/pulumi-digitalocean-native/blob/main/provider/pkg/gen/openapi_fixes.go

If there are endpoints in the spec that you don't care about and want to exclude them from Pulumi,
you can pass a [list](https://github.com/cloudy-sky-software/pulumi-provider-template/blob/main/provider/pkg/gen/schema.go#L94) of the endpoint paths exactly as they appear in the spec.

#### Build the provider and install the plugin

```bash
$ make build install
```

This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

Feel free to modify any of the Make targets (or add news ones) to fit your needs.
If you feel others might find them useful, please consider contributing it back to
this template repo. :)

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

### Importing Existing Resources

Import IDs should satisfy all ID segments in the `GET` endpoint for the resource
you are importing. The IDs required in the path should be separated by `/`.
First, start by identifying the `GET` endpoint in the OpenAPI spec
for the provider.

For example, let's assume such a GET endpoint path for some resource is: `/services/{serviceId}/someResource/{someResourceId}`.

Thus, the `pulumi import` command to run is:

```bash
# The resource type token can be easily found by using your IDEs
# Go To Definition functionality for the resource and looking at the type
# property defined in the custom resource's class definition.
pulumi import {resourceTypeToken} {resourceName} /{serviceId}/{someResourceId}
```

Alternatively, you can also import using the `import` Pulumi resource option.
Run `pulumi up` to import the resource into your stack's state. Once imported,
you should remove the `import` resource option.

```typescript
const someResource = new SomeResource(
  "myResourceName",
  { //inputs for the reosurce },
  {
    protect: true,
    import: `/{serviceId}/{someResourceId}`,
  }
);
```

Refer to the Pulumi [docs](https://www.pulumi.com/docs/iac/adopting-pulumi/import/) for importing a
resource.

## Configuring CI and releases

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

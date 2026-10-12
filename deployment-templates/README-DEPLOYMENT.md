# Configure

1. Obtain needed accounts from the language specific package managers you will use, i.e. `npm`, `pypi`, `NuGet` etc.

1. Customize the release.yml with the correct tokens using the format:

1. Customize .goreleaser.yml for your provider, replacing any instances of 'xyz' with your provider's name, and paying special attention that the `ldlflags` are set to match your provider/go.mod exactly:

   `-X github.com/pulumi/pulumi-my-provider/provider/pkg/version.Version={{.Tag}}`

1. Delete this directory if desired

## Trusted Publishing

Almost all language registries now require (or strongly suggest) using Trusted Publishing.
It's a process that allows you to exchange the OIDC token issued by your CI/CD service, GitHub,
in this case, for a short-lived token issued by the respective language registry.

Follow these guides to learn more about Trusted Publishing:

- [npm](https://docs.npmjs.com/trusted-publishers)
  - Requires generating a fine-grained access token in `npm` that will be used to **stage** the package.
- [NuGet](https://learn.microsoft.com/en-us/nuget/nuget-org/trusted-publishing)
- [pypi](https://docs.pypi.org/trusted-publishers/using-a-publisher/)

⚠️ You need to complete the steps to establish the trust in the respective
registry _before_ you can publish a version.

## Releasing A New Version

:info: Switch to the `main` branch first and get the latest `git pull origin main && git fetch`. Check what the last release tag was.

**Note**: The following steps are simply a suggestion. The only requirement is the tag formats. Feel free to follow a different
release process that suits your needs.

1. Regular releases should just increment the patch version unless a minor or a major (breaking changes) version bump is warranted.
1. Update the `CHANGELOG.md` with notes about what will be included in this release.
1. Commit the changelog with `git commit -am "vX.Y.Z"` or something similar and push `git push origin main`.
1. Tag the commit with the release version by running

   ```bash
   git tag vX.Y.Z
   git tag sdk/vX.Y.Z
   ```

1. Push the tags.

   ```bash
   git push --tags
   ```


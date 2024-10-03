# Releasing a new `wit-bindgen-go` version

This document describes the steps to release a new version of the `wit-bindgen-go` CLI.

## 1. Update the [CHANGELOG.md](./CHANGELOG.md)

Update the `CHANGELOG.md` file with the changes that are part of the new release, and make sure `Unreleased` is renamed to the new version number. Make a PR with these changes.

## 2. Create a new release

Once the PR is merged, create a new release on GitHub with the same version number as the one in the `CHANGELOG.md` file as the tag name.

```sh
git tag -a v0.3.0 -m "Release v0.3.0"
```

Push the tag to GitHub.

```sh
git push upstream v0.3.0
```

After the tag is pushed, GitHub Actions will automatically create a new release with the content of the `CHANGELOG.md` file.

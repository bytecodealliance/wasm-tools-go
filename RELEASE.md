# Releasing a new `wit-bindgen-go` version

This document describes the steps to release a new version of the `wit-bindgen-go` CLI.

## 1. Update the [CHANGELOG.md](./CHANGELOG.md)

* Add the latest changes to CHANGELOG.md.
* Rename the Unreleased section to reflect the new version number.
* Submit a pull request (PR) with these updates.

## 2. Create a new release

Once the PR is merged, tag the new version in Git:

```sh
git tag -a v0.3.0 -m "Release v0.3.0"
```

Push the tag to GitHub:

```sh
git push upstream v0.3.0
```

After the tag is pushed, GitHub Actions will automatically create a new release with the content of the `CHANGELOG.md` file.

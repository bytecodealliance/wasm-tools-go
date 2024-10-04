# Release

This document describes the steps to release a new version of the `wit-bindgen-go` CLI.

## 1. Update the [CHANGELOG.md](./CHANGELOG.md)

* Add the latest changes to CHANGELOG.md.
* Rename the Unreleased section to reflect the new version number.
	* Update the links to new version tag in the footer of CHANGELOG.md
* Submit a pull request (PR) with these updates.

## 2. Create a new release

Once the PR is merged, tag the new version in Git and push the tag to GitHub.

For example, to tag version `v0.3.0`:

```sh
git tag v0.3.0
git push origin v0.3.0
```

After the tag is pushed, GitHub Actions will automatically create a new release with the content of the `CHANGELOG.md` file.

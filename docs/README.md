# go.bytecodealliance.org Vanity URLs

This directory hosts the configuration for Go vanity URLs under `go.bytecodealliance.org` domain.

More information about vanity URLs can be found in the [Go documentation](https://golang.org/cmd/go/#hdr-Remote_import_paths).

## Adding a new vanity URL

If you want to create a vanity URL `go.bytecodealliance.org/foo` pointing to `github.com/bytecodealliance/go-foo`:

- Create directory `foo`
- Create file `foo/index.html` with the following content:

     ```html
     <html>
     	<head>
     		<meta name="go-import" content="go.bytecodealliance.org/foo git https://github.com/bytecodealliance/go-foo" />
     	</head>
     </html>
     ```

NOTE: Vanity URLs must be unique and not shadow existing package names in this repository. For ex: Creating a vanity URL `go.bytecodealliance.org/cm` would shadow the `cm` package in this repository.

# DizGit

DizGit CLI tool for getting latest GitHub Repository Releases and/or latest Pull Requests

Go standard library has a [flag pkg](https://pkg.go.dev/flag) that makes building CLIs pretty painless, and was used here.

## Prereqs:

- Have Go (v1.16.6) installed and GOPATH and or GOROOT configured: https://golang.org/doc/gopath_code

## Building Executable:

From within this directory run `go build dizgit.go`, and it should create the `dizgit` executable.
Optionally, you can move the executable to `/usr/local/bin` to make it available outside of this directory.

## Usage:

### Get 3 Latest Releases:

```golang
dizgit --repo https://github.com/user/repo.git --releases
```

### Get 3 Latest Pull Requests:

```golang
dizgit --repo https://github.com/user/repo.git --pullrequests
```

### Get 3 Latest Releases and Pull Requests:

```golang
dizgit --repo https://github.com/user/repo.git --releases --pullrequests
```

# Gist

> Command line tool for publishing gists

## Usage:

``` sh
# read from stdin
cat file.sh | gist

# set file name
cat file.sh | gist -f "myfile.sh"

# make public
cat file.sh | gist -p

# multiple files
gist *.js
```

## Install:

``` sh
go get github.com/woremacx/gist
```

For auth, the tool looks for an environment variable called `GITHUB_TOKEN_FOR_GIST`
You can generate one at: https://github.com/settings/tokens

``` sh
export GITHUB_TOKEN_FOR_GIST="blah blah blah"
```

## GitHubGist Enterprise:
To use GitHubGist Enterprise, set environment variables: `GIST_ENTERPRISE_BASE_URL` and `GIST_ENTERPRISE_UPLOAD_URL`

``` sh
export GIST_ENTERPRISE_BASE_URL="https://git.example.com/api/v3/"
export GIST_ENTERPRISE_UPLOAD_URL="https://git.example.com/api/v3/"
```

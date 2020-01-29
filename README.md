# GittyUp 

GittyUp automates updating manifests files in GitOps repositories. A typical use case is to use it as part of a CI pipeline to update versions in a manifest file after changes are committed for release. It currently supports updating JSON, YAML and HCL. It works by cloning a repo, parsing a manifest file, update values and creating a branch with the changes. If it is a GitHub repo it will also create a pull request.

## Installing

### Brew

    brew install liatrio/tap/gitty-up

### Executable

Download an executable from [releases](https://github.com/liatrio/gitty-up/releases)

### Docker

    docker pull liatrio/gitty-up:latest

## Usage



### Arguments

Most arguments can be set as flags or with environment variables.

- `--gitUrl` | `GITOPS_GIT_URL` (required)
  - URL of git repository. Example: `https://github.com/my-org/my-repo` or `git@github.com/my-org/my-repo.git`
- `--gitUsername` | `GITOPS_GIT_USERNAME` (required)
  - Username to authenticate with git.
- `--gitPassword` | `GITOPS_GIT_PASSWORD` (required)
  - Password or token to authenticate with git.
- `--repoFile` | `GITOPS_REPO_FILE` (required)
  - File in git repository to apply changes to. Example `testing/application.json`
- `--values` | `GITOPS_VALUES` (required)
  - List of variables and corresponding values to update. Variable paths are a list of keys seperated with periods. Each variable is separating with a colon. Example 'input.builder_images_version=${VERSION}:inputs.jenkins_image_version=${VERSION}'
- `--dry-run`
  - Do not push up branch or create pull request

### Command line example

    gitty-up -gitUrl=https://github.com/my-org/my-repo -gitUsername=USERNAME -gitPassword=PASSWORD -repoFile=testing/application.json --values=main.version=v0.0.42

This will change the value of _main.version_ to _v0.0.42_ in the file _testing/application.json_ for the repo _https://github.com/my-org/my-repo_

### Docker example
 
    docker run liatrio/gitty-up --gitUrl=https://github.com/my-org/my-repo --gitUsername=USERNAME --gitPassword=PASSWORD --repoFile=testing/application.json --values=version=v0.0.42

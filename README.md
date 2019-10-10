# Gitops Builder Image
The gitops builder image contains a script written in golang that can parse and edit HCL files in a specified repository and then commit the changes and trigger a PR in that same repository

### Environment Variables

The image requires a few environment variables to be set in order to run properly.

- `GITOPS_GIT_URL`
  - URL of git repository. Example: `https://github.com/liatrio/lead-environments` or `git@github.com/liatrio/lead-environments.git`
- `GITOPS_GIT_USERNAME`
  - Username to authenticate with git.
- `GITOPS_GIT_PASSWORD`
  - Password or token to authenticate with git.
- `GITOPS_REPO_FILE`
  - File in git repository to apply changes to. Example `terragrunt.hcl`
- `GITOPS_VALUES`
  - List of variables and corresponding values to update. Variable paths are a list of keys seperated with periods. Each variable is separating with a colon. Example 'input.builder_images_version=${VERSION}:inputs.jenkins_image_version=${VERSION}'


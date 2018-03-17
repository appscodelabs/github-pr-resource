[![Go Report Card](https://goreportcard.com/badge/github.com/appscodelabs/github-pr-resource)](https://goreportcard.com/report/github.com/appscodelabs/github-pr-resource)
[![Build Status](https://travis-ci.org/appscodelabs/github-pr-resource.svg?branch=master)](https://travis-ci.org/appscodelabs/github-pr-resource)
[![codecov](https://codecov.io/gh/appscodelabs/github-pr-resource/branch/master/graph/badge.svg)](https://codecov.io/gh/appscodelabs/github-pr-resource)
[![Docker Pulls](https://img.shields.io/docker/pulls/appscodelabs/github-pr-resource.svg)](https://hub.docker.com/r/appscodelabs/github-pr-resource/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)
# Concourse-ci Resource for Git Pull Requests

This resource can check for new pull requests and run test on them. After finishing tests, it can update status to pending, success or failure.

You can define organization, then only pr-s from the members of the organization will run. You can also define labels. After checking pr-s from users not in the org, you can add labels like `ok-to-test` and concourse will run test automatically. After running tests, it'll remove the label from that pr, so all future commits to that pr will not run automatically.

## Deploying to concourse

```
resource_types:
- name: pull-request
  type: docker-image
  source:
    repository: appscodelabs/github-pr-resource
    tag: 1.0.0

```

### Source Configuration

* `owner`: *Required.* example:`appscode`
* `repo`: *Required.* example: `github-pr-resource`
* `access_token`: *Required.* It is needed to change the status of pr
* `label`: Optional..
* `org`: Optional.


## Behaviour

### `check`: Checks for new pull requests

Checks for You must define `version: every` to ensure checking every pr

### `in`: Clones the repository at the given pull request ref
### `out`: Update the status of a pull request

#### Parameters

* `path`: *Required.*
* `status`: *Required.* Status can only be `success`, `failure`, `error`, or `pending`.

## Example

```
resource_types:
- name: pull-request
  type: docker-image
  source:
    repository: appscodelabs/github-pr-resource
    tag: 1.0.0

resources:
- name: pull-request
  type: pull-request
  source:
    owner: appscode
    repo: guard
    access_token: ((access_token))
    label: ok-to-test
    org: appscode

jobs:
- name: test-pr
  plan:
  - get: pull-request
    trigger: true
    version: every
  - put: pull-request
    params:
      path: pull-request
      status: pending
  - task: test-pr
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: ubuntu
      inputs:
      - name: pull-request
      run:
        path: echo
        args: ["hello world"]
    on_success:
      put: pull-request
      params:
        path: pull-request
        status: success
    on_failure:
      put: pull-request
      params:
        path: pull-request
        status: failure
```

## Testing the Code on Your Own Local Environment

Please copy `out/find_hash.sh` `out/fetch_pr.sh` `in/git_script.sh` to your rood directory (`/`)

Or,

in `in/main.go`, change `exec.Command("/git_script.sh"...` to `exec.Command("./git_script.sh"...`

in `out/main.go` change `exec.Command("/find_hash.sh",...` to `exec.Command("./find_hash.sh",...` and `exec.Command("/fetch_pr.sh",...` to `exec.Command("./fetch_pr.sh",...`

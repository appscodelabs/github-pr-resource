# Concourse-ci Resource for Git Pull Requests

This resource can check for new pull requests and run test on them. After finishing tests, it can update status to pending, success or failure.

You can define organization, then only pr-s from the members of the organization will run. You can also define labels. After checking pr-s from users not in the org, you can add labels like `ok-to-test` and concourse will run test automatically. After running tests, it'll remove the label from that pr, so all future commits to that pr will not run automatically.

## Deploying to concourse

```
resource_types:
- name: pull-request
  type: docker-image
  source:
    repository: tahsin/git-pull-resource
    tag: 1.0.0

```

### Source Configuration

* `owner`: *Required.* example:`tahsinrahman`
* `repo`: *Required.* example: `concourse-git-pr-resource`
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
    repository: tahsin/git-pull-resource
    tag: 1.0.0

resources:
- name: pull-request
  type: pull-request
  source:
    owner: tahsinrahman
    repo: test-status
    access_token: ((access_token))
    label: ok-to-test
    org: your_org

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
        path: pull-request/ci/test.sh
        args: []
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

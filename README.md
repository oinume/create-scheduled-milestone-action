# create-scheduled-milestone-action

A GitHub Action which can create a milestone with given title, state, description and due date. This GitHub Action is inspired by [github-action-create-milestone](https://github.com/marketplace/actions/create-milestone). The action doesn't support setting due_on and state so I created another one.

## Options

This action supports following options.

### title

The title (version) for the milestone.

- *Required*: `yes`
- *Type*: `string`
- *Example*: `v1.0.0`

### state

The state of the milestone. Either `open` or `closed`. Default: `open`

- *Required*: `no`
- *Type*: `string`
- *Example*: `open`

### description

An optional description for the milestone.

- *Required*: `no`
- *Type*: `string`

### due_on

An optional due date for the milestone. This is a timestamp in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format: `YYYY-MM-DDTHH:MM:SSZ`.

* *Required*: `no`
* *Type*: `string`

## Output

Set some output data which from [API response](https://developer.github.com/v3/issues/milestones/#response)

### number

The number of created milestone.

## Example

```yaml
name: Create Milestone
on:
  push:
    branches:
      - develop

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        with:
          ref: develop

      - name: "Set due_on"
        id: set-due-on
        run: echo "::set-output name=due_on::$(date --iso-8601=seconds -d '13 days')"

      - name: "Create a new milestone"
        id: create-milestone
        uses: oinume/create-scheduled-milestone-action@v1.0.0
        with:
          title: "1.0.0"
          state: "open"
          description: "v1.0.0"
          due_on: "${{ steps.set-due-on.outputs.due_on }}"
        env:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

      - name: "Output milestone number"
        run: echo ${{ steps.create-milestone.outputs.number }}
```
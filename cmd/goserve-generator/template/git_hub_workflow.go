package template

const GitHubWorkflow = `name: Testing

on:
  push:
    branches:
      - "*"
      - feat/**
jobs:
  validation:
    name: "Build validation"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: 1.24.2
      - name: Test
        run: |
          go test ./...
`

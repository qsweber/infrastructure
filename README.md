# infrastructure

![CI](https://github.com/qsweber/infrastructure/actions/workflows/production.yml/badge.svg)

An AWS CDK (Go) app with a single stack, `github-actions-deploy` (account `120356305272`,
`us-west-2`), that lets GitHub Actions deploy to AWS without storing AWS keys as secrets:

- an IAM OIDC provider for `token.actions.githubusercontent.com`
- the `github-actions-deploy` IAM role, which workflows in the repos listed in
  `githubDeployRepos` (`github_oidc.go`) can assume

Pull requests run `gofmt`, `go vet`, `cdk synth`, and `cdk diff`; merges to `main` run
`cdk deploy`. Both assume the role themselves, so this repo needs no AWS secrets either.

## Using the role from a workflow

```yaml
permissions:
  id-token: write
  contents: read

steps:
  - uses: aws-actions/configure-aws-credentials@v4
    with:
      role-to-assume: arn:aws:iam::120356305272:role/github-actions-deploy
      aws-region: us-west-2
```

To allow a new repo, add it to `githubDeployRepos` and merge to `main`.

## Local usage

Prerequisites: Go 1.25.8+, Node.js 22+, AWS credentials.

```bash
npm install
npx cdk diff github-actions-deploy
npx cdk deploy github-actions-deploy
```

The very first deploy has to be run this way, with admin credentials, because CI uses the role
that this stack creates.

package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

const (
	account = "120356305272"
	region  = "us-west-2"
)

func main() {
	app := awscdk.NewApp(nil)

	NewGitHubActionsDeployStack(app, "github-actions-deploy", &awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: jsii.String(account),
			Region:  jsii.String(region),
		},
		Description: jsii.String("GitHub Actions OIDC provider and deploy role"),
	})

	app.Synth(nil)
}

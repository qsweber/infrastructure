package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// githubDeployRepos are the GitHub repositories whose Actions workflows may assume the deploy
// role. Add a repo here (and merge to main) to let its workflows deploy to AWS without any
// AWS secrets.
var githubDeployRepos = []string{
	"qsweber/infrastructure",
	"qsweber/phone-numbers",
	"qsweber/phone-numbers-web",
	"qsweber/website-template",
	"qsweber/go-template",
}

// createGitHubActionsDeployRole sets up GitHub Actions OIDC federation so workflows get
// short-lived AWS credentials by assuming a role, instead of using long-lived access keys
// stored as repository secrets.
func createGitHubActionsDeployRole(scope constructs.Construct) awsiam.CfnRole {
	provider := awsiam.NewCfnOIDCProvider(scope, jsii.String("GitHubActionsOidcProvider"), &awsiam.CfnOIDCProviderProps{
		Url:          jsii.String("https://token.actions.githubusercontent.com"),
		ClientIdList: jsii.Strings("sts.amazonaws.com"),
	})

	subjects := make([]interface{}, 0, len(githubDeployRepos))
	for _, repo := range githubDeployRepos {
		// Any ref, pull request, or environment in the repo.
		subjects = append(subjects, fmt.Sprintf("repo:%s:*", repo))
	}

	role := awsiam.NewCfnRole(scope, jsii.String("GitHubActionsDeployRole"), &awsiam.CfnRoleProps{
		RoleName:           jsii.String("github-actions-deploy"),
		Description:        jsii.String("Assumed by GitHub Actions via OIDC to deploy CDK stacks"),
		MaxSessionDuration: jsii.Number(3600),
		AssumeRolePolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []interface{}{
				map[string]interface{}{
					"Effect":    "Allow",
					"Principal": map[string]interface{}{"Federated": provider.AttrArn()},
					"Action":    "sts:AssumeRoleWithWebIdentity",
					"Condition": map[string]interface{}{
						"StringEquals": map[string]interface{}{
							"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
						},
						"StringLike": map[string]interface{}{
							"token.actions.githubusercontent.com:sub": subjects,
						},
					},
				},
			},
		},
		// Workflows run `cdk deploy` (which creates IAM roles, Lambdas, API Gateways, S3
		// buckets, CloudFront distributions, Route 53 records, ...) plus direct S3 sync,
		// CloudFront invalidation, and CloudFormation describe calls. Scope down if desired.
		ManagedPolicyArns: &[]interface{}{jsii.String("arn:aws:iam::aws:policy/AdministratorAccess")},
	})

	return role
}

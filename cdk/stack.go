package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewGitHubActionsDeployStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, jsii.String(id), props)

	deployRole := createGitHubActionsDeployRole(stack)

	awscdk.NewCfnOutput(stack, jsii.String("GithubActionsDeployRoleArn"), &awscdk.CfnOutputProps{
		Value: deployRole.AttrArn(),
	})

	return stack
}

/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at
    http://aws.amazon.com/apache2.0/
   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Usage:
// go run sts_assume_role.go
func main() {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		fmt.Println("NewSession Error", err)
		return
	}

	// Create a STS client
	svc := sts.New(sess)

	roleToAssumeArn := "arn:aws:iam::975156237701:role/CoS-Secrets-Provisioner"
	sessionName := "test_session"
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &roleToAssumeArn,
		RoleSessionName: &sessionName,
	})

	if err != nil {
		fmt.Println("AssumeRole Error", err)
		return
	}

	fmt.Println("Assumed Credentials: ", result)

	fmt.Println("Access ID: ", result.Credentials.AccessKeyId)

	newSession, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(*result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken),
	})

	if err != nil {
		fmt.Println("NewSession Error", err)
		return
	}

	newSVC := sts.New(newSession)
	res, err := newSVC.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		fmt.Println("GetCallerIdentity Error", err)
		return
	}

	fmt.Println("Identity: ", res)

}

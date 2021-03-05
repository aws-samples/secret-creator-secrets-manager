/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: MIT-0
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the Software
 * without restriction, including without limitation the rights to use, copy, modify,
 * merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
 * SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
 
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.uber.org/ratelimit"
	"time"
)

// struct to build response for Client
type BodyResponse struct {
	ARN       string `json:"arn"`
	Name      string `json:"name"`
	VersionId string `json:"versionid"`
	Error     string `json:"error"`
}

type BodyResponses struct {
	Collection []BodyResponse
}

type Secret struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Secrets struct {
	Collection []Secret
}

// Handler function Using AWS Lambda Proxy Request
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secrets := make([]Secret, 10)

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &secrets)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	/*	fmt.Printf("%#v\n", secrets)
		fmt.Println("secrets length", len(secrets))*/

	bodyResponses := createSecrets(secrets)

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponses)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}

// Create secrets in AWS secrets manager, with rate limits
func createSecrets(secrets []Secret) []BodyResponse {
	bodyResponses := make([]BodyResponse, len(secrets))

	rl := ratelimit.New(50) // per second

	prev := time.Now()
	for i := 0; i < len(secrets); i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))

		result, err := createSecret(secrets[i])
		if err != nil {
			bodyResponses[i].Error = err.Error()
		}
		bodyResponses[i].ARN = aws.StringValue(result.ARN)
		bodyResponses[i].Name = aws.StringValue(result.Name)
		bodyResponses[i].VersionId = aws.StringValue(result.VersionId)
		prev = now
	}

	return bodyResponses
}

func createSecret(secret Secret) (*secretsmanager.CreateSecretOutput, error) {
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.CreateSecretInput{
		Name:         aws.String(secret.Name),
		SecretString: aws.String(fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", secret.Username, secret.Password)),
	}

	result, err := svc.CreateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	// fmt.Println(result)
	return result, err
}

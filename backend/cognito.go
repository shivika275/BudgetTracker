package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	cognitoClientID = "25aclbpg1cu32bu52k2fp7e1ga"
	cognitoRegion   = "us-west-2"
)

var cognitoSvc *cognitoidentityprovider.CognitoIdentityProvider

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cognitoRegion),
	})
	if err != nil {
		panic(err)
	}

	cognitoSvc = cognitoidentityprovider.New(sess)
}

func CreateUserInCognito(registerData RegisterData) error {
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(cognitoClientID),
		Username: aws.String(registerData.Username),
		Password: aws.String(registerData.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(registerData.Email),
			},
		},
	}

	_, err := cognitoSvc.SignUp(signUpInput)
	if err != nil {
		return fmt.Errorf("failed to create user in cognito: %v", err)
	}
	return nil
}

func AuthenticateUser(loginData LoginData) (*AuthResult, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(cognitoClientID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(loginData.Username),
			"PASSWORD": aws.String(loginData.Password),
		},
	}

	result, err := cognitoSvc.InitiateAuth(input)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %v", err)
	}

	if result.AuthenticationResult == nil {
		return nil, fmt.Errorf("authentication result is nil")
	}

	authResult := &AuthResult{
		AccessToken:  aws.StringValue(result.AuthenticationResult.AccessToken),
		IdToken:      aws.StringValue(result.AuthenticationResult.IdToken),
		RefreshToken: aws.StringValue(result.AuthenticationResult.RefreshToken),
		ExpiresIn:    aws.Int64Value(result.AuthenticationResult.ExpiresIn),
		TokenType:    aws.StringValue(result.AuthenticationResult.TokenType),
	}

	return authResult, nil
}

func ValidateToken(token string) (bool, error) {
	// Validate the token using Cognito
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	_, err := cognitoSvc.GetUser(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

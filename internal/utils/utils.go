package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	structs_aws "discord_bot/internal/struct_aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	uuid "github.com/google/uuid"
)

func Uuid() string {
	return uuid.New().String()
}

func DeleteTextInString(text string, oldText string) string {
	return strings.Replace(text, oldText, "", -1)
}

// function to masking sensitive string
func MaskSensitiveString(text string) string {
	xString := "xxxxxx"
	return xString + text[len(text)-10:]
}

func GetSecretAws() (structs_aws.SecretAws, error) {
	// Create a new session with the specified region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"), // Replace with your desired AWS region
	})
	if err != nil {
		fmt.Println("Error creating new session: ", err)
		return structs_aws.SecretAws{}, err
	}

	// Create Secrets Manager client
	svc := secretsmanager.New(sess)

	// Get AWS credentials from Secrets Manager
	secretName := "aws-credentials-pro" // Replace with your actual secret name
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Println("Error getting secret value: ", err)
		return structs_aws.SecretAws{}, err
	}

	secret := structs_aws.SecretAws{}
	err = json.Unmarshal([]byte(*result.SecretString), &secret)

	return secret, nil
}

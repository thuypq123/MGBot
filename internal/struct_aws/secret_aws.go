package structs_aws

type SecretAws struct {
	AccessKeyId     string `json:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"AWS_SECRET_ACCESS_KEY"`
	BotToken        string `json:"DISCORD_BOT_TOKEN"`
	Region          string `json:"REGION"`
}

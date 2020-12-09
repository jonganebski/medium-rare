package myaws

import (
	"home/jonganebski/github/medium-rare/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ConnectAws connects to the bucket of this site
func ConnectAws() *session.Session {
	secretAccessKey := config.Config("AWS_SECRET_ACCESS_KEY")
	myRegion := config.Config("AWS_REGION")
	accessKeyID := config.Config("AWS_ACCESS_KEY_ID")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(myRegion),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	if err != nil {
		panic(err)
	}

	return sess
}

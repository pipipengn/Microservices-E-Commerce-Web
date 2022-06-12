package goods

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	region          = ""
	accessKerId     = ""
	secretAccessKey = ""
	buckerName      = ""
)

func GetPresignedS3Url(c *gin.Context) {
	options := s3.Options{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKerId, secretAccessKey, "")),
	}
	client := s3.New(options)
	psClient := s3.NewPresignClient(client, s3.WithPresignExpires(time.Minute))

	tempFileName := fmt.Sprintf("%s.jpg", uuid.NewV4())
	putObject, err := psClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(buckerName),
		Key:    aws.String(tempFileName),
	})

	if err != nil {
		fmt.Println("Got an error retrieving pre-signed object:")
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"url": putObject.URL,
	})
}

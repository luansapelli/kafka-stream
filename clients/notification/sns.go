package notification

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type SNS struct {
	client   snsiface.SNSAPI
	topicArn string
}

// NewSNS creates SNS client
func NewSNS(topicArn string) *SNS {
	return &SNS{
		client:   createSnsClient(),
		topicArn: topicArn,
	}
}

func (notification *SNS) Write(content []byte) (int, error) {
	_, err := notification.client.Publish(&sns.PublishInput{
		Message:  aws.String(string(content)),
		TopicArn: &notification.topicArn,
	})
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func createSnsClient() snsiface.SNSAPI {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sns.New(awsSession)
}

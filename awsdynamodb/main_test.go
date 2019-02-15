package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams/dynamodbstreamsiface"
	"github.com/stretchr/testify/assert"
)

type mockedDynamoStreamsClient struct {
	dynamodbstreamsiface.DynamoDBStreamsAPI
	listStreamsOutput      dynamodbstreams.ListStreamsOutput
	listStreamsOutputError error

	describeStreamOutput      dynamodbstreams.DescribeStreamOutput
	describeStreamOutputError error

	getShardIteratorOutput      dynamodbstreams.GetShardIteratorOutput
	getShardIteratorOutputError error

	getRecordsOutput      dynamodbstreams.GetRecordsOutput
	getRecordsOutputError error
}

func (m mockedDynamoStreamsClient) ListStreams(in *dynamodbstreams.ListStreamsInput) (*dynamodbstreams.ListStreamsOutput, error) {
	return &m.listStreamsOutput, m.listStreamsOutputError
}

func (m mockedDynamoStreamsClient) DescribeStream(in *dynamodbstreams.DescribeStreamInput) (*dynamodbstreams.DescribeStreamOutput, error) {
	return &m.describeStreamOutput, m.describeStreamOutputError
}

func (m mockedDynamoStreamsClient) GetShardIterator(in *dynamodbstreams.GetShardIteratorInput) (*dynamodbstreams.GetShardIteratorOutput, error) {
	return &m.getShardIteratorOutput, m.getShardIteratorOutputError
}

func (m mockedDynamoStreamsClient) GetRecords(in *dynamodbstreams.GetRecordsInput) (*dynamodbstreams.GetRecordsOutput, error) {
	return &m.getRecordsOutput, m.getRecordsOutputError
}

func TestGetStreams(t *testing.T) {
	c := Client{
		StreamsClient: mockedDynamoStreamsClient{
			listStreamsOutput:      dynamodbstreams.ListStreamsOutput{},
			listStreamsOutputError: errors.New("getstreams failed"),
		},
	}

	streams, err := c.getStreams()
	assert.Error(t, err)
	assert.Equal(t, 0, len(streams))

	c = Client{
		StreamsClient: mockedDynamoStreamsClient{
			listStreamsOutput: dynamodbstreams.ListStreamsOutput{
				Streams: []*dynamodbstreams.Stream{{}, {}},
			},
			listStreamsOutputError: nil,
		},
	}

	streams, err = c.getStreams()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(streams))
}

func TestGetStreamsDescriptions(t *testing.T) {

	streams := []*dynamodbstreams.Stream{{}}

	c := Client{
		StreamsClient: mockedDynamoStreamsClient{
			describeStreamOutput:      dynamodbstreams.DescribeStreamOutput{},
			describeStreamOutputError: errors.New("get stream description failed"),
		},
	}

	streamsDescriptions, err := c.getStreamsDescriptions(streams)
	assert.Error(t, err)
	assert.Equal(t, 0, len(streamsDescriptions))

	streams = []*dynamodbstreams.Stream{
		{
			StreamArn:   aws.String("foo"),
			StreamLabel: aws.String("bar"),
			TableName:   aws.String("fooTable"),
		},
	}

	c = Client{
		StreamsClient: mockedDynamoStreamsClient{
			describeStreamOutput: dynamodbstreams.DescribeStreamOutput{
				StreamDescription: &dynamodbstreams.StreamDescription{},
			},
			describeStreamOutputError: nil,
		},
	}

	streamsDescriptions, err = c.getStreamsDescriptions(streams)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(streamsDescriptions))
}

func TestGetShardIterators(t *testing.T) {

}

func TestGetLatestRecords(t *testing.T) {

}

func TestSendCloudevent(t *testing.T) {

}

/*
Copyright (c) 2019 TriggerMesh, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/cloudevents/sdk-go"
	"github.com/jarcoal/httpmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockedClientForCommits struct {
	codecommitiface.CodeCommitAPI
	GetBranchResp codecommit.GetBranchOutput
	GetCommitResp codecommit.GetCommitOutput
	GetBranchErr  error
	GetCommitErr  error
}

type mockedClientForPRs struct {
	codecommitiface.CodeCommitAPI
	ListPRsResp codecommit.ListPullRequestsOutput
	GetPRResp   codecommit.GetPullRequestOutput
	ListPRsErr  error
	GetPRErr    error
}

func (m mockedClientForCommits) GetBranch(in *codecommit.GetBranchInput) (*codecommit.GetBranchOutput, error) {
	return &m.GetBranchResp, m.GetBranchErr
}

func (m mockedClientForCommits) GetCommit(in *codecommit.GetCommitInput) (*codecommit.GetCommitOutput, error) {
	return &m.GetCommitResp, m.GetCommitErr
}

func (m mockedClientForPRs) ListPullRequests(in *codecommit.ListPullRequestsInput) (*codecommit.ListPullRequestsOutput, error) {
	return &m.ListPRsResp, m.ListPRsErr
}

func (m mockedClientForPRs) GetPullRequest(in *codecommit.GetPullRequestInput) (*codecommit.GetPullRequestOutput, error) {
	return &m.GetPRResp, m.GetPRErr
}

func TestSendPREvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://foo.com", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://bar.com", httpmock.NewStringResponder(500, ``))

	transport, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("https://bar.com"),
	)
	assert.NoError(t, err)

	cloudClient, err := cloudevents.NewClient(transport)
	assert.NoError(t, err)

	clients := Clients{
		CloudEvents: cloudClient,
	}

	pullRequest := codecommit.PullRequest{}
	pullRequest.SetPullRequestId("12345")

	err = clients.sendPREvent(&pullRequest)
	assert.Error(t, err)

	transport, err = cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("https://foo.com"),
	)
	assert.NoError(t, err)

	cloudClient, err = cloudevents.NewClient(transport)
	assert.NoError(t, err)

	clients.CloudEvents = cloudClient

	err = clients.sendPREvent(&pullRequest)
	assert.NoError(t, err)
}

func TestSendCommitEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://foo.com", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://bar.com", httpmock.NewStringResponder(500, ``))

	transport, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("https://bar.com"),
	)
	assert.NoError(t, err)

	cloudClient, err := cloudevents.NewClient(transport)
	assert.NoError(t, err)

	clients := Clients{
		CloudEvents: cloudClient,
	}

	commit := codecommit.Commit{}
	commit.SetCommitId("12345")

	err = clients.sendCommitEvent(&commit)
	assert.Error(t, err)

	transport, err = cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("https://foo.com"),
	)
	assert.NoError(t, err)

	cloudClient, err = cloudevents.NewClient(transport)
	assert.NoError(t, err)

	clients.CloudEvents = cloudClient

	err = clients.sendCommitEvent(&commit)
	assert.NoError(t, err)
}

func TestProcessCommits(t *testing.T) {
	lastCommit = "foo"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://foo.com", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://bar.com", httpmock.NewStringResponder(500, ``))

	testCases := []struct {
		GetBranchResp codecommit.GetBranchOutput
		GetCommitResp codecommit.GetCommitOutput
		GetBranchErr  error
		GetCommitErr  error
		Sink          string
		Err           error
	}{
		{
			GetBranchResp: codecommit.GetBranchOutput{},
			GetBranchErr:  errors.New("get branch failed"),
			Err:           errors.New("get branch failed"),
			Sink:          "https://foo.com",
		},
		{
			GetBranchResp: codecommit.GetBranchOutput{
				Branch: &codecommit.BranchInfo{CommitId: aws.String("123")},
			},
			GetCommitResp: codecommit.GetCommitOutput{},
			GetBranchErr:  nil,
			GetCommitErr:  errors.New("get commit failed"),
			Err:           errors.New("get commit failed"),
			Sink:          "https://foo.com",
		},
		{
			GetBranchResp: codecommit.GetBranchOutput{
				Branch: &codecommit.BranchInfo{CommitId: aws.String("123")},
			},
			GetCommitResp: codecommit.GetCommitOutput{Commit: &codecommit.Commit{CommitId: aws.String("foo")}},
			GetBranchErr:  nil,
			GetCommitErr:  nil,
			Err:           nil,
			Sink:          "https://foo.com",
		},
		{
			GetBranchResp: codecommit.GetBranchOutput{
				Branch: &codecommit.BranchInfo{CommitId: aws.String("123")},
			},
			GetCommitResp: codecommit.GetCommitOutput{Commit: &codecommit.Commit{CommitId: aws.String("bar")}},
			GetBranchErr:  nil,
			GetCommitErr:  nil,
			Sink:          "https://bar.com",
			Err:           errors.New("error sending cloudevent: 500"),
		},
		{
			GetBranchResp: codecommit.GetBranchOutput{
				Branch: &codecommit.BranchInfo{CommitId: aws.String("123")},
			},
			GetCommitResp: codecommit.GetCommitOutput{Commit: &codecommit.Commit{CommitId: aws.String("bar")}},
			GetBranchErr:  nil,
			GetCommitErr:  nil,
			Sink:          "https://foo.com",
			Err:           nil,
		},
	}

	for _, tt := range testCases {

		transport, err := cloudevents.NewHTTPTransport(
			cloudevents.WithTarget(tt.Sink),
		)
		assert.NoError(t, err)

		cloudClient, err := cloudevents.NewClient(transport)
		assert.NoError(t, err)

		clients := Clients{
			CodeCommit: mockedClientForCommits{
				GetBranchResp: tt.GetBranchResp,
				GetCommitResp: tt.GetCommitResp,
				GetBranchErr:  tt.GetBranchErr,
				GetCommitErr:  tt.GetCommitErr,
			},
			CloudEvents: cloudClient,
		}

		err = clients.processCommits()
		assert.Equal(t, tt.Err, err)
		lastCommit = "foo"

	}

}

func TestProcessPullRequest(t *testing.T) {
	pullRequestIDs = []*string{aws.String("1")}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://foo.com", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://bar.com", httpmock.NewStringResponder(500, ``))

	testCases := []struct {
		ListPRsResp codecommit.ListPullRequestsOutput
		GetPRResp   codecommit.GetPullRequestOutput
		ListPRsErr  error
		GetPRErr    error
		Sink        string
		Err         error
	}{
		{
			ListPRsResp: codecommit.ListPullRequestsOutput{},
			ListPRsErr:  errors.New("failed to list pull requests"),
			Sink:        "https://foo.com",
			Err:         errors.New("failed to list pull requests"),
		},
		{
			ListPRsResp: codecommit.ListPullRequestsOutput{PullRequestIds: []*string{aws.String("1")}},
			ListPRsErr:  nil,
			Sink:        "https://foo.com",
			Err:         nil,
		},
		{
			ListPRsResp: codecommit.ListPullRequestsOutput{PullRequestIds: []*string{aws.String("2")}},
			GetPRResp:   codecommit.GetPullRequestOutput{PullRequest: &codecommit.PullRequest{}},
			ListPRsErr:  nil,
			GetPRErr:    errors.New("failed to get pull request"),
			Sink:        "https://foo.com",
			Err:         errors.New("failed to get pull request"),
		},
		{
			ListPRsResp: codecommit.ListPullRequestsOutput{PullRequestIds: []*string{aws.String("2")}},
			GetPRResp:   codecommit.GetPullRequestOutput{PullRequest: &codecommit.PullRequest{}},
			ListPRsErr:  nil,
			GetPRErr:    nil,
			Sink:        "https://bar.com",
			Err:         nil,
		},
		{
			ListPRsResp: codecommit.ListPullRequestsOutput{PullRequestIds: []*string{aws.String("2")}},
			GetPRResp:   codecommit.GetPullRequestOutput{PullRequest: &codecommit.PullRequest{}},
			ListPRsErr:  nil,
			GetPRErr:    nil,
			Sink:        "https://foo.com",
			Err:         nil,
		},
	}

	for _, tt := range testCases {

		transport, err := cloudevents.NewHTTPTransport(
			cloudevents.WithTarget(tt.Sink),
		)
		assert.NoError(t, err)

		cloudClient, err := cloudevents.NewClient(transport)
		assert.NoError(t, err)

		client := Clients{
			CodeCommit: mockedClientForPRs{
				ListPRsResp: tt.ListPRsResp,
				GetPRResp:   tt.GetPRResp,
				ListPRsErr:  tt.ListPRsErr,
				GetPRErr:    tt.GetPRErr,
			},
			CloudEvents: cloudClient,
		}

		_, err = client.preparePullRequests()
		assert.Equal(t, tt.Err, err)

		pullRequestIDs = []*string{aws.String("1")}

	}
}

func TestRemoveOldPRs(t *testing.T) {
	oldPRs := []*codecommit.PullRequest{
		{PullRequestId: aws.String("1"), PullRequestStatus: aws.String("CREATED")},
		{PullRequestId: aws.String("2"), PullRequestStatus: aws.String("CREATED")},
	}
	newPRs := []*codecommit.PullRequest{
		{PullRequestId: aws.String("1"), PullRequestStatus: aws.String("CREATED")},
		{PullRequestId: aws.String("2"), PullRequestStatus: aws.String("CLOSED")},
		{PullRequestId: aws.String("3"), PullRequestStatus: aws.String("CREATED")},
	}

	expectedPRs := []*codecommit.PullRequest{
		{PullRequestId: aws.String("2"), PullRequestStatus: aws.String("CLOSED")},
		{PullRequestId: aws.String("3"), PullRequestStatus: aws.String("CREATED")},
	}

	prs := removeOldPRs(oldPRs, newPRs)
	log.Info(prs)
	assert.Equal(t, 2, len(prs))
	assert.Equal(t, expectedPRs, prs)

}

package line

import (
	"context"
	"log"

	talk "../talkservice"

	"../login"
)

type Client struct {
	talk     *talk.TalkServiceClient
	poll     *talk.TalkServiceClient
	Revision int64
	ctx      context.Context
	Profile  *talk.Profile
}

func NewClient(token, appName string) Client {
	ctx := context.Background()
	talk := login.Talk(token, appName)
	poll := login.Poll(token, appName)
	profile, err := talk.GetProfile(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(profile.DisplayName + ": login success")
	return Client{talk: talk, poll: poll, Revision: -1, ctx: ctx, Profile: profile}
}

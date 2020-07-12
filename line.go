package lineGo

import (
	talk "github.com/n4tsumi/lineGo/talkservice"
	"log"
	"net/http"
)

type Client struct {
	Talk *talk.TalkServiceClient
	Poll *talk.TalkServiceClient

	Profile *talk.Profile

	Revision    int64
	OpInterrupt map[talk.Af]func(*talk.Operation)
	hc          *http.Client
}

func NewClient() *Client {
	return &Client{
		Revision:    -1,
		OpInterrupt: make(map[talk.Af]func(*talk.Operation)),
		hc:          &http.Client{Transport: &http.Transport{}},
	}
}

// Login
func (c *Client) Login(l LineLogin, options ...LoginOption) error {
	switch l.Type() {
	case authToken:
		var err error
		c.Talk = createSession(l.Value(), Normal, c.hc)
		c.Poll = createSession(l.Value(), Polling, c.hc)
		c.Profile, err = c.Talk.GetProfile(ctx, talk.SyncReason_UNSPECIFIED)
		if err != nil {
			return err
		}
		log.Println(c.Profile.DisplayName + ": login success")

	default:
		log.Fatal("under development")
	}
	return nil
}

// Op Interrupt
func (c *Client) SetOpInterrupt(opInterrupt map[talk.Af]func(*talk.Operation)) {
	c.OpInterrupt = opInterrupt
}

func (c *Client) AddOpInterrupt(optype talk.Af, fn func(*talk.Operation)) {
	c.OpInterrupt[optype] = fn
}

func (c *Client) Run() {
	for {
		ops, err := c.Poll.FetchOperations(ctx, c.Revision, 100)
		if err != nil {
			log.Printf("%#v\n", err)
			_ = c.Poll.Noop(nil)
			c.hc.CloseIdleConnections()
			continue
		}
		for _, op := range ops {
			if op.Type != talk.Af_END_OF_OPERATION {
				if c.Revision < op.Revision {
					c.Revision = op.Revision
				}
				if fn, ok := c.OpInterrupt[op.Type]; ok {
					fn(op)
				}
			}
		}
	}
}

// User
func (c *Client) AcquireEncryptedAccessToken(featureType talk.Nc) (string, error) {
	return c.Talk.AcquireEncryptedAccessToken(ctx, featureType)
}

func (c *Client) GetProfile() (*talk.Profile, error) {
	return c.Talk.GetProfile(ctx, talk.SyncReason_UNSPECIFIED)
}

func (c *Client) GetSettings() (*talk.Settings, error) {
	return c.Talk.GetSettings(ctx, talk.SyncReason_UNSPECIFIED)
}

// Message
func (c *Client) SendMessage(msg *talk.Message) (*talk.Message, error) {
	return c.Talk.SendMessage(ctx, 0, msg)
}

func (c *Client) SendText(to, text string) (*talk.Message, error) {
	msg := &talk.Message{To: to, Text: text}
	return c.Talk.SendMessage(ctx, 0, msg)
}

func (c *Client) SendChatChecked(to, lastMessageId string, sessionId int8) error {
	return c.Talk.SendChatChecked(ctx, 0, to, lastMessageId, sessionId)
}

// Chat
func (c *Client) GetChatRoomAnnouncements(groupId string) ([]*talk.ChatRoomAnnouncement, error) {
	return c.Talk.GetChatRoomAnnouncements(ctx, groupId)
}

func (c *Client) CreateChatRoomAnnouncement(reqSeq int32, chatRoomMid string, typeA1 talk.Y9, contents *talk.ChatRoomAnnouncementContents) (r *talk.ChatRoomAnnouncement, err error) {
	return c.Talk.CreateChatRoomAnnouncement(ctx, reqSeq, chatRoomMid, typeA1, contents)
}

func (c *Client) AcceptChatInvitation(groupId string) (e error) {
	_, e = c.Talk.AcceptChatInvitation(ctx, &talk.AcceptChatInvitationRequest{ChatMid: groupId})
	return
}

func (c *Client) AcceptChatInvitationByTicket(groupId, ticketId string) (e error) {
	_, e = c.Talk.AcceptChatInvitationByTicket(ctx, &talk.AcceptChatInvitationByTicketRequest{ChatMid: groupId, TicketId: ticketId})
	return
}

func (c *Client) GetChat(groupId string) (*talk.Chat, error) {
	r, e := c.Talk.GetChats(ctx, &talk.GetChatsRequest{ChatMid: []string{groupId}, WithMembers: true, WithInvitees: true})
	if e != nil {
		return nil, e
	}
	return r.Chats[0], nil
}

func (c *Client) GetChats(groupIds []string) ([]*talk.Chat, error) {
	r, e := c.Talk.GetChats(ctx, &talk.GetChatsRequest{ChatMid: groupIds, WithMembers: true, WithInvitees: true})
	if e != nil {
		return nil, e
	}
	return r.Chats, nil
}

func (c *Client) UpdateChat(chat *talk.Chat, updateType talk.P9) (e error) {
	_, e = c.Talk.UpdateChat(ctx, &talk.UpdateChatRequest{ReqSeq: 0, Chat: chat, UpdatedAttribute: updateType})
	return
}

func (c *Client) ReissueChatTicket(groupMid string) (string, error) {
	r, e := c.Talk.ReissueChatTicket(ctx, &talk.ReissueChatTicketRequest{GroupMid: groupMid})
	if e != nil {
		return "", e
	}
	return r.TicketId, nil
}

func (c *Client) InviteIntoChat(groupId string, midlist []string) (e error) {
	_, e = c.Talk.InviteIntoChat(ctx, &talk.InviteIntoChatRequest{ChatMid: groupId, TargetUserMids: midlist})
	return
}

func (c *Client) LeaveChat(groupId string) (e error) {
	_, e = c.Talk.DeleteSelfFromChat(ctx, &talk.DeleteSelfFromChatRequest{ChatMid: groupId})
	return
}

func (c *Client) KickoutFromChat(groupId, mid string) (e error) {
	_, e = c.Talk.DeleteOtherFromChat(ctx, &talk.DeleteOtherFromChatRequest{ChatMid: groupId, TargetUserMids: []string{mid}})
	return
}

func (c *Client) CancelChatInvitation(groupId, mid string) (e error) {
	_, e = c.Talk.CancelChatInvitation(ctx, &talk.CancelChatInvitationRequest{ChatMid: groupId, TargetUserMids: []string{mid}})
	return
}

func (c *Client) RejectChatInvitation(groupId string) (e error) {
	_, e = c.Talk.RejectChatInvitation(ctx, &talk.RejectChatInvitationRequest{ChatMid: groupId})
	return
}

func (c *Client) FindChatByTicket(ticketId string) (*talk.Chat, error) {
	r, e := c.Talk.FindChatByTicket(ctx, &talk.FindChatByTicketRequest{TicketId: ticketId})
	if e != nil {
		return nil, e
	}
	return r.Chat, nil
}

func (c *Client) GetAllChatIds() ([]string, []string, error) {
	r, e := c.Talk.GetAllChatMids(ctx, &talk.GetAllChatMidsRequest{WithMemberChats: true, WithInvitedChats: true}, talk.SyncReason_UNSPECIFIED)
	if e != nil {
		return nil, nil, e
	}
	return r.MemberMids, r.InviteeMids, nil
}

// Contact
func (c *Client) GetContact(mid string) (*talk.Contact, error) {
	return c.Talk.GetContact(ctx, mid)
}

func (c *Client) GetContacts(midlist []string) ([]*talk.Contact, error) {
	return c.Talk.GetContacts(ctx, midlist)
}

func (c *Client) GetAllContactIds() ([]string, error) {
	return c.Talk.GetAllContactIds(ctx, talk.SyncReason_UNSPECIFIED)
}

func (c *Client) FindAndAddContactsByMid(mid string) (map[string]*talk.Contact, error) {
	return c.Talk.FindAndAddContactsByMid(ctx, 0, mid, talk.Gb_MID, "")
}

func (c *Client) BlockContact(mid string) error {
	return c.Talk.BlockContact(ctx, 0, mid)
}

func (c *Client) UnblockContact(mid string) error {
	return c.Talk.UnblockContact(ctx, 0, mid, "")
}

func (c *Client) BlockRecommendation(mid string) error {
	return c.Talk.BlockRecommendation(ctx, 0, mid)
}

func (c *Client) UnblockRecommendation(mid string) error {
	return c.Talk.UnblockRecommendation(ctx, 0, mid)
}

func (c *Client) GetBlockedContactIds() ([]string, error) {
	return c.Talk.GetBlockedContactIds(ctx, talk.SyncReason_UNSPECIFIED)
}

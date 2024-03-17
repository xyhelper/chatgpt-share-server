package service

import (
	"backend/modules/chatgpt/model"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type ChatgptConversationsService struct {
	*cool.Service
}

func NewChatgptConversationsService() *ChatgptConversationsService {
	return &ChatgptConversationsService{
		&cool.Service{
			Model: model.NewChatgptConversations(),
			PageQueryOp: &cool.QueryOp{
				FieldEQ:      []string{"convid", "title", "chatgptaccountid", "email", "usertoken"},
				KeyWordField: []string{"convid", "title", "chatgptaccountid", "email", "usertoken"},
			},
		},
	}
}

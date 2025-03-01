package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	pb "github.com/micro/distributed-api/threads/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/gorm"
)

// Create a message within a conversation
func (s *Threads) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
	// validate the request
	if len(req.AuthorId) == 0 {
		return ErrMissingAuthorID
	}
	if len(req.ConversationId) == 0 {
		return ErrMissingConversationID
	}
	if len(req.Text) == 0 {
		return ErrMissingText
	}

	// lookup the conversation
	var conv Conversation
	if err := s.DB.Where(&Conversation{ID: req.ConversationId}).First(&conv).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// create the message
	msg := &Message{
		ID:             req.Id,
		SentAt:         s.Time(),
		Text:           req.Text,
		AuthorID:       req.AuthorId,
		ConversationID: req.ConversationId,
	}
	if len(msg.ID) == 0 {
		msg.ID = uuid.New().String()
	}
	if err := s.DB.Create(msg).Error; err == nil {
		rsp.Message = msg.Serialize()
		return nil
	} else if !strings.Contains(err.Error(), "messages_pkey") {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// a message already exists with this id
	var existing Message
	if err := s.DB.Where(&Message{ID: msg.ID}).First(&existing).Error; err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}
	rsp.Message = existing.Serialize()
	return nil
}

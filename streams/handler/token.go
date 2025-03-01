package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/micro/distributed-api/streams/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

func (s *Streams) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	// construct the token and write it to the database
	t := Token{
		Token:     uuid.New().String(),
		ExpiresAt: s.Time().Add(TokenTTL),
		Topic:     req.Topic,
	}
	if err := s.DB.Create(&t).Error; err != nil {
		logger.Errorf("Error creating token in store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error writing token to database")
	}

	// return the token in the response
	rsp.Token = t.Token
	return nil
}

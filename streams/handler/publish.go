package handler

import (
	"context"

	pb "github.com/micro/distributed-api/streams/proto"
	"github.com/micro/micro/v3/service/logger"
)

func (s *Streams) Publish(ctx context.Context, req *pb.Message, rsp *pb.PublishResponse) error {
	// validate the request
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}
	if len(req.Message) == 0 {
		return ErrMissingMessage
	}

	// publish the message
	logger.Infof("Publishing message to topic: %v", req.Topic)
	return s.Events.Publish(req.Topic, req.Message)
}

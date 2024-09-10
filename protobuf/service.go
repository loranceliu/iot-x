package protobuf

import (
	"context"
)

type Service struct {
}

func (s *Service) mustEmbedUnimplementedInstanceServer() {
}

func (s *Service) Transfer(ctx context.Context, input *MessageInput) (*MessageOutput, error) {
	return &MessageOutput{
		RequestId: "1",
		Type:      0,
		Msg:       "2",
		Data:      nil,
	}, nil
}

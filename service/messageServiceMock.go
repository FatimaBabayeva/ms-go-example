package service

import (
	"context"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/stretchr/testify/mock"
)

type MessageServiceMock struct {
	mock.Mock
}

func (s *MessageServiceMock) SaveMessage(ctx context.Context, message model.Message) (*model.Message, error) {
	// args hold the arguments that should be returned when this method is called.
	args := s.Called(ctx, message)
	return checkArguments(args)
}

func (s *MessageServiceMock) GetMessageById(ctx context.Context, id int64) (*model.Message, error) {
	args := s.Called(ctx, id)
	return checkArguments(args)
}

func (s *MessageServiceMock) UpdateMessageById(ctx context.Context, id int64, message model.Message) (*model.Message, error) {
	args := s.Called(ctx, id, message)
	return checkArguments(args)
}

func (s *MessageServiceMock) DeleteMessageById(ctx context.Context, id int64) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func checkArguments(args mock.Arguments) (*model.Message, error) {
	firstArg := args.Get(0)
	if firstArg != nil {
		return firstArg.(*model.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

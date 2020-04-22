package service

import (
	"context"
	"github.com/FatimaBabayeva/ms-go-example/ctmerror"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/FatimaBabayeva/ms-go-example/repo"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

// MessageService is an interface to operate with messages
type MessageService interface {
	SaveMessage(ctx context.Context, message model.Message) (*model.Message, error)
	GetMessageById(ctx context.Context, id int64) (*model.Message, error)
	UpdateMessageById(ctx context.Context, id int64, message model.Message) (*model.Message, error)
	DeleteMessageById(ctx context.Context, id int64) error
}

// MessageServiceImpl is an implementation of MessageService
type MessageServiceImpl struct {
	MsgRepo repo.MessageRepo
}

func (s *MessageServiceImpl) SaveMessage(ctx context.Context, message model.Message) (*model.Message, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.SaveMessage.start")

	message.Id = 0
	message.Status = model.CREATED
	result, err := s.MsgRepo.Save(&message)
	if err != nil {
		logger.Errorf("ActionLog.SaveMessage.error : Error saving message %v,\n%s", err, string(debug.Stack()))
		return nil, ctmerror.NewMessageError(err)
	}

	logger.Info("ActionLog.SaveMessage.end")
	return result, nil
}

func (s *MessageServiceImpl) GetMessageById(ctx context.Context, id int64) (*model.Message, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetMessageById.start")

	result, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.GetMessageById.error : Error getting message with id = %d, %v,\n%s", id, err, string(debug.Stack()))
		return nil, ctmerror.NewMessageError(err)
	}

	logger.Info("ActionLog.GetMessageById.end")
	return result, nil
}

func (s *MessageServiceImpl) UpdateMessageById(ctx context.Context, id int64, message model.Message) (*model.Message, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.UpdateMessageById.start")

	originalMsg, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.UpdateMessageById.error : Error getting message with id = %d, %v,\n%s", id, err, string(debug.Stack()))
		return nil, ctmerror.NewMessageError(err)
	}

	if message.Text != "" {
		originalMsg.Text = message.Text
		originalMsg.UpdatedAt = time.Now()
	}

	result, err := s.MsgRepo.Update(originalMsg)
	if err != nil {
		logger.Errorf("ActionLog.UpdateMessageById.error : Error updating message with id = %d, %v,\n%s", id, err, string(debug.Stack()))
		return nil, ctmerror.NewMessageError(err)
	}

	logger.Info("ActionLog.UpdateMessageById.end")
	return result, nil
}

func (s *MessageServiceImpl) DeleteMessageById(ctx context.Context, id int64) error {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.DeleteMessageById.start")

	originalMsg, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.DeleteMessageById.error : Error getting message with id = %d, %v,\n%s", id, err, string(debug.Stack()))
		return ctmerror.NewMessageError(err)
	}

	originalMsg.UpdatedAt = time.Now()
	originalMsg.Status = model.DELETED
	_, err = s.MsgRepo.Update(originalMsg)
	if err != nil {
		logger.Errorf("ActionLog.DeleteMessageById.error : Error deleting message with id = %d, %v,\n%s", id, err, string(debug.Stack()))
		return ctmerror.NewMessageError(err)
	}

	logger.Info("ActionLog.DeleteMessageById.end")
	return nil
}

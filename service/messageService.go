package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"ms-go-example/model"
	"ms-go-example/repo"
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
		logger.Errorf("ActionLog.SaveMessage.error : Error saving message %v", err)
		return nil, err
	}

	logger.Info("ActionLog.SaveMessage.end")
	return result, nil
}

func (s *MessageServiceImpl) GetMessageById(ctx context.Context, id int64) (*model.Message, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetMessageById.start")

	result, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.GetMessageById.error : Error getting message with id = %d, %v", id, err)
		return nil, err
	}

	logger.Info("ActionLog.GetMessageById.end")
	return result, nil
}

func (s *MessageServiceImpl) UpdateMessageById(ctx context.Context, id int64, message model.Message) (*model.Message, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.UpdateMessageById.start")

	originalMsg, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.UpdateMessageById.error : Error getting message with id = %d, %v", id, err)
		return nil, err
	}

	if message.Text != "" {
		originalMsg.Text = message.Text
		originalMsg.UpdatedAt = time.Now()
	}

	result, err := s.MsgRepo.Update(originalMsg)
	if err != nil {
		logger.Errorf("ActionLog.UpdateMessageById.error : Error updating message with id = %d, %v", id, err)
		return nil, err
	}

	logger.Info("ActionLog.UpdateMessageById.end")
	return result, nil
}

func (s *MessageServiceImpl) DeleteMessageById(ctx context.Context, id int64) error {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.DeleteMessageById.start")

	originalMsg, err := s.MsgRepo.Get(id)
	if err != nil {
		logger.Errorf("ActionLog.DeleteMessageById.error : Error getting message with id = %d, %v", id, err)
		return err
	}

	originalMsg.UpdatedAt = time.Now()
	originalMsg.Status = model.DELETED
	_, err = s.MsgRepo.Update(originalMsg)
	if err != nil {
		logger.Errorf("ActionLog.DeleteMessageById.error : Error deleting message with id = %d, %v", id, err)
		return err
	}

	logger.Info("ActionLog.DeleteMessageById.end")
	return nil
}

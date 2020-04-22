package service

import (
	"context"
	"github.com/FatimaBabayeva/ms-go-example/ctmerror"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/FatimaBabayeva/ms-go-example/repo"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	mockRepo      = repo.MessageRepoMock{}
	s             = MessageServiceImpl{MsgRepo: &mockRepo}
	unexpectedErr = ctmerror.NewMessageErrorBuilder("error.go-example.unexpected-error", assert.AnError, 500)
	notFoundErr   = ctmerror.NewMessageErrorBuilder("error.go-example.message-not-found", pg.ErrNoRows, 404)

	id int64 = 1

	errorTable = []struct {
		repoError error
		msgError  *ctmerror.MessageError
	}{
		{assert.AnError, unexpectedErr},
		{pg.ErrNoRows, notFoundErr},
	}
)

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.ContextLogger, log.WithContext(ctx))
	return ctx
}

func TestMessageServiceImpl_SaveMessage_Ok(t *testing.T) {
	// given:
	message := model.Message{
		Id:     0,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	mockRepo.On("Save", &message).Once().Return(&message, nil)

	// when:
	result, err := s.SaveMessage(mockContext(), message)

	// then:
	assert.Nil(t, err)
	assert.NotEqual(t, 0, result.Id)
	assert.Equal(t, message.Text, result.Text)
	assert.Equal(t, message.Status, result.Status)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_GetMessageById_Ok(t *testing.T) {
	// given:
	message := model.Message{Id: id}
	mockRepo.On("Get", id).Once().Return(&message, nil)

	// when:
	result, err := s.GetMessageById(mockContext(), id)

	// then:
	assert.Nil(t, err)
	assert.Equal(t, message.Id, result.Id)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_UpdateMessageById_Ok(t *testing.T) {
	// given:
	message := model.Message{Text: "UPDATED_TEXT"}

	originalMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	updatedMessage := model.Message{
		Id:     id,
		Text:   "UPDATED_TEXT",
		Status: "CREATED",
	}

	mockRepo.On("Get", id).Once().Return(&originalMessage, nil)
	mockRepo.On("Update", mock.MatchedBy(func(msg *model.Message) bool {
		return msg.Id == id &&
			msg.Text == message.Text &&
			msg.Status == originalMessage.Status
	})).Once().Return(&updatedMessage, nil)

	// when:
	result, err := s.UpdateMessageById(mockContext(), id, message)

	// then:
	assert.Nil(t, err)
	assert.Equal(t, id, result.Id)
	assert.Equal(t, message.Text, result.Text)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_DeleteMessageById_Ok(t *testing.T) {
	// given:
	originalMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	deletedMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "DELETED",
	}

	mockRepo.On("Get", id).Once().Return(&originalMessage, nil)
	mockRepo.On("Update", mock.MatchedBy(func(msg *model.Message) bool {
		return msg.Id == id &&
			msg.Text == originalMessage.Text &&
			msg.Status == deletedMessage.Status
	})).Once().Return(&deletedMessage, nil)

	// when:
	err := s.DeleteMessageById(mockContext(), id)

	// then:
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_SaveMessage_Error(t *testing.T) {
	// given:
	message := model.Message{}
	mockRepo.On("Save", mock.Anything).Once().Return(nil, assert.AnError)

	// when:
	result, err := s.SaveMessage(mockContext(), message)

	// then:
	assert.NotNil(t, err)
	assert.Equal(t, err, unexpectedErr)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_GetMessageById_Error(t *testing.T) {
	for _, errCase := range errorTable {
		// given:
		mockRepo.On("Get", id).Once().Return(nil, errCase.repoError)

		// when:
		result, err := s.GetMessageById(mockContext(), id)

		// then:
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err, errCase.msgError)
		mockRepo.AssertExpectations(t)
	}
}

func TestMessageServiceImpl_UpdateMessageById_MessageNotFound(t *testing.T) {
	// given:
	message := model.Message{Text: "UPDATED_TEXT"}
	mockRepo.On("Get", id).Once().Return(nil, pg.ErrNoRows)

	// when:
	result, err := s.UpdateMessageById(mockContext(), id, message)

	// then:
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err, notFoundErr)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_UpdateMessageById_AnyError(t *testing.T) {
	// given:
	message := model.Message{Text: "UPDATED_TEXT"}

	originalMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	mockRepo.On("Get", id).Once().Return(&originalMessage, nil)
	mockRepo.On("Update", mock.Anything).Once().Return(nil, assert.AnError)

	// when:
	result, err := s.UpdateMessageById(mockContext(), id, message)

	// then:
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err, unexpectedErr)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_DeleteMessageById_MessageNotFound(t *testing.T) {
	// given:
	mockRepo.On("Get", id).Once().Return(nil, pg.ErrNoRows)

	// when:
	err := s.DeleteMessageById(mockContext(), id)

	// then:
	assert.NotNil(t, err)
	assert.Equal(t, err, notFoundErr)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_DeleteMessageById_AnyError(t *testing.T) {
	// given:
	originalMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	mockRepo.On("Get", id).Once().Return(&originalMessage, nil)
	mockRepo.On("Update", mock.Anything).Once().Return(nil, assert.AnError)

	// when:
	err := s.DeleteMessageById(mockContext(), id)

	// then:
	assert.NotNil(t, err)
	assert.Equal(t, err, unexpectedErr)
	mockRepo.AssertExpectations(t)
}

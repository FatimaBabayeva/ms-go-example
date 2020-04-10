package service

import (
	"context"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type messageRepoMock struct {
	mock.Mock
}

func (r *messageRepoMock) Save(m *model.Message) (*model.Message, error) {
	args := r.Called(m)
	return checkArguments(args)
}

func (r *messageRepoMock) Update(m *model.Message) (*model.Message, error) {
	args := r.Called(m)
	return checkArguments(args)
}

func (r *messageRepoMock) Get(id int64) (*model.Message, error) {
	args := r.Called(id)
	return checkArguments(args)
}

func checkArguments(args mock.Arguments) (*model.Message, error) {
	firstArg := args.Get(0)
	if firstArg != nil {
		return firstArg.(*model.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.ContextLogger, log.WithContext(ctx))
	return ctx
}

var (
	mockRepo        = messageRepoMock{}
	s               = MessageServiceImpl{MsgRepo: &mockRepo}
	id        int64 = 1
)

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

func TestMessageServiceImpl_SaveMessage_AnyError(t *testing.T) {
	// given:
	message := model.Message{}
	mockRepo.On("Save", mock.Anything).Once().Return(nil, assert.AnError)

	// when:
	result, err := s.SaveMessage(mockContext(), message)

	// then:
	assert.NotNil(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_GetMessageById_AnyError(t *testing.T) {
	// given:
	mockRepo.On("Get", id).Once().Return(nil, assert.AnError)

	// when:
	result, err := s.GetMessageById(mockContext(), id)

	// then:
	assert.Nil(t, result)
	assert.NotNil(t, err)
	mockRepo.AssertExpectations(t)
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
	assert.Equal(t, err, pg.ErrNoRows)
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
	mockRepo.AssertExpectations(t)
}

func TestMessageServiceImpl_DeleteMessageById_MessageNotFound(t *testing.T) {
	// given:
	mockRepo.On("Get", id).Once().Return(nil, pg.ErrNoRows)

	// when:
	err := s.DeleteMessageById(mockContext(), id)

	// then:
	assert.NotNil(t, err)
	assert.Equal(t, err, pg.ErrNoRows)
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
	mockRepo.AssertExpectations(t)
}

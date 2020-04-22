package repo

import (
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/stretchr/testify/mock"
)

type MessageRepoMock struct {
	mock.Mock
}

func (r *MessageRepoMock) Save(m *model.Message) (*model.Message, error) {
	args := r.Called(m)
	return checkArguments(args)
}

func (r *MessageRepoMock) Update(m *model.Message) (*model.Message, error) {
	args := r.Called(m)
	return checkArguments(args)
}

func (r *MessageRepoMock) Get(id int64) (*model.Message, error) {
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

package repo

import (
	"github.com/FatimaBabayeva/ms-go-example/model"
)

// MessageRepo is an interface to operate with messages on Db level
type MessageRepo interface {
	Save(m *model.Message) (*model.Message, error)
	Update(m *model.Message) (*model.Message, error)
	Get(id int64) (*model.Message, error)
}

// MessageRepoImpl is an implementation of MessageRepo
type MessageRepoImpl struct {
}

func (r *MessageRepoImpl) Save(m *model.Message) (*model.Message, error) {
	_, err := Db.Model(m).Insert()
	return m, err
}

func (r *MessageRepoImpl) Update(m *model.Message) (*model.Message, error) {
	_, err := Db.Model(m).WherePK().Update()
	return m, err
}

func (r *MessageRepoImpl) Get(id int64) (*model.Message, error) {
	res := model.Message{Id: id}
	err := Db.Model(&res).WherePK().Select()
	return &res, err
}
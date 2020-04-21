package handler

import (
	"bytes"
	"encoding/json"
	"github.com/FatimaBabayeva/ms-go-example/ctmerror"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/FatimaBabayeva/ms-go-example/properties"
	"github.com/FatimaBabayeva/ms-go-example/service"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mockService   = service.MessageServiceMock{}
	handler       = messageHandler{&mockService}
	unexpectedErr = ctmerror.NewMessageErrorBuilder("error.go-example.unexpected-error", assert.AnError, 500)
	notFoundErr   = ctmerror.NewMessageErrorBuilder("error.go-example.message-not-found", pg.ErrNoRows, 404)

	id int64 = 1

	errorTable = []struct {
		msgError  *ctmerror.MessageError
		errorCode string
		httpCode  int
	}{
		{unexpectedErr, unexpectedErr.Error(), unexpectedErr.HttpCode()},
		{notFoundErr, notFoundErr.Error(), notFoundErr.HttpCode()},
	}
)

func TestSaveMessage_Ok(t *testing.T) {
	// given:
	message := model.Message{Text: "MOCK_TEXT"}
	savedMessage := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	mockService.On("SaveMessage", mock.Anything, message).Once().Return(&savedMessage, nil)

	requestJson, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", properties.RootPath+"/message", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.saveMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	result := model.Message{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, savedMessage, result)
	mockService.AssertExpectations(t)
}

func TestGetMessage_Ok(t *testing.T) {
	// given:
	message := model.Message{
		Id:     id,
		Text:   "MOCK_TEXT",
		Status: "CREATED",
	}
	mockService.On("GetMessageById", mock.Anything, id).Once().Return(&message, nil)

	req, err := http.NewRequest("GET", properties.RootPath+"/message/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	// when:
	handler := http.HandlerFunc(handler.getMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	result := model.Message{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, message, result)
	mockService.AssertExpectations(t)
}

func TestEditMessage_Ok(t *testing.T) {
	// given:
	message := model.Message{Text: "UPDATED_TEXT"}
	updatedMessage := model.Message{
		Id:     id,
		Text:   "UPDATED",
		Status: "CREATED",
	}
	mockService.On("UpdateMessageById", mock.Anything, id, message).Once().Return(&updatedMessage, nil)

	requestJson, _ := json.Marshal(message)
	req, err := http.NewRequest("PUT", properties.RootPath+"/message/{id}", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	// when:
	handler := http.HandlerFunc(handler.editMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	result := model.Message{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, updatedMessage, result)
	mockService.AssertExpectations(t)
}

func TestDeleteMessage_Ok(t *testing.T) {
	// given:
	mockService.On("DeleteMessageById", mock.Anything, id).Once().Return(nil)

	req, err := http.NewRequest("DELETE", properties.RootPath+"/message/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	// when:
	handler := http.HandlerFunc(handler.deleteMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	mockService.AssertExpectations(t)
}

func TestSaveMessage_InvalidBody(t *testing.T) {
	// given:
	requestJson, _ := json.Marshal("INVALID_DATA")
	req, err := http.NewRequest("POST", properties.RootPath+"/message", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.saveMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSaveMessage_ServiceError(t *testing.T) {
	// given:
	message := model.Message{Text: "MOCK_TEXT"}
	mockService.On("SaveMessage", mock.Anything, message).Once().Return(nil, unexpectedErr)

	requestJson, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", properties.RootPath+"/message", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.saveMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, "error.go-example.unexpected-error", strings.TrimSpace(w.Body.String()))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetMessage_NoParams(t *testing.T) {
	// given:
	req, err := http.NewRequest("GET", properties.RootPath+"/message/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.getMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditMessage_NoParams(t *testing.T) {
	// given:
	message := model.Message{Text: "UPDATED_TEXT"}
	requestJson, _ := json.Marshal(message)
	req, err := http.NewRequest("PUT", properties.RootPath+"/message/{id}", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.editMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditMessage_InvalidBody(t *testing.T) {
	// given:
	requestJson, _ := json.Marshal("INVALID_DATA")
	req, err := http.NewRequest("PUT", properties.RootPath+"/message/{id}", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	// when:
	handler := http.HandlerFunc(handler.editMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditMessage_ServiceError(t *testing.T) {
	for _, errCase := range errorTable {
		// given:
		message := model.Message{Text: "UPDATED_TEXT"}
		mockService.On("UpdateMessageById", mock.Anything, id, message).Once().Return(nil, errCase.msgError)

		requestJson, _ := json.Marshal(message)
		req, err := http.NewRequest("PUT", properties.RootPath+"/message/{id}", bytes.NewBuffer(requestJson))
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})

		// when:
		handler := http.HandlerFunc(handler.editMessage)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// then:
		assert.Equal(t, errCase.httpCode, w.Code)
		assert.Equal(t, errCase.errorCode, strings.TrimSpace(w.Body.String()))
		mockService.AssertExpectations(t)
	}
}

func TestDeleteMessage_NoParams(t *testing.T) {
	// given:
	req, err := http.NewRequest("DELETE", properties.RootPath+"/message/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}

	// when:
	handler := http.HandlerFunc(handler.deleteMessage)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then:
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteMessage_ServiceError(t *testing.T) {
	for _, errCase := range errorTable {
		// given:
		mockService.On("DeleteMessageById", mock.Anything, id).Once().Return(errCase.msgError)

		req, err := http.NewRequest("DELETE", properties.RootPath+"/message/{id}", nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})

		// when:
		handler := http.HandlerFunc(handler.deleteMessage)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// then:
		assert.Equal(t, errCase.httpCode, w.Code)
		assert.Equal(t, errCase.errorCode, strings.TrimSpace(w.Body.String()))
		mockService.AssertExpectations(t)
	}
}

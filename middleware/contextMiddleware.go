package middleware

import (
	"context"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var headers = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
	"User-Agent",
	"X-Forwarded-For",
	"requestid",
}

// RequestParamsMiddleware is middleware function for context time logger and header transport
func RequestParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		requestID := r.Header.Get(model.HeaderKeyRequestID)
		operation := r.RequestURI
		userAgent := r.Header.Get(model.HeaderKeyUserAgent)
		userIP := r.Header.Get(model.HeaderKeyUserIP)

		if len(requestID) == 0 {
			requestID = uuid.New().String()
		}
		fields := log.Fields{}
		addLoggerParam(fields, model.LoggerKeyRequestID, requestID)
		addLoggerParam(fields, model.LoggerKeyOperation, operation)
		addLoggerParam(fields, model.LoggerKeyUserAgent, userAgent)
		addLoggerParam(fields, model.LoggerKeyUserIP, userIP)

		logger := log.WithFields(fields)
		header := http.Header{}

		for _, v := range headers {
			header.Add(v, r.Header.Get(v))
		}

		ctx = context.WithValue(ctx, model.ContextLogger, logger)
		ctx = context.WithValue(ctx, model.ContextHeader, header)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addLoggerParam(fields log.Fields, field string, value string) {
	if len(value) > 0 {
		fields[field] = value
	}
}

package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	shortenerMock "github.com/0xc00000f/shortener-tpl/internal/shortener/mocks"
)

func TestSaveURL_UserNil_Positive(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	baseURL := "http://example.com"
	expectedShort := "5ZytxbC"
	long := "https://dzen.ru/"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	encoder.EXPECT().Short(uuid.Nil, long).Return(expectedShort, nil)

	serverFunc := handlers.SaveURL(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(long),
	)
	req.Header.Set("content-type", "text/html")

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("url", "")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("content-type"))
	assert.Equal(t, fmt.Sprintf("%s/%s", baseURL, expectedShort), string(result))
}

func TestSaveURLJson_UserNil_Positive(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	baseURL := "http://example.com"
	expectedShort := "5ZytxbC"
	long := "https://dzen.ru/"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	encoder.EXPECT().Short(uuid.Nil, long).Return(expectedShort, nil)

	serverFunc := handlers.SaveURLJson(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, long)),
	)

	req.Header.Set("content-type", "application/json")

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("content-type"))
	assert.Equal(t, fmt.Sprintf(`{"result":"%s/%s"}`, baseURL, expectedShort), string(result))
}

package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	shortenerMock "github.com/0xc00000f/shortener-tpl/internal/shortener/mocks"
)

func TestRedirect_Positive(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	baseURL := "http://example.com"
	short := "5ZytxbC"
	expectedLong := "https://dzen.ru/"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	encoder.EXPECT().Get(short).Return(expectedLong, nil)

	serverFunc := handlers.Redirect(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/%s", short),
		nil,
	)
	req.Header.Set("content-type", "application/json")

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("url", short)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("content-type"))
	assert.Equal(t, expectedLong, res.Header.Get("Location"))
	assert.Empty(t, result)
}

func TestRedirect_EncoderGetError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	baseURL := "http://example.com"
	short := "5ZytxbC"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	getErr := errors.New("db is down")
	encoder.EXPECT().Get(short).Return("", getErr)

	serverFunc := handlers.Redirect(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/%s", short),
		nil,
	)
	req.Header.Set("content-type", "application/json")

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("url", short)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("content-type"))
}

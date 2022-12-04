package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	shortenerMock "github.com/0xc00000f/shortener-tpl/internal/shortener/mocks"
	"github.com/0xc00000f/shortener-tpl/internal/user"
)

var errStorageOutOfReach = errors.New("db is down")

func TestGetSavedData_Positive_201(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	const baseURL = "http://example.com"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	exp := map[string]string{
		"5ZytxbC": "https://dzen.ru/",
	}
	encoder.EXPECT().GetAll(user.Nil.UserID).Return(exp, nil)

	prepareMap := []handlers.Result{
		{
			Short: fmt.Sprintf("%s/%s", baseURL, "5ZytxbC"),
			Long:  "https://dzen.ru/",
		},
	}
	expectedResult, err := json.MarshalIndent(prepareMap, "", " ")
	require.NoError(t, err)

	serverFunc := handlers.GetSavedData(ns).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		"/user/urls",
		nil,
	)
	req.Header.Set("content-type", "application/json")

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("content-type"))

	require.JSONEq(t, string(expectedResult), string(result))
}

func TestGetSavedData_Positive_204(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	const baseURL = "http://example.com"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	exp := map[string]string{}
	encoder.EXPECT().GetAll(user.Nil.UserID).Return(exp, nil)

	serverFunc := handlers.GetSavedData(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/user/urls",
		nil,
	)
	req.Header.Set("content-type", "application/json")
	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Empty(t, result)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("content-type"))
}

func TestGetSavedData_Negative_GetAllError(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	const baseURL = "http://example.com"

	encoder := shortenerMock.NewMockShortener(ctl)
	ns := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(baseURL),
		shortener.SetLogger(zap.L()),
	)

	encoder.EXPECT().GetAll(user.Nil.UserID).Return(nil, errStorageOutOfReach)

	serverFunc := handlers.GetSavedData(ns).ServeHTTP
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/user/urls",
		nil,
	)
	req.Header.Set("content-type", "application/json")
	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

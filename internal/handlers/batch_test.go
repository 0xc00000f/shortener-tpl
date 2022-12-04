package handlers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestBatch_UserNil_Positive(t *testing.T) {
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

	prepareMap := []handlers.IOBatch{
		{
			CorrelationID: "1",
			ShortURL:      fmt.Sprintf("%s/%s", ns.BaseURL, "5ZytxbC"),
		},
		{
			CorrelationID: "2",
			ShortURL:      fmt.Sprintf("%s/%s", ns.BaseURL, "RmOSY54"),
		},
	}
	expectedResult, err := json.MarshalIndent(prepareMap, "", " ")
	require.NoError(t, err)

	serverFunc := handlers.Batch(ns).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		strings.NewReader(
			`[
			{
				"correlation_id": "1",
				"original_url": "https://dzen.ru/"
			},
			{
				"correlation_id": "2",
				"original_url": "https://ya.ru/"
			}
		]
		`),
	)
	req.Header.Set("content-type", "application/json")

	first := encoder.EXPECT().Short(
		req.Context(),
		user.Nil.UserID,
		"https://dzen.ru/",
	).Return("5ZytxbC", nil)
	second := encoder.EXPECT().Short(
		req.Context(),
		user.Nil.UserID,
		"https://ya.ru/",
	).Return("RmOSY54", nil)

	gomock.InOrder(
		first,
		second,
	)

	serverFunc(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("content-type"))

	require.JSONEq(t, string(expectedResult), string(result))
}

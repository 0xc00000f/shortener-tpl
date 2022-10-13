package handlers

//import (
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//
//	"github.com/0xc00000f/shortener-tpl/internal/utils"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestSaveURL(t *testing.T) {
//	r := NewRouter()
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//
//	statusCode, body := testRequest(t, ts, "POST", "/", strings.NewReader("https://vk.com"), true)
//	assert.Equal(t, http.StatusCreated, statusCode)
//	assert.True(t, utils.IsURL(body))
//}

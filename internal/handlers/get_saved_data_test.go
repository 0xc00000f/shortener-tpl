package handlers

//func TestGetSavedData_Positive_201(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	baseURL := "http://example.com"
//
//	encoder := shortenerMock.NewMockShortener(ctl)
//	ns := shortener.New(
//		shortener.SetEncoder(encoder),
//		shortener.InitBaseURL(baseURL),
//		shortener.SetLogger(zap.L()),
//	)
//
//	exp := map[string]string{
//		"5ZytxbC": "https://dzen.ru/",
//		"RmOSY54": "https://ya.ru/",
//	}
//	encoder.EXPECT().GetAll().Return(exp, nil)
//	prepareMap := []result{
//		{
//			Short: fmt.Sprintf("%s/%s", baseURL, "5ZytxbC"),
//			Long:  "https://dzen.ru/",
//		},
//		{
//			Short: fmt.Sprintf("%s/%s", baseURL, "RmOSY54"),
//			Long:  "https://ya.ru/",
//		},
//	}
//	expectedResult, err := json.MarshalIndent(prepareMap, "", " ")
//	require.NoError(t, err)
//
//	serverFunc := GetSavedData(ns).ServeHTTP
//
//	rec := httptest.NewRecorder()
//
//	req := httptest.NewRequest(
//		http.MethodGet,
//		"/user/urls",
//		nil,
//	)
//	req.Header.Set("content-type", "application/json")
//
//	serverFunc(rec, req)
//
//	res := rec.Result()
//	defer res.Body.Close()
//
//	result, err := io.ReadAll(res.Body)
//	require.NoError(t, err)
//
//	assert.Equal(t, http.StatusCreated, res.StatusCode)
//	assert.Equal(t, "application/json", res.Header.Get("content-type"))
//
//	require.JSONEq(t, string(expectedResult), string(result))
//}
//
//func TestGetSavedData_Positive_204(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	baseURL := "http://example.com"
//
//	encoder := shortenerMock.NewMockShortener(ctl)
//	ns := shortener.New(
//		shortener.SetEncoder(encoder),
//		shortener.InitBaseURL(baseURL),
//		shortener.SetLogger(zap.L()),
//	)
//
//	exp := map[string]string{}
//	encoder.EXPECT().GetAll().Return(exp, nil)
//
//	serverFunc := GetSavedData(ns).ServeHTTP
//	rec := httptest.NewRecorder()
//	req := httptest.NewRequest(
//		http.MethodGet,
//		"/user/urls",
//		nil,
//	)
//	req.Header.Set("content-type", "application/json")
//	serverFunc(rec, req)
//
//	res := rec.Result()
//	defer res.Body.Close()
//
//	result, err := io.ReadAll(res.Body)
//	require.NoError(t, err)
//
//	assert.Empty(t, result)
//	assert.Equal(t, http.StatusNoContent, res.StatusCode)
//	assert.Equal(t, "application/json", res.Header.Get("content-type"))
//}
//
//func TestGetSavedData_Negative_GetAllError(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	baseURL := "http://example.com"
//
//	encoder := shortenerMock.NewMockShortener(ctl)
//	ns := shortener.New(
//		shortener.SetEncoder(encoder),
//		shortener.InitBaseURL(baseURL),
//		shortener.SetLogger(zap.L()),
//	)
//
//	storageErr := errors.New("db is down")
//	encoder.EXPECT().GetAll().Return(nil, storageErr)
//
//	serverFunc := GetSavedData(ns).ServeHTTP
//	rec := httptest.NewRecorder()
//	req := httptest.NewRequest(
//		http.MethodGet,
//		"/user/urls",
//		nil,
//	)
//	req.Header.Set("content-type", "application/json")
//	serverFunc(rec, req)
//
//	res := rec.Result()
//	defer res.Body.Close()
//
//	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
//}
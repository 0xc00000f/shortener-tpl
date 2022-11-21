package encoder

import (
	"errors"
	"testing"

	storageMock "github.com/0xc00000f/shortener-tpl/internal/encoder/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/rand"
	"github.com/stretchr/testify/assert"
)

func TestURLEncoder_Encode(t *testing.T) {
	r := rand.New(false)
	tests := []struct {
		name    string
		letters int
	}{
		{
			name:    "6 letters url",
			letters: 6,
		},
		{
			name:    "72 letters url",
			letters: 72,
		},
		{
			name:    "0 letters url",
			letters: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ue := URLEncoder{length: tt.letters, rand: r}
			url := ue.encode()
			assert.Equal(t, len(url), tt.letters)
		})
	}
}

func TestURLEncoder_Short_Positive(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	tests := []struct {
		name  string
		short string
		long  string
	}{
		{
			name:  "positive #1",
			short: "BpLnfg", // first predictable result of ue.encode()
			long:  "https://google.com",
		},
		{
			name:  "positive #2",
			short: "Dsc2WD", // second predictable result of ue.encode()
			long:  "https://dzen.ru/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.EXPECT().IsKeyExist(tt.short).Return(false, nil)
			storage.EXPECT().Store(uuid.Nil, tt.short, tt.long).Return(nil)

			short, err := ue.Short(uuid.Nil, tt.long)
			require.NoError(t, err)
			assert.Equal(t, tt.short, short)
		})
	}
}

func TestURLEncoder_Short_IsKeyExist_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	expectedShort := "BpLnfg" // first predictable result of ue.encode()
	long := "https://dzen.ru/"
	storageErr := errors.New("db is down")

	storage.EXPECT().IsKeyExist(expectedShort).Return(false, storageErr)
	short, err := ue.Short(uuid.Nil, long)

	require.ErrorIs(t, err, storageErr)
	assert.Equal(t, "", short)

}

func TestURLEncoder_Short_Positive_IsKeyExist_IfExist(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	firstShort := "BpLnfg"  // first predictable result of ue.encode()
	secondShort := "Dsc2WD" // second predictable result of ue.encode()
	long := "https://dzen.ru/"

	storage.EXPECT().IsKeyExist(firstShort).Return(true, nil)
	storage.EXPECT().IsKeyExist(secondShort).Return(false, nil)
	storage.EXPECT().Store(uuid.Nil, secondShort, long).Return(nil)
	short, err := ue.Short(uuid.Nil, long)

	require.NoError(t, err)
	assert.Equal(t, secondShort, short)

}

func TestURLEncoder_Short_Store_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	expectedShort := "BpLnfg" // first predictable result of ue.encode()
	long := "https://dzen.ru/"
	storageErr := errors.New("db is down")

	storage.EXPECT().IsKeyExist(expectedShort).Return(false, nil)
	storage.EXPECT().Store(uuid.Nil, expectedShort, long).Return(storageErr)

	short, err := ue.Short(uuid.Nil, long)
	require.ErrorIs(t, err, storageErr)
	assert.Equal(t, "", short)

}

func TestURLEncoder_Get(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	short := "BpLnfg" // first predictable result of ue.encode()
	expectedLong := "https://dzen.ru/"
	storage.EXPECT().Get(short).Return(expectedLong, nil)

	long, err := ue.Get(short)
	require.NoError(t, err)
	assert.Equal(t, expectedLong, long)
}

func TestURLEncoder_Get_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := New(
		SetLength(PreferredLength),
		SetStorage(storage),
		SetLogger(zap.L()),
		SetRandom(rand.New(true)),
	)

	short := "BpLnfg" // first predictable result of ue.encode()
	expectedLong := ""
	storageErr := errors.New("db is down")
	storage.EXPECT().Get(short).Return("", storageErr)

	long, err := ue.Get(short)
	require.ErrorIs(t, err, storageErr)
	assert.Equal(t, expectedLong, long)
}

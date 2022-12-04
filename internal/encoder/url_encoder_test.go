package encoder_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	storageMock "github.com/0xc00000f/shortener-tpl/internal/encoder/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/0xc00000f/shortener-tpl/internal/rand"
)

var errStorageOutOfReach = errors.New("db is down")

func TestURLEncoder_Short_Positive(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	t.Cleanup(func() { ctl.Finish() })

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage.EXPECT().IsKeyExist(tt.short).Return(false, nil)
			storage.EXPECT().Store(uuid.Nil, tt.short, tt.long).Return(nil)

			short, err := ue.Short(uuid.Nil, tt.long)
			require.NoError(t, err)
			assert.Equal(t, tt.short, short)
		})
	}
}

func TestURLEncoder_Short_IsKeyExist_Error(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
	)

	const (
		expectedShort = "BpLnfg" // first predictable result of ue.encode()
		long          = "https://dzen.ru/"
	)

	storage.EXPECT().IsKeyExist(expectedShort).Return(false, errStorageOutOfReach)

	short, err := ue.Short(uuid.Nil, long)

	require.ErrorIs(t, err, errStorageOutOfReach)
	assert.Equal(t, "", short)
}

func TestURLEncoder_Short_Positive_IsKeyExist_IfExist(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
	)

	const (
		firstShort  = "BpLnfg" // first predictable result of ue.encode()
		secondShort = "Dsc2WD" // second predictable result of ue.encode()
		long        = "https://dzen.ru/"
	)

	storage.EXPECT().IsKeyExist(firstShort).Return(true, nil)
	storage.EXPECT().IsKeyExist(secondShort).Return(false, nil)
	storage.EXPECT().Store(uuid.Nil, secondShort, long).Return(nil)
	short, err := ue.Short(uuid.Nil, long)

	require.NoError(t, err)
	assert.Equal(t, secondShort, short)
}

func TestURLEncoder_Short_Store_Error(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
	)

	const (
		expectedShort = "BpLnfg" // first predictable result of ue.encode()
		long          = "https://dzen.ru/"
	)

	storage.EXPECT().IsKeyExist(expectedShort).Return(false, nil)
	storage.EXPECT().Store(uuid.Nil, expectedShort, long).Return(errStorageOutOfReach)

	short, err := ue.Short(uuid.Nil, long)
	require.ErrorIs(t, err, errStorageOutOfReach)
	assert.Equal(t, "", short)
}

func TestURLEncoder_Get(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
	)

	short := "BpLnfg" // first predictable result of ue.encode()
	expectedLong := "https://dzen.ru/"
	storage.EXPECT().Get(short).Return(expectedLong, nil)

	long, err := ue.Get(short)
	require.NoError(t, err)
	assert.Equal(t, expectedLong, long)
}

func TestURLEncoder_Get_Error(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	storage := storageMock.NewMockURLStorager(ctl)
	ue := encoder.New(
		encoder.SetLength(encoder.PreferredLength),
		encoder.SetStorage(storage),
		encoder.SetLogger(zap.L()),
		encoder.SetRandom(rand.New(true)),
	)

	short := "BpLnfg" // first predictable result of ue.encode()
	expectedLong := ""
	storage.EXPECT().Get(short).Return("", errStorageOutOfReach)

	long, err := ue.Get(short)
	require.ErrorIs(t, err, errStorageOutOfReach)
	assert.Equal(t, expectedLong, long)
}

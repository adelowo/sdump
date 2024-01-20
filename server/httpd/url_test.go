package httpd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/mocks"
	"github.com/sebdah/goldie/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func verifyMatch(t *testing.T, v interface{}) {
	g := goldie.New(t, goldie.WithFixtureDir("./testdata"))

	b := new(bytes.Buffer)

	var err error

	if d, ok := v.(*httptest.ResponseRecorder); ok {
		_, err = io.Copy(b, d.Body)
	} else {
		err = json.NewEncoder(b).Encode(v)
	}

	require.NoError(t, err)
	g.Assert(t, t.Name(), b.Bytes())
}

func TestURLHandler(t *testing.T) {
	tt := []struct {
		name               string
		mockFn             func(urlRepo *mocks.MockURLRepository)
		expectedStatusCode int
	}{
		{
			name: "could not create url",
			mockFn: func(urlRepo *mocks.MockURLRepository) {
				urlRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("could not create dump"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))

			logrus.SetOutput(io.Discard)

			logger := logrus.WithField("module", "test")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlRepo := mocks.NewMockURLRepository(ctrl)

			v.mockFn(urlRepo)

			u := &urlHandler{
				logger:  logger,
				cfg:     config.Config{},
				urlRepo: urlRepo,
			}

			u.create(recorder, req)

			require.Equal(t, v.expectedStatusCode, recorder.Result().StatusCode)

			verifyMatch(t, recorder)
		})
	}
}

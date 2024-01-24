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

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/mocks"
	"github.com/r3labs/sse/v2"
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

func TestURLHandler_Create(t *testing.T) {
	tt := []struct {
		name               string
		mockFn             func(urlRepo *mocks.MockURLRepository)
		expectedStatusCode int

		// sometimes data changes in the response, if this
		// field is set to true, we will skip matching golden files
		// technically it can be reworked to provide an implementation that never
		// changes during tests but I can always come back to taht
		hasDynamicData bool
	}{
		{
			name: "could not create url",
			mockFn: func(urlRepo *mocks.MockURLRepository) {
				urlRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("could not create dump"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "url was created",
			mockFn: func(urlRepo *mocks.MockURLRepository) {
				urlRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			hasDynamicData:     true,
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
				logger:    logger,
				cfg:       config.Config{},
				urlRepo:   urlRepo,
				sseServer: sse.New(),
			}

			u.create(recorder, req)

			require.Equal(t, v.expectedStatusCode, recorder.Result().StatusCode)

			if !v.hasDynamicData {
				verifyMatch(t, recorder)
			}
		})
	}
}

func TestURLHandler_Ingest(t *testing.T) {
	tt := []struct {
		name               string
		mockFn             func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository)
		expectedStatusCode int
		requestBody        io.Reader
		requestBodySize    int64
	}{
		{
			name: "url reference not found",
			mockFn: func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository) {
				urlRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&sdump.URLEndpoint{}, sdump.ErrURLEndpointNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			requestBody:        strings.NewReader(``),
			requestBodySize:    10,
		},
		{
			name: "error while fetching url",
			mockFn: func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository) {
				urlRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&sdump.URLEndpoint{}, errors.New("could not fetch url"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			requestBody:        strings.NewReader(``),
			requestBodySize:    10,
		},
		{
			name: "http request body too large",
			mockFn: func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository) {
				urlRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&sdump.URLEndpoint{}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			requestBody:        strings.NewReader(`{"name" : "Lanre", "occupation" :"Software"}`),
			requestBodySize:    10,
		},
		{
			name: "could not create ingestion",
			mockFn: func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository) {
				urlRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&sdump.URLEndpoint{}, nil)

				requestRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not insert into the database"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			requestBody:        strings.NewReader(`{"name" : "Lanre", "occupation" :"Software"}`),
			requestBodySize:    100,
		},
		{
			name: "ingested correctly",
			mockFn: func(urlRepo *mocks.MockURLRepository, requestRepo *mocks.MockIngestRepository) {
				urlRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&sdump.URLEndpoint{}, nil)

				requestRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
			requestBody:        strings.NewReader(`{"name" : "Lanre", "occupation" :"Software"}`),
			requestBodySize:    100,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", v.requestBody)

			logrus.SetOutput(io.Discard)

			logger := logrus.WithField("module", "test")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlRepo := mocks.NewMockURLRepository(ctrl)

			requestRepo := mocks.NewMockIngestRepository(ctrl)

			v.mockFn(urlRepo, requestRepo)

			u := &urlHandler{
				logger: logger,
				cfg: config.Config{
					HTTP: config.HTTPConfig{
						MaxRequestBodySize: v.requestBodySize,
					},
				},
				urlRepo:    urlRepo,
				ingestRepo: requestRepo,
				sseServer:  sse.New(),
			}

			u.ingest(recorder, req)

			require.Equal(t, v.expectedStatusCode, recorder.Result().StatusCode)
			verifyMatch(t, recorder)
		})
	}
}

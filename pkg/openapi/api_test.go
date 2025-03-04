package openapi

import (
	"bytes"
	"dvnetman/pkg/utils"
	"encoding/json"
	"github.com/google/uuid"
	mux2 "github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var testUUID1 = uuid.MustParse("637C7185-AEF8-42A9-91FC-F369E3970729")
var testTime1, _ = time.Parse(time.DateTime, time.DateTime)

type BufferCloser struct {
	*strings.Reader
}

func (b BufferCloser) Close() error {
	return nil
}

type LenReader interface {
	Len() int
}

func responseEqual(t testing.TB, expected, actual *http.Response) {
	t.Helper()
	actual.Header.Del("Date")
	actualBuffer := &bytes.Buffer{}
	require.NoError(t, actual.Write(actualBuffer))
	expectedBuffer := &bytes.Buffer{}
	expected.ProtoMajor = actual.ProtoMajor
	expected.ProtoMinor = actual.ProtoMinor
	if expected.Body != nil {
		expected.ContentLength = int64(expected.Body.(LenReader).Len())
	}
	require.NoError(t, expected.Write(expectedBuffer))
	require.Equal(t, expectedBuffer.String(), actualBuffer.String())
}

type responseBuilder struct {
	res *http.Response
}

func newResponseBuilder(code int) *responseBuilder {
	return &responseBuilder{res: &http.Response{StatusCode: code, Header: http.Header{}}}
}

func (r *responseBuilder) header(key, value string) *responseBuilder {
	r.res.Header.Add(key, value)
	return r
}

func (r *responseBuilder) body(body string) *responseBuilder {
	r.res.Body = &BufferCloser{strings.NewReader(body)}
	r.res.ContentLength = int64(len(body))
	return r
}

func (r *responseBuilder) json(j interface{}) *responseBuilder {
	body, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return r.header("Content-Type", "application/json").body(string(body))
}

func (r *responseBuilder) build() *http.Response {
	return r.res
}

func toJSONReadCloser(data interface{}) io.ReadCloser {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return &BufferCloser{strings.NewReader(string(b))}
}

func TestListDevices(t *testing.T) {
	searchResults1 := &DeviceSearchResults{
		Count: 1,
		Items: []*DeviceResult{
			{
				Id:      testUUID1,
				Name:    utils.ToPtr("test"),
				Created: &testTime1,
				Updated: &testTime1,
				DeviceType: &ObjectReference{
					Id:          testUUID1,
					DisplayName: utils.ToPtr("some-modal"),
				},
			},
		},
	}

	tests := []struct {
		name     string
		expect   func(s *MockDeviceAPI, e *MockErrorHandler)
		request  *http.Request
		response *http.Response
	}{
		{
			name: "no options",
			expect: func(s *MockDeviceAPI, e *MockErrorHandler) {
				s.EXPECT().ListDevices(mock.Anything, &ListDevicesOpts{Body: &DeviceSearchBody{}}).Return(
					&Response{
						Code:   200,
						Object: searchResults1,
					}, nil,
				)
			},
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/api/v1/device/search"},
				Body:   toJSONReadCloser(&DeviceSearchBody{}),
			},
			response: newResponseBuilder(200).json(searchResults1).build(),
		},
		{
			name: "all options",
			expect: func(s *MockDeviceAPI, e *MockErrorHandler) {
				s.EXPECT().ListDevices(
					mock.Anything, &ListDevicesOpts{
						Page:    utils.ToPtr(1),
						PerPage: utils.ToPtr(10),
						Sort:    utils.ToPtr("name"),
						Body: &DeviceSearchBody{
							Ids:        []uuid.UUID{uuid.MustParse("637C7185-AEF8-42A9-91FC-F369E3970729")},
							Fields:     []string{"name", "device_type", "status"},
							NameRegex:  utils.ToPtr("test1"),
							Name:       utils.ToPtr("test2"),
							Status:     utils.ToPtr("test3"),
							DeviceType: []uuid.UUID{uuid.MustParse("637C7185-AEF8-42A9-91FC-F369E3970730")},
						},
					},
				).Return(
					&Response{
						Code:   200,
						Object: searchResults1,
					}, nil,
				)
			},
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/api/v1/device/search", RawQuery: "page=1&per_page=10&sort=name"},
				Body: toJSONReadCloser(
					&DeviceSearchBody{
						Ids:        []uuid.UUID{uuid.MustParse("637C7185-AEF8-42A9-91FC-F369E3970729")},
						Fields:     []string{"name", "device_type", "status"},
						NameRegex:  utils.ToPtr("test1"),
						Name:       utils.ToPtr("test2"),
						Status:     utils.ToPtr("test3"),
						DeviceType: []uuid.UUID{uuid.MustParse("637C7185-AEF8-42A9-91FC-F369E3970730")},
					},
				),
			},
			response: newResponseBuilder(200).json(searchResults1).build(),
		},
		{
			name: "invalid page number",
			expect: func(s *MockDeviceAPI, e *MockErrorHandler) {
				e.EXPECT().ErrorHandler(
					mock.Anything, mock.Anything,
					mock.MatchedBy(
						func(err error) bool {
							return err.Error() == "invalid Query param page: strconv.Atoi: parsing \"abc\": invalid syntax"
						},
					),
				).Return()
			},
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/api/v1/device/search", RawQuery: "page=abc"},
				Body:   toJSONReadCloser(&DeviceSearchBody{}),
			},
			response: newResponseBuilder(200).build(),
		},
		{
			name: "invalid per page number",
			expect: func(s *MockDeviceAPI, e *MockErrorHandler) {
				e.EXPECT().ErrorHandler(
					mock.Anything, mock.Anything,
					mock.MatchedBy(
						func(err error) bool {
							return err.Error() == "invalid Query param per_page: strconv.Atoi: parsing \"abc\": invalid syntax"
						},
					),
				).Return()
			},
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/api/v1/device/search", RawQuery: "per_page=abc"},
				Body:   toJSONReadCloser(&DeviceSearchBody{}),
			},
			response: newResponseBuilder(200).build(),
		},
		{
			name: "handler error",
			expect: func(s *MockDeviceAPI, e *MockErrorHandler) {
				s.EXPECT().ListDevices(mock.Anything, mock.Anything).Return(nil, errors.New("handler error"))
				e.EXPECT().ErrorHandler(
					mock.Anything, mock.Anything,
					mock.MatchedBy(
						func(err error) bool {
							return err.Error() == "handler error"
						},
					),
				).Return()
			},
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/api/v1/device/search"},
				Body:   toJSONReadCloser(&DeviceSearchBody{}),
			},
			response: newResponseBuilder(200).build(),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				r := require.New(t)
				s := NewMockDeviceAPI(t)
				e := NewMockErrorHandler(t)
				mux := mux2.NewRouter()
				AttachDeviceAPI(s, e, mux)
				srv := httptest.NewServer(mux)
				defer srv.Close()
				test.expect(s, e)
				c := srv.Client()
				u, err := url.Parse(srv.URL)
				r.NoError(err)
				test.request.URL = u.ResolveReference(test.request.URL)
				res, err := c.Do(test.request)
				r.NoError(err)
				responseEqual(t, test.response, res)
			},
		)
	}
}

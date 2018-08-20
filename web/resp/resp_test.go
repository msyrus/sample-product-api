package resp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRender(t *testing.T) {
	type args struct {
		r    *http.Request
		resp Response
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody string
	}{
		{
			args: args{
				r: httptest.NewRequest("GET", "/test", nil),
				resp: Response{
					Code: 200,
					Data: true,
				},
			},
			wantCode: 200,
			wantBody: `{"data":true}`,
		},
		{
			args: args{
				r: httptest.NewRequest("GET", "/test", nil),
				resp: Response{
					Code: 200,
					Data: true,
				},
			},
			wantCode: 200,
			wantBody: `{"data":true}`,
		},
		{
			args: args{
				r: httptest.NewRequest("GET", "/test", nil),
				resp: Response{
					Code: 200,
					Data: []string{"test1", "test2"},
					Meta: NewPager(0, 1, 2),
				},
			},
			wantCode: 200,
			wantBody: `{"data":["test1","test2"],"meta":{"offset":1,"take":0,"total":0}}`,
		},
		{
			args: args{
				r: httptest.NewRequest("GET", "/test", nil),
				resp: Response{
					Code: 200,
					Data: []map[string]string{{"key1": "val1"}, {"key2": "val2"}},
					Meta: NewPager(3, 1, 4),
				},
			},
			wantCode: 200,
			wantBody: `{"data":[{"key1":"val1"},{"key2":"val2"}],"meta":{"offset":1,"take":2,"total":3}}`,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			Render(w, tt.args.r, tt.args.resp)
			if w.Code != tt.wantCode {
				t.Errorf("Render() code got = %v, wantCode = %v", w.Code, tt.wantCode)
			}
			if body := w.Body.String(); !reflect.DeepEqual(body, tt.wantBody) {
				t.Errorf("Render() body got = %v, wantBody = %v", string(body), tt.wantBody)
			}
		})
	}
}

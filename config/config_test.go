package config

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	testYml := `
host: "example.com"
port: 500
gracefulWait: 10
readTimeout: 5
writeTimeout: 0
idleTimeout: 2
postgres:
  uri: "postgres://db:5432/test"
`
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Application
		wantErr bool
	}{
		{
			args: args{
				r: &bytes.Buffer{},
			},
			want:    Application{},
			wantErr: true,
		},
		{
			args: args{
				r: bytes.NewBufferString("asd: 123"),
			},
			want:    Application{},
			wantErr: false,
		},
		{
			args: args{
				r: bytes.NewBufferString(testYml),
			},
			want: Application{
				GracefulWait: 10 * time.Second,
				ReadTimeout:  5 * time.Second,
				IdleTimeout:  2 * time.Second,
				Host:         "example.com",
				Port:         500,
				Postgres: Postgres{
					URI: "postgres://db:5432/test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

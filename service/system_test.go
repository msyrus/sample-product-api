package service

import (
	"reflect"
	"testing"

	"github.com/msyrus/simple-product-inv/version"
)

func TestNewSystem(t *testing.T) {
	tests := []struct {
		name string
		want *System
	}{
		{
			want: &System{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_Ready(t *testing.T) {
	tests := []struct {
		name    string
		s       *System
		want    bool
		wantErr bool
	}{
		{
			s:       NewSystem(),
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Ready()
			if (err != nil) != tt.wantErr {
				t.Errorf("System.Ready() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("System.Ready() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_Health(t *testing.T) {
	tests := []struct {
		name    string
		s       *System
		want    bool
		wantErr bool
	}{
		{
			s:       NewSystem(),
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Health()
			if (err != nil) != tt.wantErr {
				t.Errorf("System.Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("System.Health() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_Version(t *testing.T) {
	tests := []struct {
		name string
		s    *System
		want string
	}{
		{
			s:    NewSystem(),
			want: version.Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Version(); got != tt.want {
				t.Errorf("System.Version() = %v, want %v", got, tt.want)
			}
		})
	}
}

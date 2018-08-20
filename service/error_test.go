package service

import "testing"

func TestNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    NotFoundError
		want string
	}{
		{
			e:    ErrProductNotFound,
			want: "product not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("NotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

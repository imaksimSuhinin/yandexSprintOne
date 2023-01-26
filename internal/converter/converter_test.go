package converter

import (
	"reflect"
	"testing"
)

func TestFloat64ToBytes(t *testing.T) {
	type want struct {
		value [8]byte
	}

	var tests = []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				value: [8]byte{0, 0, 0, 0, 0, 0, 105, 64},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64ToBytes(200); !reflect.DeepEqual(got, tt.want.value) {
				t.Errorf("Float64ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

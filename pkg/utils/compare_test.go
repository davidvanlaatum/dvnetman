package utils

import "testing"

func TestComparePointers(t *testing.T) {
	type args struct {
		a *int
		b *int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test 1",
			args: args{
				a: nil,
				b: nil,
			},
			want: true,
		},
		{
			name: "Test 2",
			args: args{
				a: nil,
				b: new(int),
			},
			want: false,
		},
		{
			name: "Test 3",
			args: args{
				a: new(int),
				b: nil,
			},
			want: false,
		},
		{
			name: "Test 4",
			args: args{
				a: new(int),
				b: new(int),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePointers(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("ComparePointers() = %v, want %v", got, tt.want)
			}
		})
	}
}

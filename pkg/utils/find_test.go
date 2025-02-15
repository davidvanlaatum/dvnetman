package utils

import "testing"

func TestFindFirst(t *testing.T) {
	type args struct {
		a []int
		f func(int) bool
	}
	tests := []struct {
		name   string
		args   args
		want   int
		wantOk bool
	}{
		{
			name: "Test 1",
			args: args{
				a: []int{1, 2, 3, 4, 5},
				f: func(i int) bool {
					return i == 3
				},
			},
			want:   3,
			wantOk: true,
		},
		{
			name: "Test 2",
			args: args{
				a: []int{1, 2, 3, 4, 5},
				f: func(i int) bool {
					return i == 6
				},
			},
			want:   0,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := FindFirst(tt.args.a, tt.args.f)
			if got != tt.want {
				t.Errorf("FindFirst() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindFirst() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

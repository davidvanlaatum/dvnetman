package utils

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapTo_ConvertsElements(t *testing.T) {
	r := require.New(t)
	in := []int{1, 2, 3}
	want := []string{"1", "2", "3"}
	got := MapTo(in, func(i int) string {
		return fmt.Sprintf("%d", i)
	})
	r.Equal(want, got)
}

func TestMapTo_ConvertsElements_EmptyInput(t *testing.T) {
	r := require.New(t)
	in := []int{}
	want := []string{}
	got := MapTo(in, func(i int) string {
		return fmt.Sprintf("%d", i)
	})
	r.Equal(want, got)
}

func TestMapErr_ConvertsElements(t *testing.T) {
	r := require.New(t)
	in := []int{1, 2, 3}
	want := []string{"1", "2", "3"}
	got, err := MapErr(in, func(i int) (string, error) {
		return fmt.Sprintf("%d", i), nil
	})
	r.NoError(err)
	r.Equal(want, got)
}

func TestMapErr_ConvertsElements_EmptyInput(t *testing.T) {
	r := require.New(t)
	in := []int{}
	want := []string{}
	got, err := MapErr(in, func(i int) (string, error) {
		return fmt.Sprintf("%d", i), nil
	})
	r.NoError(err)
	r.Equal(want, got)
}

func TestMapErr_ErrorInFunction(t *testing.T) {
	r := require.New(t)
	in := []int{1, 2, 3}
	wantErr := errors.New("error")
	_, err := MapErr(in, func(i int) (string, error) {
		if i == 2 {
			return "", wantErr
		}
		return fmt.Sprintf("%d", i), nil
	})
	r.ErrorIs(err, wantErr)
}

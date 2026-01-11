package operations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedule_isMessagesEqual(t *testing.T) {
	table := []struct {
		A    string
		B    string
		want bool
	}{
		{
			A:    "",
			B:    "",
			want: true,
		},
		{
			A:    `Є відключення`,
			B:    `Є відключення`,
			want: true,
		},
		{
			A:    `Інформація станом на 10:00`,
			B:    `Інформація станом на 11:00`,
			want: true,
		},
		{
			A:    "Є відключення\nІнформація станом на 10:00",
			B:    "Є відключення\nІнформація станом на 11:00",
			want: true,
		},
		{
			A:    `Є відключення`,
			B:    `Нема відключення`,
			want: false,
		},
		{
			A:    `Є відключення`,
			B:    ``,
			want: false,
		},
	}

	for i, tt := range table {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			have := isMessagesEqual(tt.A, tt.B)
			assert.Equal(t, tt.want, have)
		})
	}
}

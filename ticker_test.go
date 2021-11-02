package time32

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestTicker(t *testing.T) {
	t.Run("reuse-time", func(t *testing.T) {
		tt := Now()
		reusedTt := ReuseTime()
		assert.Equal(t, tt.Unix(), reusedTt.Unix())
	})
	t.Run("reuse-unix", func(t *testing.T) {
		tt := Now()
		reusedTt := ReuseUnix()
		assert.Equal(t, tt.Unix(), reusedTt)
	})
	t.Run("reuse-nanos", func(t *testing.T) {
		tt := Now()
		reusedTt := ReuseUnixNano()
		diff := math.Abs(float64(tt.UnixNano()-reusedTt))
		fmt.Println(diff)
		assert.True(t, diff < 0.2*1000*1000)
	})
}
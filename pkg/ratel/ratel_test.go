package ratel

import (
	"ratel/pkg/ratel/key"
	"ratel/pkg/ratel/rconfig"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRatel(t *testing.T) {
	ratel := NewLimiter(
		WithDefaultRule(rconfig.ConfigItem{
			Capacity:            1,
			TimeWindowSecond:    1,
			BlockDurationSecond: 1,
		}),
		WithRuleItem(key.TokenKey("token"), rconfig.ConfigItem{
			Capacity:            2,
			TimeWindowSecond:    1,
			BlockDurationSecond: 1,
		}),
	)
	allow, err := ratel.Allow("token")
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = ratel.Allow("token")
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = ratel.Allow("token")
	assert.NoError(t, err)
	assert.False(t, allow)
	allow, err = ratel.Allow("ip")
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = ratel.Allow("ip")
	assert.NoError(t, err)
	assert.False(t, allow)
	allow, err = ratel.Allow("ip2")
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = ratel.Allow("ip")
	assert.NoError(t, err)
	assert.False(t, allow)

	time.Sleep(2 * time.Second)
	allow, err = ratel.Allow("ip")
	assert.NoError(t, err)
	assert.True(t, allow)
}

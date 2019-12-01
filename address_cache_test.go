package zstack

import (
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZStack_networkAddressCache(t *testing.T) {
	t.Run("added cache entry can be retrieved", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		z := New(unpiMock)
		defer unpiMock.Stop()

		ieee := zigbee.IEEEAddress(2754945194579457)
		expectedNetworkAddress := zigbee.NetworkAddress(426)

		z.addressUpdate(ieee, expectedNetworkAddress)

		actualNetworkAddress, found := z.addressLookup(ieee)
		assert.True(t, found)
		assert.Equal(t, expectedNetworkAddress, actualNetworkAddress)
	})

	t.Run("nonexistent cache entry returns not found", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		z := New(unpiMock)
		defer unpiMock.Stop()

		ieee := zigbee.IEEEAddress(2754945194579457)

    	_, found := z.addressLookup(ieee)
		assert.False(t, found)
	})

	t.Run("cache removal works", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		z := New(unpiMock)
		defer unpiMock.Stop()

		ieee := zigbee.IEEEAddress(2754945194579457)
		na := zigbee.NetworkAddress(426)

		z.addressUpdate(ieee, na)
		z.addressRemove(ieee)

		_, found := z.addressLookup(ieee)
		assert.False(t, found)
	})
}
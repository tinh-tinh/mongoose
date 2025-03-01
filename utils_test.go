package mongoose_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
)

func Test_IsValidateObjectID(t *testing.T) {
	require.False(t, mongoose.IsValidateObjectID("5d"))
	require.True(t, mongoose.IsValidateObjectID("5e874a551e9c3e148a8c530e"))
}

func Test_ToObjectID(t *testing.T) {
	require.Panics(t, func() {
		mongoose.ToObjectID("5d")
	})

	require.NotPanics(t, func() {
		mongoose.ToObjectID("5e874a551e9c3e148a8c530e")
	})
}

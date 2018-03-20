<%=licenseText%>
package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FromContext_NoFields(t *testing.T) {
	ctx := context.Background()

	logger := FromContext(ctx)

	assert.NotNil(t, logger)
	assert.Empty(t, logger.Data)
}

func Test_FromContext_Fields(t *testing.T) {
	ctx := context.Background()
	ctx = ContextWithFields(ctx, Fields{"SomeKey": "SomeValue"})

	logger := FromContext(ctx)

	require.NotNil(t, logger, "logger should not be nil")
	require.Equal(t, 1, len(logger.Data), "Fields does not contain the expected number of items")

	v, ok := logger.Data["SomeKey"]
	require.True(t, ok, "Fields does not contain expected key")
	require.Equal(t, "SomeValue", v)
}

func Test_AddFields(t *testing.T) {
	ctx := context.Background()

	// Service 1
	logger := FromContext(ctx)
	ctx, logger = AddFieldToCtx(ctx, "Service1Key", "Service1Value")

	require.Equal(t, 1, len(logger.Data), "Fields does not contain the expected number of items")

	// Service 2 (this is called from service 1, with the ctx)
	ctx, logger = AddFieldToCtx(ctx, "Service2Key", "Service2Value")

	require.Equal(t, 2, len(logger.Data), "Fields does not contain the expected number of items")
}

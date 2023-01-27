package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	var m Metrics
	m.UpdateMetrics()

	assert.Equal(t, 1, int(m.PollCount))
}

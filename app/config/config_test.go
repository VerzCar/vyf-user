package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const ttlDefault uint = 300

func TestNewConfig(t *testing.T) {
	c := NewConfig(".")

	require.Equal(t, c.Ttl.Default, ttlDefault)
}

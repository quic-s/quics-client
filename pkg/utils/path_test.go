package utils_test

import (
	//"path"
	"testing"

	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLocalRelToRoot(t *testing.T) {
	abs := utils.LocalRelToRoot("../qic", "/workspaces/quics-client")
	assert.Equal(t, "/quics-client/pkg/qic", abs)
}

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTemplates_None(t *testing.T) {
	box := packr.NewBox("/nonexistent")
	require.NotNil(t, box)

	tmpls, err := loadTemplates(box)
	require.Error(t, err)
	assert.Empty(t, tmpls)
}

func TestLoadTemplates_Empty(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "loadResources")
	require.NoError(t, err)

	defer os.RemoveAll(tmpdir)

	box := packr.NewBox(tmpdir)
	require.NotNil(t, box)

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)
	assert.NotNil(t, tmpls)
	assert.Empty(t, tmpls)
}

func TestLoadTemplates_Good(t *testing.T) {
	box := packr.NewBox("./files")

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)

	assert.NotNil(t, tmpls)
	assert.NotEmpty(t, tmpls)
}

func TestLoadTemplates_GoodDebug(t *testing.T) {
	fDebug = true
	box := packr.NewBox("./files")

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)
	assert.NotNil(t, tmpls)
	assert.NotEmpty(t, tmpls)

	fDebug = false
}

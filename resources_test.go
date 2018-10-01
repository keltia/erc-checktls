package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadContractFile(t *testing.T) {
	// We embed the file now
	box := packr.NewBox("./files")

	cntrs, err := readContractFile(box)
	assert.NoError(t, err)
	assert.NotEmpty(t, cntrs)
}

func TestLoadTemplates(t *testing.T) {
	// We embed the file now
	box := packr.NewBox("./files")

	str, err := loadTemplates(box)
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
	assert.Equal(t, 2, len(str))
}

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

	fl, err := filepath.Glob("files/*.html")
	require.NoError(t, err)

	assert.Equal(t, len(fl), len(tmpls))

	assert.NotEmpty(t, tmpls["templ.html"])
	assert.NotEmpty(t, tmpls["summaries.html"])
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

func TestLoadResources_GoodDebug(t *testing.T) {
	err := loadResources(resourcesPath)
	assert.NoError(t, err)

	assert.NotEmpty(t, tmpls)
	assert.NotEmpty(t, contracts)
}

func TestLoadResources_None(t *testing.T) {
	err := loadResources("/nonexistent")
	assert.Error(t, err)

	assert.Empty(t, tmpls)
	assert.Empty(t, contracts)
}

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadContractFile(t *testing.T) {
	// We embed the file now
	box := packr.New("test", "./files")

	cntrs, err := readContractFile(box)
	assert.NoError(t, err)
	assert.NotEmpty(t, cntrs)
}

func TestLoadTemplates(t *testing.T) {
	// We embed the file now
	box := packr.New("test", "./files")

	str, err := loadTemplates(box)
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
	assert.Equal(t, 2, len(str))
}

func TestLoadTemplates_None(t *testing.T) {
	box := packr.New("testnone", "/nonexistent")
	require.NotNil(t, box)

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)
	assert.Empty(t, tmpls)
}

func TestLoadTemplates_Empty(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "loadResources")
	require.NoError(t, err)

	defer os.RemoveAll(tmpdir)

	box := packr.New("testempty", tmpdir)
	require.NotNil(t, box)

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)
	assert.NotNil(t, tmpls)
	assert.Empty(t, tmpls)
}

func TestLoadTemplates_Good(t *testing.T) {
	box := packr.New("testgood", "./files")

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
	box := packr.New("test", "./files")

	tmpls, err := loadTemplates(box)
	require.NoError(t, err)
	assert.NotNil(t, tmpls)
	assert.NotEmpty(t, tmpls)

	fDebug = false
}

func TestLoadResources_GoodDebug(t *testing.T) {
	fDebug = true
	c, tt, err := loadResources()
	assert.NoError(t, err)

	assert.NotEmpty(t, tt)
	assert.NotEmpty(t, c)
	fDebug = false
}

func TestLoadResources_Good(t *testing.T) {
	c, tt, err := loadResources()
	assert.NoError(t, err)

	assert.NotEmpty(t, tt)
	assert.NotEmpty(t, c)
}

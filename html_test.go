package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBooleanT(t *testing.T) {
	assert.Equal(t, white("TRUE"), booleanT(true))
	assert.Equal(t, red("FALSE"), booleanT(false))
}

func TestBooleanF(t *testing.T) {
	assert.Equal(t, white("FALSE"), booleanF(false))
	assert.Equal(t, red("TRUE"), booleanF(true))
}

func TestRed(t *testing.T) {
	str := ""
	assert.Equal(t, `<td class=xl64 align=center></td>`, red(str))
}

func TestRed1(t *testing.T) {
	str := "foobar"
	assert.Equal(t, `<td class=xl64 align=center>foobar</td>`, red(str))
}

func TestGrade(t *testing.T) {
	td := []struct {
		In   string
		HTML string
	}{
		{"A+", "<td class=xl65 align=center>A+</td>"},
		{"A", "<td class=xl65 align=center>A</td>"},
		{"A-", "<td class=xl631 align=center>A-</td>"},
		{"B+", "<td class=xl63 align=center>B+</td>"},
		{"B", "<td class=xl63 align=center>B</td>"},
		{"B-", "<td class=xl63 align=center>B-</td>"},
		{"C+", "<td class=xl64 align=center>C+</td>"},
		{"C", "<td class=xl64 align=center>C</td>"},
		{"C-", "<td class=xl64 align=center>C-</td>"},
		{"D+", "<td class=xl64 align=center>D+</td>"},
		{"D", "<td class=xl64 align=center>D</td>"},
		{"D-", "<td class=xl64 align=center>D-</td>"},
		{"E+", "<td class=xl64 align=center>E+</td>"},
		{"E", "<td class=xl64 align=center>E</td>"},
		{"E-", "<td class=xl64 align=center>E-</td>"},
		{"F", "<td class=xl64 align=center>F</td>"},
		{"Z", "<td class=xl661 align=center>&nbsp;</td>"},
	}
	for _, tst := range td {
		assert.Equal(t, tst.HTML, grade(tst.In))
	}
}

func TestText(t *testing.T) {
	str := "foobar"
	assert.Equal(t, "<td height=21 style='height:16.0pt'>foobar</td>", text(str))
}

func TestProto(t *testing.T) {
	td := []struct {
		In   string
		HTML string
	}{
		{"TLSv1.2", "<td class=xl65 align=center>TLSv1.2</td>"},
		{"TLSv1.1,TLSv1.2", "<td class=xl631 align=center>TLSv1.1,TLSv1.2</td>"},
		{"TLSv1.0,TLSv1.1,TLSv1.2", "<td class=xl631 align=center>TLSv1.0,TLSv1.1,TLSv1.2</td>"},
		{"SSLv3.0,TLSv1.0", "<td class=xl64 align=center>SSLv3.0,TLSv1.0</td>"},
		{"foobar", "<td class=xl661 align=center>foobar</td>"},
	}
	for _, tst := range td {
		assert.Equal(t, tst.HTML, proto(tst.In))
	}
}

func TestHSTS(t *testing.T) {
	td := []struct {
		In   int64
		HTML string
	}{
		{-1, "<td class=xl64 align=center>NO</td>"},
		{3600, "<td class=xl63 align=center>3600</td>"},
		{5184100, "<td class=xl63 align=center>5184100</td>"},
		{15550000, "<td class=xl631 align=center>15550000</td>"},
		{15552000, "<td class=xl65 align=center>15552000</td>"},
	}
	for _, tst := range td {
		assert.Equal(t, tst.HTML, hsts(tst.In))
	}
}

func TestTLSReport_ToHTML(t *testing.T) {
	var buf strings.Builder

	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewTLSReport(all)
	require.NoError(t, err)

	raw, err := ioutil.ReadFile("files/templ.html")
	tmpl := string(raw)
	require.NoError(t, err)
	require.NotEmpty(t, tmpl)

	err = sites.ToHTML(&buf, tmpl)
	assert.NoError(t, err)
}

func TestWriteHTML(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	https := map[string]int{
		"A":  666,
		"B+": 37,
		"F":  42,
	}

	err := WriteHTML(os.Stderr, nil, cntrs, https)
	assert.Error(t, err)
}

func TestWriteHTML2(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	https := map[string]int{
		"A":  666,
		"B+": 37,
		"F":  42,
	}

	err := WriteHTML(os.Stderr, &TLSReport{}, cntrs, https)
	assert.Error(t, err)
}

func TestWriteHTML3(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	https := map[string]int{
		"A":  666,
		"B+": 37,
		"F":  42,
	}

	file := "testdata/site.json"
	raw, err := getResults(file)
	require.NoError(t, err)

	allSites, err := ssllabs.ParseResults(raw)
	require.NoError(t, err)

	box := packr.NewBox("./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewTLSReport(allSites)
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(t, err)

	err = WriteHTML(null, final, cntrs, https)
	assert.NoError(t, err)
}

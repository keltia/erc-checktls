package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
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
		{"C+", "<td class=xl63 align=center>C+</td>"},
		{"C", "<td class=xl63 align=center>C</td>"},
		{"C-", "<td class=xl63 align=center>C-</td>"},
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

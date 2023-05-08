package cpre

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testIncluder(filePath string, global bool) (string, []byte, error) {
	p := filepath.Join("examples", filePath)

	id, err := filepath.Abs(p)
	if err != nil {
		return "", nil, err
	}

	source, err := os.ReadFile(p)

	return id, source, err
}

func testPreprocess(t *testing.T, name string) {
	p := NewPreprocessor(PreprocessorConfig{
		Include: testIncluder,
	})

	testPreprocessWithPreprocessor(t, name, p)
}

func testPreprocessWithPreprocessor(t *testing.T, name string, p *Preprocessor) {
	source, err := os.ReadFile("examples/" + name + ".cpp")
	assert.NoError(t, err)
	expected, err := os.ReadFile("examples/" + name + ".pre.cpp")
	assert.NoError(t, err)

	actual := p.Process(string(source))

	assert.Equal(t, string(expected), actual)
}

func TestSingleLineComment(t *testing.T) {
	testPreprocess(t, "single_line_comment")
}

func TestMultiLineComment(t *testing.T) {
	testPreprocess(t, "multi_line_comment")
}

func TestSpacesInDefine(t *testing.T) {
	testPreprocess(t, "spaces_in_define")
}

func TestMulti(t *testing.T) {
	testPreprocess(t, "multi")
}

func TestCircular(t *testing.T) {
	testPreprocess(t, "circular")
}

func TestCircularValue(t *testing.T) {
	testPreprocess(t, "circular_value")
}

func TestIf(t *testing.T) {
	testPreprocess(t, "if")
}

func TestOnce(t *testing.T) {
	testPreprocess(t, "once")
}

func TestNotOnce(t *testing.T) {
	testPreprocess(t, "not_once")
}

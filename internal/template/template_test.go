package template

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultTemplate(t *testing.T) {
	tmpDir := "../../testdata/tmp"
	t.Run("success", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		Setup(tmpDir)
		err := CreateDefaultTemplate()
		assert.Nil(t, err)
		os.RemoveAll(tmpDir)
	})

	t.Run("already existing", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		Setup(tmpDir)
		err := CreateDefaultTemplate()
		assert.Nil(t, err)

		Setup(tmpDir)
		err = CreateDefaultTemplate()
		assert.Nil(t, err)
		os.RemoveAll(tmpDir)
	})

	t.Run("stat error", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		// create unwritable dir
		err := os.Mkdir(tmpDir, 0444)
		assert.Nil(t, err)

		Setup(tmpDir)
		err = createDefaultTemplate(func(path string) (os.FileInfo, error) {
			return nil, errors.New("Mock")
		})
		assert.NotNil(t, err)
		os.RemoveAll(tmpDir)
	})
}

func TestExecuteTemplateFromFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Setup("../../testdata")
		res, err := ExecuteTemplateFromFile("template.html", map[string]string{"var": "test"})
		assert.Nil(t, err)
		assert.Equal(t, "<test>", res)
	})

	t.Run("error", func(t *testing.T) {
		res, err := ExecuteTemplateFromFile("invalid-template.html", nil)
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})

	t.Run("error2", func(t *testing.T) {
		res, err := ExecuteTemplateFromString("{{ .var.var }}", map[string]string{"var": "test"})
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})
}

func TestExecuteTemplateFromString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := ExecuteTemplateFromString("<{{ .var }}>", map[string]string{"var": "test"})
		assert.Nil(t, err)
		assert.Equal(t, "<test>", res)
	})

	t.Run("success with print", func(t *testing.T) {
		res, err := ExecuteTemplateFromString("<{{ print .var }}>", map[string][]string{"var": {"test"}})
		assert.Nil(t, err)
		assert.Equal(t, "<test>", res)
	})

	t.Run("error", func(t *testing.T) {
		res, err := ExecuteTemplateFromString("{{notvalid}}", nil)
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})
}

func TestPrint(t *testing.T) {
	assert.Equal(t, "test", print("test"))
	assert.Equal(t, "testtest", print([]string{"test", "test"}))
	assert.Equal(t, "", print(nil))
}

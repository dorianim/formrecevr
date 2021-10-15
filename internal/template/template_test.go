package template_test

import (
	"os"
	"testing"

	"github.com/dorianim/formrecevr/internal/template"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewTemplate(t *testing.T) {
	tmpDir := "../../testdata/tmp"
	t.Run("success", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		err := template.CreateDefaultTemplate(tmpDir)
		assert.Nil(t, err)
		os.RemoveAll(tmpDir)
	})

	t.Run("already existing", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		err := template.CreateDefaultTemplate(tmpDir)
		assert.Nil(t, err)

		err = template.CreateDefaultTemplate(tmpDir)
		assert.Nil(t, err)
		os.RemoveAll(tmpDir)
	})

	t.Run("unwritable dir", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		// create unwritable dir
		err := os.Mkdir(tmpDir, 0444)
		assert.Nil(t, err)

		err = template.CreateDefaultTemplate(tmpDir)
		assert.NotNil(t, err)
		os.RemoveAll(tmpDir)
	})
}

func TestExecuteTemplateFromFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := template.ExecuteTemplateFromFile("../../testdata/template.html", map[string]string{"var": "test"})
		assert.Nil(t, err)
		assert.Equal(t, "<test>", res)
	})

	t.Run("error", func(t *testing.T) {
		res, err := template.ExecuteTemplateFromFile("../../testdata/invalid-template.html", nil)
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})

	t.Run("error2", func(t *testing.T) {
		res, err := template.ExecuteTemplateFromString("{{ .var.var }}", map[string]string{"var": "test"})
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})
}

func TestExecuteTemplateFromString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := template.ExecuteTemplateFromString("<{{ .var }}>", map[string]string{"var": "test"})
		assert.Nil(t, err)
		assert.Equal(t, "<test>", res)
	})

	t.Run("error", func(t *testing.T) {
		res, err := template.ExecuteTemplateFromString("{{notvalid}}", nil)
		assert.NotNil(t, err)
		assert.Equal(t, "", res)
	})
}

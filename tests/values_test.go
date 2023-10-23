package tests

import (
	"testing"

	"main/render"
)

func TestRunRender(t *testing.T) {

	t.Run("template-values.yaml", func(t *testing.T) {
		valuesRender := new(render.ValuesRenderer)
		valuesRender.Run("template-values.yaml", "")
	})

	t.Run("base1-values1.yaml", func(t *testing.T) {
		valuesRender := new(render.ValuesRenderer)
		valuesRender.Run("base-values1.yaml", "")
	})

	t.Run("test-glob-values.yaml", func(t *testing.T) {
		valuesRender := new(render.ValuesRenderer)
		valuesRender.Run("test-glob-values.yaml", "")
	})

	// t.Run("extended-values1.yaml", func(t *testing.T) {
	// 	valuesRender := new(render.ValuesRenderer)
	// 	valuesRender.Run("extended-values1.yaml", "")
	// })
}

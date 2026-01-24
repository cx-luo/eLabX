// Package enote coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/6 14:10
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : img_render.go
// @Software: GoLand
package enote

import (
	"github.com/cx-luo/go-indigo/core"
)

func RenderReaction(indigoInit *core.Indigo, rxnSmiles string) ([]byte, error) {
	indigoRenderer, err := indigoInit.InitRenderer()
	if err != nil {
		return nil, err
	}
	defer indigoRenderer.DisposeRenderer()
	indigoRenderer.ResetRenderer()
	// Set Indigo chemical processing and rendering options as per the provided prompt

	// Chemical options
	indigoRenderer.SetRenderOptionBool("smart-layout", true)
	indigoRenderer.SetRenderOptionBool("ignore-stereochemistry-errors", true)
	indigoRenderer.SetRenderOptionBool("mass-skip-error-on-pseudoatoms", false)
	indigoRenderer.SetRenderOptionBool("gross-formula-add-rsites", true)
	indigoRenderer.SetRenderOptionBool("aromatize-skip-superatoms", true)
	indigoRenderer.SetRenderOptionBool("dearomatize-on-load", false)
	indigoRenderer.SetRenderOptionBool("ignore-no-chiral-flag", false)
	indigoRenderer.SetRenderOption("aromaticity-model", "generic")
	indigoRenderer.SetRenderOptionBool("gross-formula-add-isotopes", true)

	// Rendering options
	indigoRenderer.SetRenderOptionBool("render-coloring", true)
	indigoRenderer.SetRenderOptionInt("render-font-size", 13)
	indigoRenderer.SetRenderOption("render-font-size-unit", "px")
	indigoRenderer.SetRenderOptionInt("render-font-size-sub", 13)
	indigoRenderer.SetRenderOption("render-font-size-sub-unit", "px")
	indigoRenderer.SetRenderOptionInt("image-resolution", 72)
	indigoRenderer.SetRenderOption("bond-length-unit", "px")
	indigoRenderer.SetRenderOptionInt("bond-length", 40)
	indigoRenderer.SetRenderOptionInt("render-bond-thickness", 2)
	indigoRenderer.SetRenderOption("render-bond-thickness-unit", "px")
	indigoRenderer.SetRenderOptionFloat("render-bond-spacing", 0.15)
	indigoRenderer.SetRenderOptionInt("render-stereo-bond-width", 6)
	indigoRenderer.SetRenderOption("render-stereo-bond-width-unit", "px")
	indigoRenderer.SetRenderOptionFloat("render-hash-spacing", 1.2)
	indigoRenderer.SetRenderOption("render-hash-spacing-unit", "px")
	indigoRenderer.SetRenderOption("render-label-mode", "terminal-hetero")

	if err := indigoRenderer.SetRenderOption("render-output-format", "svg"); err != nil {
		return nil, err
	}
	if err := indigoRenderer.SetRenderOption("render-background-color", "255, 255, 255"); err != nil {
		return nil, err
	}

	rxn, err := indigoInit.LoadReactionFromString(rxnSmiles)
	if err != nil {
		return nil, err
	}
	defer rxn.Close()

	outputHandle, err := indigoInit.CreateWriteBuffer()
	if err != nil {
		return nil, err
	}

	defer indigoInit.FreeObject(outputHandle)

	err = indigoRenderer.Render(rxn.Handle, outputHandle)
	if err != nil {
		return nil, err
	}

	imageSvgData, err := indigoInit.GetBufferData(outputHandle)
	if err != nil {
		return nil, err
	}

	return imageSvgData, nil
}

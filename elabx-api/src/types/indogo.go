// Package types coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/6 14:04
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : indogo.go
// @Software: GoLand
package types

type MoleculeCalculated struct {
	Gross            string `json:"gross,omitempty"`
	MolecularWeight  string `json:"molecular-weight,omitempty"`
	MostAbundantMass string `json:"most-abundant-mass,omitempty"`
	MonoisotopicMass string `json:"monoisotopic-mass,omitempty"`
	MassComposition  string `json:"mass-composition,omitempty"`
}

type Molecule struct {
	Name             string  `json:"name"`
	InChi            string  `json:"inChi"`
	InChiKey         string  `json:"inChiKey"`
	MoleWeight       float64 `json:"moleWeight"`
	MostAbundantMass float64 `json:"mostAbundantMass"`
	Formula          string  `json:"formula"`
}

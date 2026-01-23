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
	Name              string  `json:"name"`
	InChi             string  `json:"inChi"`
	InChiKey          string  `json:"inChiKey"`
	Formula           string  `json:"formula"`
	CxSmiles          string  `json:"cxSmiles"`
	ExactMass         float64 `json:"exactMass"`
	MolecularWeight   float64 `json:"molecularWeight"`
	MostAbundantMass  float64 `json:"mostAbundantMass"`
	MonoisotopicMass  float64 `json:"monoisotopicMass"`
	MassComposition   string  `json:"massComposition"`
	TPSA              float64 `json:"tpsa"`
	NumRotatableBonds int     `json:"numRotatableBonds"`
	NumHeteroatoms    int     `json:"numHeteroatoms"`
	NumHeavyAtoms     int     `json:"numHeavyAtoms"`
	NumAtoms          int     `json:"numAtoms"`
	NumBonds          int     `json:"numBonds"`
}

type Reaction struct {
	ReactionId     int64      `json:"reactionId"`
	ReactionSmiles string     `json:"reactionSmiles"`
	CdStructure    string     `json:"cdStructure"`
	CxSmiles       string     `json:"cxSmiles"`
	Reactants      []Molecule `json:"reactants"`
	Products       []Molecule `json:"products"`
	Reagents       []Molecule `json:"reagents"`
}

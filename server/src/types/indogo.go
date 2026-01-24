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
	ImageSvg       []byte     `json:"imageSvg"`
	Reactants      []Molecule `json:"reactants"`
	Products       []Molecule `json:"products"`
	Reagents       []Molecule `json:"reagents"`
}

type ReagentTableDataStruct struct {
	ReagentId        int64   `json:"reagentId"`
	ReagentName      string  `json:"reagentName"`
	ReagentSmiles    string  `json:"reagentSmiles"`
	Mw               float64 `json:"mw"`
	ReagentRole      string  `json:"reagentRole"`
	Formula          string  `json:"formula"`
	Cas              string  `json:"cas"`
	Eq               float64 `json:"eq"`
	Purity           float64 `json:"purity"`
	Quantity         float64 `json:"quantity"`
	QuantityUnit     string  `json:"quantityUnit"`
	Concentration    float64 `json:"concentration"`
	Density          float64 `json:"density"`
	DensityUnit      string  `json:"densityUnit"`
	CompoundId       int64   `json:"compoundId"`
	Yield            float64 `json:"yield"`
	IsLimiting       int     `json:"isLimiting"`
	IsChiral         int     `json:"isChiral"`
	StereoCentersCnt int     `json:"stereoCentersCnt"`
	ChiralDescriptor string  `json:"chiralDescriptor"`
	ReagentImg       []byte  `json:"reagentImg"`
	ProductAlias     string  `json:"productAlias"`
	Moles            float64 `json:"moles"`
	MolesUnit        string  `json:"molesUnit"`
	Volume           float64 `json:"volume"`
	VolumeUnit       string  `json:"volumeUnit"`
	CdStructure      string  `json:"cdStructure,omitempty"`
}

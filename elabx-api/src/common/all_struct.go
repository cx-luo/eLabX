// Package common coding=utf-8
// @Project : eLabX
// @Time    : 2024/2/23 14:48
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : all_struct.go
// @Software: GoLand
package common

type indigoApi struct {
	Render string
}

// Indigo todo: finalExport API
var Indigo = indigoApi{
	Render: "http://192.168.2.139:18100/v2/indigo/render",
}

type jchemApi struct {
	GetReactionDetailsApi string
	GetStereoApi          string
	ReactionExport        string
	GetSingleMolApi       string
	MolExport             string
}

var JChemApi = jchemApi{
	GetReactionDetailsApi: "http://192.168.1.79:32527/v1/converter/getReactionDetails",
	GetStereoApi:          "http://192.168.2.139:8015/rest-v1/jws/stereo/cip",
	ReactionExport:        "http://10.10.5.72:8080/rest-v1/util/calculate/reactionExport",
	GetSingleMolApi:       "http://192.168.1.79:32527/v1/converter/resolveSingle",
	MolExport:             "http://10.10.5.72:8080/rest-v1/util/calculate/molExport",
}

type UsedProps struct {
	Cid              int     `json:"cid"`
	Cmpdname         string  `json:"cmpdname"`
	Mf               string  `json:"mf"`
	Mw               float64 `json:"mw"`
	Smiles           string  `json:"smiles"`
	Exactmass        float64 `json:"exactmass,omitempty"`
	Monoisotopicmass float64 `json:"monoisotopicmass,omitempty"`
	Inchi            string  `json:"inchi"`
	Inchikey         string  `json:"inchikey"`
	Iupacname        string  `json:"iupacname"`
}

type ReagentInfo struct {
	ReagentID        int64   `json:"reagentId,omitempty" db:"reagent_id"`
	ReagentName      string  `json:"reagentName,omitempty" db:"reagent_name"`
	ReagentSmiles    string  `json:"reagentSmiles,omitempty" db:"reagent_smiles"`
	CdStructure      string  `json:"cdStructure,omitempty" db:"cd_structure"`
	Mw               float64 `json:"mw,omitempty" db:"mw"`
	ReagentRole      string  `json:"reagentRole,omitempty" db:"reagent_role"`
	Formula          string  `json:"formula,omitempty" db:"formula"`
	Cas              string  `json:"cas,omitempty" db:"cas"`
	Purity           float64 `json:"purity,omitempty" db:"purity"`
	Quantity         float64 `json:"quantity,omitempty" db:"quantity"`
	QuantityUnit     string  `json:"quantityUnit,omitempty" db:"quantity_unit"`
	Eq               float64 `json:"eq,omitempty" db:"equiv"`
	Concentration    float64 `json:"concentration,omitempty" db:"concentration"`
	Density          float64 `json:"density,omitempty" db:"density"`
	CompoundID       int     `json:"compoundId,omitempty" db:"compound_id"`
	IsLimiting       int     `json:"isLimiting,omitempty" db:"is_limiting"`
	IsChiral         int     `json:"isChiral,omitempty" db:"is_chiral"`
	StereoCentersCnt int     `json:"stereoCentersCnt" db:"stereo_centers_cnt"`
	ChiralDescriptor string  `json:"chiralDescriptor,omitempty" db:"chiral_descriptor"`
	ProductAlias     string  `json:"productAlias,omitempty" db:"product_alias"`
	ReagentImg       string  `json:"reagentImg,omitempty" db:"reagent_img"`
	Moles            float64 `json:"moles,omitempty" db:"moles"`
	MolesUnit        string  `json:"molesUnit,omitempty" db:"moles_unit"`
	Volume           float64 `json:"volume,omitempty" db:"volume"`
	VolumeUnit       string  `json:"volumeUnit,omitempty" db:"volume_unit"`
}

type ElnRxnReagents struct {
	ReagentID        int64   `json:"reagentId,omitempty" db:"reagent_id"`
	ReactionId       int64   `json:"reactionId,omitempty" db:"reaction_id"`
	ReagentName      string  `json:"reagentName,omitempty" db:"reagent_name"`
	ReagentSmiles    string  `json:"reagentSmiles,omitempty" db:"reagent_smiles"`
	Mw               float32 `json:"mw,omitempty" db:"mw"`
	Exactmass        float32 `json:"exactmass,omitempty" db:"exactmass"`
	MonoisotopicMass float32 `json:"monoisotopicMass,omitempty" db:"monoisotopic_mass"`
	Formula          string  `json:"formula,omitempty" db:"formula"`
	ReagentRole      string  `json:"reagentRole,omitempty" db:"reagent_role"`    // role
	Equiv            float32 `json:"equiv,omitempty" db:"equiv"`                 // eq
	Cas              string  `json:"cas,omitempty" db:"cas"`                     // CAS#
	Concentration    float32 `json:"concentration,omitempty" db:"concentration"` // Conv. 浓度
	CdStructure      string  `json:"cdStructure,omitempty" db:"cd_structure"`
	Inchi            string  `json:"inchi,omitempty" db:"inchi"`
	Inchikey         string  `json:"inchikey,omitempty" db:"inchikey"`
	Cxsmiles         string  `json:"cxsmiles,omitempty" db:"cxsmiles"`
	Density          float32 `json:"density,omitempty" db:"density"` // \'g/mL\'
	Quantity         float32 `json:"quantity,omitempty" db:"quantity"`
	QuantityUnit     string  `json:"quantityUnit,omitempty" db:"quantity_unit"`
	Purity           float32 `json:"purity,omitempty" db:"purity"`
	CompoundId       int64   `json:"compoundId,omitempty" db:"compound_id"`
	IsChiral         int8    `json:"isChiral,omitempty" db:"is_chiral"`
	ReagentHash      int64   `json:"reagentHash,omitempty" db:"reagent_hash"`
	IsLimiting       int8    `json:"isLimiting,omitempty" db:"is_limiting"`
	StereoCentersCnt int64   `json:"stereoCentersCnt,omitempty" db:"stereo_centers_cnt"`
	ProductAlias     string  `json:"productAlias,omitempty" db:"product_alias"`
	ChiralDescriptor string  `json:"chiralDescriptor,omitempty" db:"chiral_descriptor"`
	Moles            float32 `json:"moles,omitempty" db:"moles"`
	MolesUnit        string  `json:"molesUnit,omitempty" db:"moles_unit"`
	Volume           float32 `json:"volume,omitempty" db:"volume"`
	VolumeUnit       string  `json:"volumeUnit,omitempty" db:"volume_unit"`
}
type ElnRxnBasicInfo struct {
	ReactionId    int64  `json:"reactionId" db:"reaction_id"`
	ProjectName   string `json:"projectName,omitempty" db:"project_name"`
	Batch         int64  `json:"batch" db:"batch"`
	AuthorId      int64  `json:"authorId,omitempty" db:"author_id"`
	StepId        string `json:"stepId" db:"step_id"`
	AuthorName    string `json:"authorName,omitempty" db:"author_name"`
	WitnessId     int64  `json:"witnessId,omitempty" db:"witness_id"`
	WitnessName   string `json:"witnessName,omitempty" db:"witness_name"`
	RxnType       string `json:"rxnType,omitempty" db:"rxn_type"`
	Reference     string `json:"reference,omitempty" db:"reference"`
	Doi           string `json:"doi,omitempty" db:"doi"`
	CreationDate  string `json:"creationDate,omitempty" db:"creation_date"`
	StartDate     string `json:"startDate" db:"start_date"`
	RxnStatus     string `json:"rxnStatus" db:"rxn_status"`
	Comment       string `json:"comment,omitempty" db:"comment"`
	RxnTypeFromAi string `json:"rxnTypeFromAi,omitempty" db:"rxn_type_from_ai"`
	RxnTypeCode   string `json:"rxnTypeCode,omitempty" db:"rxn_type_code"`
	PageName      string `json:"pageName,omitempty" db:"page_name"`
	CommitDate    string `json:"commitDate,omitempty" db:"commit_date"`
}

type Samples struct {
	ReactionId          int64   `json:"reactionId,omitempty" db:"reaction_id"`
	ReagentId           int64   `json:"reagentId,omitempty" db:"reagent_id"`
	ReagentName         string  `json:"reagentName,omitempty" db:"reagent_name"`
	SampleStatus        string  `json:"sampleStatus,omitempty" db:"sample_status"`
	IsTestSamples       int     `json:"isTestSamples" db:"is_test_samples"`
	SampleType          string  `json:"sampleType,omitempty" db:"sample_type"`
	SampleId            string  `json:"sampleId,omitempty" db:"sample_id"`
	Mf                  string  `json:"mf,omitempty" db:"mf"`
	Mass                float64 `json:"mass,omitempty" db:"mass"`
	Purity              float64 `json:"purity,omitempty" db:"purity"`
	Synthesised         float64 `json:"synthesised,omitempty" db:"synthesised"`
	SynthesisedUnit     string  `json:"synthesisedUnit,omitempty" db:"synthesised_unit"`
	SampleYield         float64 `json:"sampleYield,omitempty" db:"sample_yield"`
	SampleReference     string  `json:"sampleReference,omitempty" db:"sample_reference"`
	Moles               float64 `json:"moles,omitempty" db:"moles"`
	MolesUnit           string  `json:"molesUnit,omitempty" db:"moles_unit"`
	UseInYield          int     `json:"useInYield,omitempty" db:"use_in_yield"`
	AmountSubmitted     float64 `json:"amountSubmitted,omitempty" db:"amount_submitted"`
	AmountSubmittedUnit string  `json:"amountSubmittedUnit,omitempty" db:"amount_submitted_unit"`
	Barcode             string  `json:"barcode,omitempty" db:"barcode"`
	EnantiomericPurity  float64 `json:"enantiomericPurity,omitempty" db:"enantiomeric_purity"`
	ProductName         string  `json:"productName,omitempty" db:"product_name"`
	QaReason            string  `json:"qaReason,omitempty" db:"qa_reason"`
	ProductAlias        string  `json:"productAlias" db:"product_alias"`
	Color               string  `json:"color,omitempty" db:"color"`
	Qualifier           string  `json:"qualifier,omitempty" db:"qualifier"`
}

type UsedRows struct {
	Compound UsedProps `json:"compound"`
	Cas      []string  `json:"cas"`
}
type UsedCmpd struct {
	TotalCount int        `json:"totalCount"`
	Rows       []UsedRows `json:"rows"`
}

type ElnReactionNote struct {
	ReactionId       int64   `json:"reactionId,omitempty" db:"reaction_id"`
	ReactionSmiles   string  `json:"reactionSmiles,omitempty" db:"reaction_smiles"`
	DaylightSmiles   string  `json:"daylightSmiles,omitempty" db:"daylight_smiles"`
	CdStructure      string  `json:"cdStructure,omitempty" db:"cd_structure"`
	Temperature      float32 `json:"temperature,omitempty" db:"temperature"`
	RxnMolarity      float32 `json:"rxnMolarity,omitempty" db:"rxn_molarity"`
	Pressure         string  `json:"pressure,omitempty" db:"pressure"`
	Time             float32 `json:"time,omitempty" db:"time"`
	GetTargetProduct int8    `json:"getTargetProduct,omitempty" db:"get_target_product"` // 是否获取到目标产物
	Comments         string  `json:"comments,omitempty" db:"comments"`
	ProcedureTxt     string  `json:"procedureTxt,omitempty" db:"procedure_txt"`
	ProcedureHtml    string  `json:"procedureHtml,omitempty" db:"procedure_html"`
}

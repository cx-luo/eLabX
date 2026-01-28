// Package enote coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/6 14:10
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction.go
// @Software: GoLand
package enote

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"time"

	"github.com/cx-luo/go-indigo/core"
	"github.com/gin-gonic/gin"
)

func parseMolecule(indigoInit *core.Indigo, indigoInchi *core.IndigoInchi, molHandle int) (types.Molecule, error) {
	molMolecule, err := indigoInit.LoadMoleculeFromHandle(molHandle)
	if err != nil {
		return types.Molecule{}, err
	}
	defer molMolecule.Close()
	molecularWeight, _ := molMolecule.MolecularWeight()
	molecularWeight = float64(int(molecularWeight*1000+0.5)) / 1000
	cxSmiles, _ := molMolecule.ToCXSmiles()
	smiles, _ := molMolecule.ToSmiles()
	formula, _ := molMolecule.MolecularFormula()
	inchi, _ := indigoInchi.GenerateInChI(molMolecule)
	inchiKey, _ := indigoInchi.InChIToKey(inchi)
	mostAbundantMass, _ := molMolecule.MostAbundantMass()
	monoisotopicMass, _ := molMolecule.MonoisotopicMass()
	massComposition, _ := molMolecule.MassComposition()
	tpsa, _ := molMolecule.TPSA(false)
	numRotatableBonds, _ := molMolecule.NumRotatableBonds()
	numHeavyAtoms, _ := molMolecule.CountHeavyAtoms()
	numBonds, _ := molMolecule.CountBonds()
	numAtoms, _ := molMolecule.CountAtoms()
	return types.Molecule{
		Name:              smiles,
		Formula:           formula,
		InChi:             inchi,
		InChiKey:          inchiKey,
		MolecularWeight:   molecularWeight,
		CxSmiles:          cxSmiles,
		MostAbundantMass:  mostAbundantMass,
		MonoisotopicMass:  monoisotopicMass,
		MassComposition:   massComposition,
		TPSA:              tpsa,
		NumRotatableBonds: numRotatableBonds,
		NumHeavyAtoms:     numHeavyAtoms,
		NumAtoms:          numAtoms,
		NumBonds:          numBonds,
	}, nil
}

func parseReaction(rxnSmiles string) (*types.Reaction, error) {
	reaction := &types.Reaction{
		ReactionId:     utils.GenerateSnowflakeID(),
		ReactionSmiles: rxnSmiles,
		CdStructure:    "",
		CxSmiles:       "",
		ImageSvg:       []byte(""),
		Reactants:      make([]types.Molecule, 0),
		Products:       make([]types.Molecule, 0),
		Reagents:       make([]types.Molecule, 0),
	}

	indigoInit, err := core.IndigoInit()
	if err != nil {
		return nil, err
	}

	indigoInchi, err := indigoInit.InchiInit()
	if err != nil {
		return nil, err
	}
	defer indigoInchi.InchiDispose()

	indigoReaction, err := indigoInit.LoadReactionFromString(rxnSmiles)
	if err != nil {
		return nil, err
	}

	cdStructure, err := indigoReaction.ToCML()
	if err != nil {
		return nil, err
	}

	reaction.CdStructure = cdStructure
	cxSmiles, err := indigoReaction.ToCXSmiles()
	if err != nil {
		return nil, err
	}

	reaction.CxSmiles = cxSmiles

	imageSvg, err := RenderReaction(indigoInit, rxnSmiles)
	if err != nil {
		return nil, err
	}

	reaction.ImageSvg = imageSvg

	reactants, err := indigoReaction.GetAllReactants()
	if err != nil {
		return nil, err
	}

	products, err := indigoReaction.GetAllProducts()
	if err != nil {
		return nil, err
	}

	reagents, err := indigoReaction.GetAllCatalysts()
	if err != nil {
		return nil, err
	}

	for _, reactant := range reactants {
		molecule, err := parseMolecule(indigoInit, indigoInchi, reactant)
		if err != nil {
			return nil, err
		}
		reaction.Reactants = append(reaction.Reactants, molecule)
	}

	for _, product := range products {
		molecule, err := parseMolecule(indigoInit, indigoInchi, product)
		if err != nil {
			return nil, err
		}
		reaction.Products = append(reaction.Products, molecule)
	}

	for _, reagent := range reagents {
		molecule, err := parseMolecule(indigoInit, indigoInchi, reagent)
		if err != nil {
			return nil, err
		}
		reaction.Reagents = append(reaction.Reagents, molecule)
	}

	return reaction, nil
}

func SaveRxnToServer(c *gin.Context) {
	var param struct {
		RxnSmiles string `json:"rxnSmiles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	rxn, err := parseReaction(param.RxnSmiles)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	now := time.Now()
	elnRxn := &types.ElnReaction{
		ReactionId:     rxn.ReactionId,
		ReactionSmiles: rxn.ReactionSmiles,
		CxSmiles:       rxn.CxSmiles,
		CdStructure:    rxn.CdStructure,
		GmtCreate:      now,
		GmtModified:    now,
	}

	if err := dao.OBCursor.Model(&types.ElnReaction{}).Create(elnRxn).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	indigoInit, err := core.IndigoInit()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	type moleculeRole struct {
		Mols []types.Molecule
		Role string
	}

	roles := []moleculeRole{
		{rxn.Reactants, "reactant"},
		{rxn.Products, "product"},
		{rxn.Reagents, "reagent"},
	}

	var compounds []*types.ReagentTableDataStruct

	for _, role := range roles {
		for _, mol := range role.Mols {
			reagentId := utils.GenerateSnowflakeID()
			elnReagent := &types.ElnRxnReagents{
				ReagentId:        reagentId,
				ReactionId:       rxn.ReactionId,
				ReagentName:      mol.Name,
				ReagentSmiles:    mol.CxSmiles,
				Mw:               mol.MolecularWeight,
				MonoisotopicMass: mol.MonoisotopicMass,
				Formula:          mol.Formula,
				ReagentRole:      role.Role,
				Cxsmiles:         mol.CxSmiles,
				Inchi:            mol.InChi,
				Inchikey:         mol.InChiKey,
				GmtCreate:        now,
				GmtModified:      now,
			}
			if err := dao.OBCursor.Model(&types.ElnRxnReagents{}).Create(elnReagent).Error; err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			reagentImg, err := RenderCompound(indigoInit, mol.CxSmiles)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			compounds = append(compounds, &types.ReagentTableDataStruct{
				ReagentId:     elnReagent.ReagentId,
				ReagentName:   elnReagent.ReagentName,
				ReagentSmiles: elnReagent.ReagentSmiles,
				Mw:            elnReagent.Mw,
				ReagentRole:   elnReagent.ReagentRole,
				Formula:       elnReagent.Formula,
				ReagentImg:    reagentImg,
			})
		}
	}

	utils.SuccessWithData(c, "success", map[string]interface{}{
		"reactionId": rxn.ReactionId,
		"imageSvg":   rxn.ImageSvg,
		"compounds":  compounds,
	})
}

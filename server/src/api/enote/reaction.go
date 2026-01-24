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

	rxnSmiles := param.RxnSmiles
	rxn, err := parseReaction(rxnSmiles)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Convert types.Reaction to types.ElnReaction
	now := time.Now()
	elnRxn := &types.ElnReaction{
		ReactionId:     rxn.ReactionId,
		ReactionSmiles: rxn.ReactionSmiles,
		CxSmiles:       rxn.CxSmiles,
		CdStructure:    rxn.CdStructure,
		GmtCreate:      now,
		GmtModified:    now,
	}

	err = dao.OBCursor.Model(&types.ElnReaction{}).Create(elnRxn).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Save reactants
	for _, reactant := range rxn.Reactants {
		reagentId := utils.GenerateSnowflakeID()
		elnReagent := &types.ElnRxnReagents{
			ReagentId:        reagentId,
			ReactionId:       rxn.ReactionId,
			ReagentName:      reactant.Name,
			ReagentSmiles:    reactant.Name, // Name is SMILES
			Mw:               float32(reactant.MolecularWeight),
			MonoisotopicMass: float32(reactant.MonoisotopicMass),
			Formula:          reactant.Formula,
			ReagentRole:      "reactant",
			Cxsmiles:         reactant.CxSmiles,
			Inchi:            reactant.InChi,
			Inchikey:         reactant.InChiKey,
			GmtCreate:        now,
			GmtModified:      now,
		}
		err = dao.OBCursor.Model(&types.ElnRxnReagents{}).Create(elnReagent).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	// Save products
	for _, product := range rxn.Products {
		reagentId := utils.GenerateSnowflakeID()
		elnReagent := &types.ElnRxnReagents{
			ReagentId:        reagentId,
			ReactionId:       rxn.ReactionId,
			ReagentName:      product.Name,
			ReagentSmiles:    product.Name, // Name is SMILES
			Mw:               float32(product.MolecularWeight),
			MonoisotopicMass: float32(product.MonoisotopicMass),
			Formula:          product.Formula,
			ReagentRole:      "product",
			Cxsmiles:         product.CxSmiles,
			Inchi:            product.InChi,
			Inchikey:         product.InChiKey,
			GmtCreate:        now,
			GmtModified:      now,
		}
		err = dao.OBCursor.Model(&types.ElnRxnReagents{}).Create(elnReagent).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	// Save reagents
	for _, reagent := range rxn.Reagents {
		reagentId := utils.GenerateSnowflakeID()
		elnReagent := &types.ElnRxnReagents{
			ReagentId:        reagentId,
			ReactionId:       rxn.ReactionId,
			ReagentName:      reagent.Name,
			ReagentSmiles:    reagent.Name, // Name is SMILES
			Mw:               float32(reagent.MolecularWeight),
			MonoisotopicMass: float32(reagent.MonoisotopicMass),
			Formula:          reagent.Formula,
			ReagentRole:      "reagent",
			Cxsmiles:         reagent.CxSmiles,
			Inchi:            reagent.InChi,
			Inchikey:         reagent.InChiKey,
			GmtCreate:        now,
			GmtModified:      now,
		}
		err = dao.OBCursor.Model(&types.ElnRxnReagents{}).Create(elnReagent).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	utils.SuccessWithData(c, "success", map[string]interface{}{
		"reactionId": rxn.ReactionId,
		"imageSvg":   rxn.ImageSvg,
	})
}

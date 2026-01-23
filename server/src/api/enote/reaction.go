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
	"strconv"

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
		ExactMass:         monoisotopicMass,
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

	err = dao.OBCursor.Model(&types.ElnReaction{}).Create(rxn).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	for _, reactant := range rxn.Reactants {
		err = dao.OBCursor.Model(&types.ElnRxnReagents{ReagentRole: "reactant"}).Create(reactant).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}
	for _, product := range rxn.Products {
		err = dao.OBCursor.Model(&types.ElnRxnReagents{ReagentRole: "product"}).Create(product).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}
	for _, reagent := range rxn.Reagents {
		err = dao.OBCursor.Model(&types.ElnRxnReagents{ReagentRole: "reagent"}).Create(reagent).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	utils.Success(c, strconv.FormatInt(rxn.ReactionId, 10))
}

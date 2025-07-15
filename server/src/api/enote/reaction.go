// Package enote coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/6 14:10
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction.go
// @Software: GoLand
package enote

import (
	"eLabX/src/types"
	"regexp"
	"strconv"
	"strings"
)

func parseReaction(data map[string]string) ([]types.Molecule, error) {
	re := regexp.MustCompile(`$$(.*?)$$`)
	molecules := make([]types.Molecule, 0)

	// 解析 gross 字段，获取所有分子名
	grossGroups := re.FindAllStringSubmatch(data["gross"], -1)
	var allGross []string
	for _, group := range grossGroups {
		parts := strings.Split(group[1], ";")
		for _, part := range parts {
			cleaned := strings.TrimSpace(part)
			if cleaned != "" {
				allGross = append(allGross, cleaned)
			}
		}
	}

	// 解析 molecular-weight
	mwGroups := re.FindAllStringSubmatch(data["molecular-weight"], -1)
	var allMW []float64
	for _, group := range mwGroups {
		parts := strings.Split(group[1], ";")
		for _, part := range parts {
			val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				return nil, err
			}
			allMW = append(allMW, val)
		}
	}

	// 解析 most-abundant-mass
	mamGroups := re.FindAllStringSubmatch(data["most-abundant-mass"], -1)
	var allMAM []float64
	for _, group := range mamGroups {
		parts := strings.Split(group[1], ";")
		for _, part := range parts {
			val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				return nil, err
			}
			allMAM = append(allMAM, val)
		}
	}

	// 合并结果
	for i := 0; i < len(allGross); i++ {
		mol := types.Molecule{
			Formula:          allGross[i],
			MoleWeight:       allMW[i],
			MostAbundantMass: allMAM[i],
		}
		molecules = append(molecules, mol)
	}

	return molecules, nil
}

// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/5/7 14:43
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : pubchem.go
// @Software: GoLand
package api

import (
	"eLabX/src/common"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Compound struct {
	Cid                    int     `json:"cid,omitempty"`
	Mw                     float64 `json:"mw,omitempty"`
	Polararea              float64 `json:"polararea,omitempty"`
	Complexity             float64 `json:"complexity,omitempty"`
	Xlogp                  float64 `json:"xlogp,omitempty"`
	Heavycnt               int     `json:"heavycnt,omitempty"`
	Hbonddonor             int     `json:"hbonddonor,omitempty"`
	Hbondacc               int     `json:"hbondacc,omitempty"`
	Rotbonds               int     `json:"rotbonds,omitempty"`
	Annothitcnt            int     `json:"annothitcnt,omitempty"`
	Charge                 int     `json:"charge,omitempty"`
	Covalentunitcnt        int     `json:"covalentunitcnt,omitempty"`
	Isotopeatomcnt         int     `json:"isotopeatomcnt,omitempty"`
	Totalatomstereocnt     int     `json:"totalatomstereocnt,omitempty"`
	Definedatomstereocnt   int     `json:"definedatomstereocnt,omitempty"`
	Undefinedatomstereocnt int     `json:"undefinedatomstereocnt,omitempty"`
	Totalbondstereocnt     int     `json:"totalbondstereocnt,omitempty"`
	Definedbondstereocnt   int     `json:"definedbondstereocnt,omitempty"`
	Undefinedbondstereocnt int     `json:"undefinedbondstereocnt,omitempty"`
	Pclidcnt               int     `json:"pclidcnt,omitempty"`
	Gpidcnt                int     `json:"gpidcnt,omitempty"`
	Gpfamilycnt            int     `json:"gpfamilycnt,omitempty"`
	Aids                   string  `json:"aids,omitempty"`
	Cmpdname               string  `json:"cmpdname,omitempty"`
	Cmpdsynonym            string  `json:"cmpdsynonym,omitempty"`
	Inchi                  string  `json:"inchi,omitempty"`
	Inchikey               string  `json:"inchikey,omitempty"`
	Smiles                 string  `json:"smiles,omitempty"`
	Iupacname              string  `json:"iupacname,omitempty"`
	Mf                     string  `json:"mf,omitempty"`
	Sidsrcname             string  `json:"sidsrcname,omitempty"`
	Annotation             string  `json:"annotation,omitempty"`
	Cidcdate               string  `json:"cidcdate,omitempty"`
	Depcatg                string  `json:"depcatg,omitempty"`
	Meshheadings           string  `json:"meshheadings,omitempty"`
	Annothits              string  `json:"annothits,omitempty"`
	Exactmass              string  `json:"exactmass,omitempty"`
	Monoisotopicmass       string  `json:"monoisotopicmass,omitempty"`
}

type SDQSet struct {
	Status struct {
		Code  int    `json:"code,omitempty"`
		Error string `json:"error,omitempty"`
	} `json:"status"`
	InputCount int        `json:"inputCount,omitempty"`
	TotalCount int        `json:"totalCount,omitempty"`
	Collection string     `json:"collection,omitempty"`
	Type       string     `json:"type,omitempty"`
	Rows       []Compound `json:"rows,omitempty"`
}

// SDQOutputSet for search from cachekey
type SDQOutputSet struct {
	SDQOutputSet []SDQSet `json:"SDQOutputSet"`
}

func urlGet(epUrl string) ([]byte, error) {
	// 创建一个自定义的 Transport 实例
	transport := &http.Transport{
		//Proxy: func(req *http.Request) (*url.URL, error) {
		//	return url.Parse("http://27.79.147.195:4005") // 设置代理
		//},
		MaxIdleConnsPerHost: 5,  // 每个主机最大空闲连接数
		MaxIdleConns:        20, // 最大空闲连接数
	}

	// 创建一个自定义的 Client 实例
	client := &http.Client{
		Transport: transport,        // 设置 Transport
		Timeout:   time.Second * 30, // 设置超时
	}

	req, err := http.NewRequest("GET", epUrl, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Error status : %s, url : %s ", response.Status, epUrl))
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	bodyText, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyText, nil
}

func (s *SDQSet) InsertCompoundsToDB() error {
	insertSql := `insert INTO elabx.compound_from_pubchem(
		cid, mw, polararea, complexity, exactmass,
		monoisotopicmass, heavycnt, hbonddonor, hbondacc, 
		rotbonds, annothitcnt, charge, covalentunitcnt, 
		isotopeatomcnt, totalatomstereocnt, definedatomstereocnt, 
		undefinedatomstereocnt,
		totalbondstereocnt,
		definedbondstereocnt,
		undefinedbondstereocnt,
		pclidcnt,
		gpidcnt,
		gpfamilycnt,
		cmpdname,
        cmpdsynonym,
		inchi,
		inchikey,
		isosmiles,
		iupacname,
		mf,
		sidsrcname,
		cidcdate,
		depcatg,
		annothits)
		VALUES (:cid,
		:mw,
		:polararea,
		:complexity,
		:exactmass,
		:monoisotopicmass,
		:heavycnt,
		:hbonddonor,
		:hbondacc,
		:rotbonds,
		:annothitcnt,
		:charge,
		:covalentunitcnt,
		:isotopeatomcnt,
		:totalatomstereocnt,
		:definedatomstereocnt,
		:undefinedatomstereocnt,
		:totalbondstereocnt,
		:definedbondstereocnt,
		:undefinedbondstereocnt,
		:pclidcnt,
		:gpidcnt,
		:gpfamilycnt,
		:cmpdname,
		:cmpdsynonym,
		:inchi,
		:inchikey,
		:smiles,
		:iupacname,
		:mf,
		:sidsrcname,
		:cidcdate,
		:depcatg,
		:annothits) on duplicate key update cmpdsynonym=:cmpdsynonym, isosmiles = :smiles`

	for i := 0; i < len(s.Rows); i++ {
		row := s.Rows[i]
		sema.Acquire(1)
		go func() {
			defer sema.Release()
			err := func(c Compound) error {
				_, err := dao.OBCursor.NamedExec(insertSql, c)
				if err != nil {
					return err
				}
				return nil
			}(row)
			if err != nil {
				utils.Logger.Error(err.Error())
			}
		}()
	}
	sema.Wait()
	return nil
}

func InsertSDQToDB(s *[]SDQSet) error {
	for i := 0; i < len(*s); i++ {
		err := (*s)[i].InsertCompoundsToDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSDQOutputSetFromCid(cid int, limit int, start int) SDQOutputSet {
	jsData := fmt.Sprintf(`{"select":"*","collection":"compound","where":{"ands":[{"cid":"%d"}]},"order":["cid,asc"],"start":%d,"limit":%d,"width":1000000,"listids":0}`, cid, start, limit)
	currUrl := "https://pubchem.ncbi.nlm.nih.gov/sdq/sdqagent.cgi?"
	params := url.Values{}
	params.Set("infmt", "json")
	params.Set("outfmt", "json")
	params.Set("query", jsData)
	currUrl = currUrl + params.Encode()
	utils.Logger.Info(currUrl)

	var sdq SDQOutputSet
	bodyText, err := urlGet(currUrl)
	if err != nil {
		utils.Logger.Warn(err.Error())
		return SDQOutputSet{SDQOutputSet: nil}
	}
	err = json.Unmarshal(bodyText, &sdq)
	if err != nil {
		utils.Logger.Error(err.Error())
		return SDQOutputSet{SDQOutputSet: nil}
	}

	if sdq.SDQOutputSet[0].Status.Code != 0 {
		utils.Logger.Error(fmt.Sprintf("GetSDQOutputSetFromCid : %d, %d, %s", cid, sdq.SDQOutputSet[0].Status.Code, sdq.SDQOutputSet[0].Status.Error))
		return SDQOutputSet{SDQOutputSet: nil}
	}
	return sdq
}

func GetCIDFromSmiles(smiles string) []int {
	//curl := fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/smiles/%s/cids/json", smiles)
	params := url.Values{}
	params.Set("smiles", smiles)
	curl := "https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/smiles/cids/json?" + params.Encode()
	res, err := urlGet(curl)

	var idenCids struct {
		IdentifierList struct {
			Cid []int `json:"CID"`
		} `json:"IdentifierList"`
	}

	if err != nil {
		utils.Logger.Warn("res is nil or " + err.Error())
		return nil
	}
	err = json.Unmarshal(res, &idenCids)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error on Unmarshal:  %s, %s, %s", curl, smiles, err))
		return nil
	}
	if len(idenCids.IdentifierList.Cid) >= 1 && idenCids.IdentifierList.Cid[0] != 0 {
		return idenCids.IdentifierList.Cid
	}
	return nil
}

func GetCIDFromName(name string) []int {
	type CIDs struct {
		ConceptsAndCIDs struct {
			CID []int `json:"CID" gorm:"column:CID"`
		} `json:"ConceptsAndCIDs" gorm:"column:ConceptsAndCIDs"`
	}
	curl := "https://pubchem.ncbi.nlm.nih.gov/rest/pug/concepts/name/JSON?"

	params := url.Values{}
	params.Set("name", name)
	curl = curl + params.Encode()

	var cIds CIDs
	res, err := urlGet(curl)
	if res == nil {
		utils.Logger.Warn("The res is nil or " + err.Error())
		return nil
	}

	err = json.Unmarshal(res, &cIds)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s, %s, %s", curl, name, err))
		return nil
	}

	//for _, cid := range cIds.ConceptsAndCIDs.CID {
	//	sqdSet := GetSDQOutputSetFromCid(cid, 10, 1).SDQOutputSet
	//	err = InsertSDQToDB(&sqdSet)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	if len(cIds.ConceptsAndCIDs.CID) >= 1 && cIds.ConceptsAndCIDs.CID[0] != 0 {
		return cIds.ConceptsAndCIDs.CID
	}

	return nil
}

func GetSDQOutputSetFromQuery(cName string, limit int, start int) SDQOutputSet {
	jsData := fmt.Sprintf(`{"select":"*","collection":"compound","where":{"ands":[{"*":"%s"}]},"order":["relevancescore,desc"],"start":%d,"limit":%d,"width":1000000,"listids":0}
`, cName, start, limit)
	currUrl := "https://pubchem.ncbi.nlm.nih.gov/sdq/sdqagent.cgi?"
	params := url.Values{}
	params.Set("infmt", "json")
	params.Set("outfmt", "json")
	params.Set("query", jsData)
	currUrl = currUrl + params.Encode()

	var sdq SDQOutputSet
	bodyText, err := urlGet(currUrl)
	if err != nil {
		utils.Logger.Warn(err.Error())
		return SDQOutputSet{SDQOutputSet: nil}
	}

	err = json.Unmarshal(bodyText, &sdq)
	if err != nil {
		utils.Logger.Error(err.Error())
		return SDQOutputSet{SDQOutputSet: nil}
	}
	return sdq
}

func getCasByRegexp(s string) []string {
	casRegex := regexp.MustCompile(`\b\d{2,7}-\d{2}-\d\b`)
	casNumbers := casRegex.FindAllString(s, -1)
	var validCasNumbers []string
	for _, casNumber := range casNumbers {
		if calculateChecksum(casNumber) {
			validCasNumbers = append(validCasNumbers, casNumber)
		}
	}
	return validCasNumbers
}

func GetCmpdInfoBySmiles(c *gin.Context) {
	var r struct {
		Smiles string `json:"smiles"`
	}
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var compounds common.UsedCmpd

	cids := GetCIDFromSmiles(r.Smiles)
	if len(cids) == 0 {
		compounds.TotalCount = 0
		c.JSON(http.StatusOK, utils.BaseResponse{
			StatusCode: 200, Msg: "not Found, please try a search from other ways", Data: gin.H{"compounds": compounds},
		})
		return
	}

	for _, cid := range cids {
		sqdSet := GetSDQOutputSetFromCid(cid, 10, 1).SDQOutputSet
		err = InsertSDQToDB(&sqdSet)
		if err != nil {
			utils.Logger.Error(err.Error())
			utils.InternalRequestErr(c, err)
			return
		}
		compounds.TotalCount += sqdSet[0].TotalCount
		for _, row := range sqdSet[0].Rows {
			cas := getCasByRegexp(row.Cmpdsynonym)
			var u = common.UsedRows{
				Compound: common.UsedProps{
					Cid: row.Cid,
					Mf:  row.Mf,
					Mw:  row.Mw,
					//Exactmass:        row.Exactmass,
					//Monoisotopicmass: row.Monoisotopicmass,
					Cmpdname:  row.Cmpdname,
					Inchi:     row.Inchi,
					Inchikey:  row.Inchikey,
					Smiles:    row.Smiles,
					Iupacname: row.Iupacname,
				},
				Cas: cas,
			}
			compounds.Rows = append(compounds.Rows, u)
		}
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{"compounds": compounds},
	})

	return
}

func GetCmpdInfoByName(c *gin.Context) {
	var r struct {
		Name string `json:"name"`
	}
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var compounds common.UsedCmpd
	cids := GetCIDFromName(r.Name)
	if len(cids) == 0 {
		utils.BadRequestErr(c, errors.New("not Found, please try a search from other ways"))
		return
	}

	for _, cid := range cids {
		sqdSet := GetSDQOutputSetFromCid(cid, 10, 1).SDQOutputSet
		err = InsertSDQToDB(&sqdSet)
		if err != nil {
			utils.Logger.Error(err.Error())
			utils.InternalRequestErr(c, err)
			return
		}
		compounds.TotalCount += sqdSet[0].TotalCount
		for _, row := range sqdSet[0].Rows {
			cas := getCasByRegexp(row.Cmpdsynonym)
			var u = common.UsedRows{
				Compound: common.UsedProps{
					Cid: row.Cid,
					Mf:  row.Mf,
					Mw:  row.Mw,
					//Exactmass:        row.Exactmass,
					//Monoisotopicmass: row.Monoisotopicmass,
					Cmpdname:  row.Cmpdname,
					Inchi:     row.Inchi,
					Inchikey:  row.Inchikey,
					Smiles:    row.Smiles,
					Iupacname: row.Iupacname,
				},
				Cas: cas,
			}
			compounds.Rows = append(compounds.Rows, u)
		}
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{"compounds": compounds},
	})

	return
}

func processAndWriteToRedis(searchKey string) (*common.UsedCmpd, error) {
	var compounds common.UsedCmpd
	var sqdSet []SDQSet
	//CIds := GetCIDFromSmiles(searchKey)
	//CIds = append(CIds, GetCIDFromName(searchKey)...)
	CIds := GetCIDFromName(searchKey)

	compounds.TotalCount += len(CIds)
	for _, cid := range CIds {
		sqdSet = append(sqdSet, GetSDQOutputSetFromCid(cid, 10, 1).SDQOutputSet...)
	}

	sqdSet = append(sqdSet, GetSDQOutputSetFromQuery(searchKey, 10, 1).SDQOutputSet...)
	err := InsertSDQToDB(&sqdSet)
	if err != nil {
		return nil, err
	}

	for _, row := range sqdSet[0].Rows {
		cas := getCasByRegexp(row.Cmpdsynonym)
		var u = common.UsedRows{
			Compound: common.UsedProps{
				//Exactmass:        row.Exactmass,
				//Monoisotopicmass: row.Monoisotopicmass,
				Cid:       row.Cid,
				Mf:        row.Mf,
				Mw:        row.Mw,
				Cmpdname:  row.Cmpdname,
				Inchi:     row.Inchi,
				Inchikey:  row.Inchikey,
				Smiles:    row.Smiles,
				Iupacname: row.Iupacname,
			},
			Cas: cas,
		}
		compounds.Rows = append(compounds.Rows, u)
	}

	compounds.TotalCount = len(compounds.Rows)
	err = dao.WriteSearchResult(searchKey, compounds)
	if err != nil {
		return nil, err
	}

	go func() {
		err := dao.CleanUpKeys(dao.RedisClient)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("clean data error: %s.", err.Error()))
		}
	}()

	return &compounds, nil
}

func CmpdComboSearch(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		utils.Logger.Info("CmpdComboSearch completed", zap.Duration("duration", duration))
	}()
	var q struct {
		Query string `json:"query"`
	}
	err := c.ShouldBind(&q)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	searchKey := strings.TrimSpace(q.Query)

	err, resFromRedis := dao.GetRSearchResult(searchKey)
	if err != nil {
		if errors.Is(err, dao.ErrKeyNotExist) || errors.Is(err, dao.ErrValueIsNull) {
			compounds, err := processAndWriteToRedis(searchKey)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			c.JSON(http.StatusOK, utils.BaseResponse{
				StatusCode: 200, Msg: "", Data: gin.H{"compounds": compounds},
			})
			return
		}
		utils.InternalRequestErr(c, err)
		return
	} else {
		var compounds common.UsedCmpd
		err = json.Unmarshal([]byte(resFromRedis), &compounds)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		utils.SuccessWithData(c, "", gin.H{"compounds": compounds})
		return
	}
}

func calculateChecksum(cas string) bool {
	casRegex := regexp.MustCompile(`(\d{2,7})-(\d{2})-(\d)`)
	match := casRegex.FindStringSubmatch(cas)
	// 将part1和part2拼接起来
	number := []rune(match[1] + match[2])

	// 初始化校验码计算的和
	sum := 0

	// 从最低位开始计算，即从字符串的末尾开始
	for i, r := range number {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			utils.Logger.Error(err.Error())
			return false
		}
		sum += (len(number) - i) * digit
	}

	// 计算最终的校验码
	checksum := sum % 10
	atoi, err := strconv.Atoi(match[3])
	if err != nil {
		utils.Logger.Error(err.Error())
		return false
	}

	return checksum == atoi
}

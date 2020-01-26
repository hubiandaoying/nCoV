package ncov

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// APIROOT API的根目录
	APIROOT = "/2019-nCoV"
)

// MainStatus 病毒状态对象
type MainStatus struct {
	name         string  // 病毒名称
	url          string  // 请求地址
	ConfirmCount float64 `json:"confirmCount"` // 确诊数量
	SuspectCount float64 `json:"suspectCount"` // 疑患数量
	DeadCount    float64 `json:"deadCount"`    // 死亡数量
	HintWords    string  `json:"hintWords"`    // 当日关键词
	RecentTime   string  `json:"recentTime"`   // 更新时间
	Cure         float64 `json:"cure"`         // 治愈数量
	UseTotal     bool    `json:"useTotal"`
}

// GetString 返回字符串信息
func (v *MainStatus) GetString() string {
	var str string
	format := "[%s][%s] 已经确诊感染%0.f人，疑患病%0.f人，目前已死亡%0.f人、治愈%0.f人,当前提示：%s"
	str = fmt.Sprintf(format, v.RecentTime, v.name, v.ConfirmCount, v.SuspectCount, v.DeadCount, v.Cure, v.HintWords)
	return str
}

// GetVirusStatusJSON 返回生成病毒状态的JSON
func (v *MainStatus) GetVirusStatusJSON() []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(v)
	dealErr(err)
	return jsonData
}

// GetVirusStatus 获取病毒状态返回实例化对象
func GetVirusStatus() *MainStatus {
	v := &MainStatus{}
	v.name = "新型冠状病毒 2019-nCoV"
	v.url = "https://view.inews.qq.com/g2/getOnsInfo?name=wuwei_ww_global_vars"

	// 获取数据
	data := requestData(v.url)

	// 实例化JSON对象
	pp := parseTencentJSON(data)

	// 赋值
	if pp[0]["confirmCount"] != nil {
		v.ConfirmCount = pp[0]["confirmCount"].(float64)
	}
	if pp[0]["suspectCount"] != nil {
		v.SuspectCount = pp[0]["suspectCount"].(float64)
	}
	if pp[0]["deadCount"] != nil {
		v.DeadCount = pp[0]["deadCount"].(float64)
	}
	if pp[0]["useTotal"] != nil {
		//v.UseTotal = pp[0]["useTotal"].(bool)
	}
	if pp[0]["hintWords"] != nil {
		v.HintWords = pp[0]["hintWords"].(string)
	}
	if pp[0]["recentTime"] != nil {
		v.RecentTime = pp[0]["recentTime"].(string)
	}
	if pp[0]["cure"] != nil {
		v.Cure = pp[0]["cure"].(float64)
	}

	return v
}

// AreaStatus 区域或者国家信息
type AreaStatus struct {
	/*{    "day": "1.12",    "time": "",    "country": "中国",    "area": "湖北武汉",    "dead": 1,    "confirm": 41,    "suspect": 0,    "heal": null,    "city": null,    "district": null  }*/
	Count int                      `json:"count"`
	Datas []map[string]interface{} `json:"data"`
	Areas []Area                   `json:"areas"`
}

// GetAreaStatus 根据名称获取单个区域的状态
func (a *AreaStatus) GetAreaStatus(name string) Area {
	var area Area
	for _, v := range a.Areas {
		if v.Area == name {
			area = v
		}
	}
	return area
}

// Area 区域信息
type Area struct {
	Area         string  `json:"area"`
	ConfirmCount float64 `json:"confirmCount"` // 确诊数量
	SuspectCount float64 `json:"suspectCount"` // 疑患数量
	DeadCount    float64 `json:"deadCount"`    // 死亡数量
	HealCount    float64 `json:"heal"`         // 治愈数量
}

// GetAllAreaStatus 获取区域的病毒状态返回实例化对象
func GetAllAreaStatus() *AreaStatus {
	as := &AreaStatus{}
	url := "https://view.inews.qq.com/g2/getOnsInfo?name=wuwei_ww_area_datas"

	// 获取数据
	data := requestData(url)

	// 解析数据为JSON对象

	jsonObj := parseTencentJSON(data)
	as.Datas = jsonObj
	as.Count = len(jsonObj)

	areas := generateAreas(jsonObj)
	for i, v := range areas {
		maxComfirm := 0.0
		maxSuspect := 0.0
		maxHeal := 0.0
		maxDead := 0.0
		for _, vv := range jsonObj {
			if vv["area"].(string) == v.Area || vv["country"].(string) == v.Area {

				if vv["confirm"] != nil {
					if vv["confirm"].(float64) > maxComfirm {
						maxComfirm = vv["confirm"].(float64)
					}
				}

				if vv["suspect"] != nil {
					if vv["suspect"].(float64) > maxSuspect {
						maxSuspect = vv["suspect"].(float64)
					}
				}

				if vv["heal"] != nil {
					if vv["heal"].(float64) > maxHeal {
						maxHeal = vv["heal"].(float64)
					}
				}

				if vv["dead"] != nil {
					if vv["dead"].(float64) > maxDead {
						maxDead = vv["dead"].(float64)
					}
				}
			}
		}
		areas[i].ConfirmCount = maxComfirm
		areas[i].SuspectCount = maxSuspect
		areas[i].HealCount = maxHeal
		areas[i].DeadCount = maxDead

	}
	as.Areas = areas
	return as
}

// GetAllAreaStatusJSON 返回生成病毒区域状态的JSON
func (a *AreaStatus) GetAllAreaStatusJSON() []byte {
	jsonData, err := json.Marshal(a)
	dealErr(err)
	return jsonData
}

// GetAreaStatusJSON 获取区域的JSON
func (a *Area) GetAreaStatusJSON() []byte {
	jsonData, err := json.Marshal(a)
	dealErr(err)
	return jsonData
}

// Dump Dump current data into a file
func Dump(data interface{}) {
	f, err := os.OpenFile("./dump-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND, 0755)
	dealErr(err)
	defer f.Close()
	io.WriteString(f, fmt.Sprintf("%v", data))
}

// DumpSimple Dump current main data into a file
func DumpSimple(data interface{}) {
	f, err := os.OpenFile("./dump-simple-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND, 0755)
	dealErr(err)
	defer f.Close()
	io.WriteString(f, fmt.Sprintf("%v", data))
}

// -----------------------------
// GET 数据
func requestData(url string) []byte {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	dealErr(err)

	// Add Header
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0")
	req.Header.Add("Referer", "https://news.qq.com/zt2020/page/feiyan.htm")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	// Do
	resp, err := client.Do(req)
	dealErr(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	dealErr(err)
	return data
}

// 处理错误信息
func dealErr(err error) {
	if err != nil {
		log.Panicf("[ERR] %s", err.Error())
	}
}

// 处理重复字段
func dealRepeat(ori string, repeat string) string {
	strl := strings.Split(ori, repeat)
	str := ""
	for _, v := range strl {
		str += v
	}
	return str
}

// 解析JSON返回对象
func parseTencentJSON(data []byte) []map[string]interface{} {
	str := dealRepeat(string(data), "\\n")

	// 对返回的数据进行解析
	var raw map[string]interface{}
	dealErr(json.Unmarshal(json.RawMessage(str), &raw))

	var pp []map[string]interface{}
	dealErr(json.Unmarshal(json.RawMessage(raw["data"].(string)), &pp))
	return pp
}

// 生成空的区域
func generateAreas(jsonObj []map[string]interface{}) []Area {
	var areas []Area
	// 生成区域信息
	for _, v := range jsonObj {
		var area Area
		name := ""
		if v["area"] != "" {
			name = v["area"].(string)
		} else if v["country"] != "" {
			name = v["country"].(string)
		}

		if checkAreaExists(areas, name) {
			continue
		}
		area.Area = name

		areas = append(areas, area)
	}
	return areas
}

// 检测区域是否已添加
func checkAreaExists(areas []Area, key string) bool {
	for _, v := range areas {
		if v.Area == key {
			return true
		}
	}
	return false
}

// APIS

// API 服务器接口
type API struct {
	root string
}

// ServeHTTP 提供一个服务器接口
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// APIROOT: /2019-nCoV
	a.root = APIROOT
	u := req.URL.String()
	w.Header().Add("Content-Type", "application/json")

	dumpFile, err := os.OpenFile("./request.log", os.O_CREATE|os.O_APPEND, 0755)
	defer dumpFile.Close()
	dealErr(err)
	_, err = io.WriteString(dumpFile, fmt.Sprintf("[%s]%s\n", time.Now().Format("2006-01-02 15:04:05"), req.RemoteAddr))
	dealErr(err)

	if u == APIROOT+"/status" {
		fmt.Fprintf(w, "%v", string(GetVirusStatus().GetVirusStatusJSON()))
		return
	}

	if strings.Contains(u, APIROOT+"/areas") {
		values := req.URL.Query()
		city := values.Get("area")
		if city == "" {
			fmt.Fprintf(w, "%v", string(GetAllAreaStatus().GetAllAreaStatusJSON()))
			return
		}
		area := GetAllAreaStatus().GetAreaStatus(city)
		fmt.Fprintf(w, "%v", string(area.GetAreaStatusJSON()))
		return
	}
}

// TODO: rebuild cache
package ipctrl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/CodeInEverest/golib/cache"
	"github.com/gwaylib/errors"
)

var defController = NewAccessController(time.Second/2, time.Minute, 20)

// read ip from haproxy or origin ip
func ReadIp(req *http.Request) (ip []string, port string) {
	return defController.ReadIp(req)
}

// 以IP组数为key进行缓存记录调用次数，
// 通过调用的次数及频率检查是否是攻击行为
// 返回1分钟内剩余可用的次数
// 如果任一IP在白名单中，将不进行检查,永返回999次可用
// 如果任一IP在黑名单中，将永返回true,永返回-999次可用
func AttackCheck(ipGroup []string) int {
	return defController.AttackCheck(ipGroup)
}

// 检查是否在黑名单中
func InBlackList(ipGroup []string) (string, bool) {
	return defController.InBlackList(ipGroup)
}

// 检查是否在白名单中
func InWhiteList(ipGroup []string) (string, bool) {
	return defController.InWhiteList(ipGroup)
}

// 通过配置文件重置白名单, 该文件将被热监听
// 文件须写全路径文件地址
func SetWhiteList(fileName string) error {
	return defController.SetWhiteList(fileName)
}
func AddWhiteList(ip string) {
	defController.AddWhiteList(ip)
}
func DelWhiteList(ip string) {
	defController.DelWhiteList(ip)
}

// 正则查找白名单数据
func RegexpWhiteList(partern string) []string {
	return defController.RegexpWhiteList(partern)
}

// 通过配置文件重置黑名单, 该文件将被热监听
// 文件须写全路径文件地址
func SetBlackList(fileName string) error {
	return defController.SetBlackList(fileName)
}
func AddBlackList(ip string) {
	defController.AddBlackList(ip)
}
func DelBlackList(ip string) {
	defController.DelBlackList(ip)
}

// 正则查找黑名单数据
func RegexpBlackList(parttern string) []string {
	return defController.RegexpBlackList(parttern)
}

type NameList map[string]bool

type AccessControl interface {
	// 通过配置文件重置白名单, 该文件将被热监听
	// 文件须写全路径文件地址
	SetWhiteList(fileName string) error
	// 添加一个ip白名单
	AddWhiteList(ip string)
	// 删除一个ip白名单
	DelWhiteList(ip string)
	// 正则查找白名单数据
	RegexpWhiteList(partern string) []string

	// 通过配置文件重置黑名单, 该文件将被热监听
	// 文件须写全路径文件地址
	SetBlackList(fileName string) error
	// 添加一个ip黑名单
	AddBlackList(ip string)
	// 删除一个ip白名单
	DelBlackList(ip string)
	// 正则查找黑名单数据
	RegexpBlackList(parttern string) []string

	// 从请求头中读取ip数据
	// 为兼容haproxy，会读取"X-Forwarded-For“字段；
	// 如果没有，将调用req.RemoteAddr进行读取
	ReadIp(req *http.Request) (ipGroup []string, port string)

	// 读取黑名单数据，如果参数数组中有任一名单在黑名单中，返回真
	InBlackList(ipGroup []string) (ip string, in bool)
	// 读取白名单数据，如果参数数组中有任一名单在白名单中，返回真
	InWhiteList(ipGroup []string) (ip string, in bool)

	// 以IP组数为key进行缓存记录调用次数，
	// 通过调用的次数及频率检查是否是攻击行为
	// 返回1分钟内剩余可用的次数
	// 如果任一IP在白名单中，将不进行检查,永返回999次可用
	// 如果任一IP在黑名单中，将永返回true,永返回-999次可用
	AttackCheck(ipGroup []string) int
}

type AttackCondition struct {
	Times    int
	LastTime time.Time
}

type accessControl struct {
	attackTimesCacher *cache.MemoryCache
	interval          time.Duration
	cacheTime         time.Duration
	limit             int
	whiteList         NameList
	blackList         NameList
}

func NewAccessController(interval, cacheTime time.Duration, limit int) AccessControl {
	control := &accessControl{
		attackTimesCacher: cache.NewMemoryCache(),
		interval:          interval,
		cacheTime:         cacheTime,
		limit:             limit,
		whiteList:         NameList{},
		blackList:         NameList{},
	}
	return control
}

func (ctrl *accessControl) SetWhiteList(fileName string) error {
	// 暂无需考虑并发安全问题
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.As(err, fileName)
	}
	nameList := NameList{}
	if err := json.Unmarshal(data, &nameList); err != nil {
		return errors.As(err, fileName)
	}
	ctrl.whiteList = nameList
	return nil
}

// 添加一个ip白名单
func (ctrl *accessControl) AddWhiteList(ip string) {
	ctrl.whiteList[ip] = true
}
func (ctrl *accessControl) DelWhiteList(ip string) {
	delete(ctrl.whiteList, ip)
}

// 正则查找白名单数据
func (ctrl *accessControl) RegexpWhiteList(pattern string) []string {
	result := []string{}
	for key, _ := range ctrl.whiteList {
		if match, _ := regexp.MatchString(pattern, key); match {
			result = append(result, key)
		}
	}
	return result
}

func (ctrl *accessControl) SetBlackList(fileName string) error {
	// 暂无需考虑并发安全问题
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.As(err, fileName)
	}
	nameList := NameList{}
	if err := json.Unmarshal(data, &nameList); err != nil {
		return errors.As(err, fileName)
	}

	ctrl.blackList = nameList
	return nil
}

// 添加一个ip黑名单
func (ctrl *accessControl) AddBlackList(ip string) {
	ctrl.blackList[ip] = true
}
func (ctrl *accessControl) DelBlackList(ip string) {
	delete(ctrl.blackList, ip)
}

// 正则查找黑名单数据
func (ctrl *accessControl) RegexpBlackList(pattern string) []string {
	result := []string{}
	for key, _ := range ctrl.blackList {
		if match, _ := regexp.MatchString(pattern, key); match {
			result = append(result, key)
		}
	}
	return result
}

func (ctrl *accessControl) ReadIp(req *http.Request) (ip []string, port string) {
	remoteAddr := strings.Split(req.RemoteAddr, ":")

	ip = req.Header["X-Forwarded-For"]
	if len(ip) == 0 {
		ip = []string{remoteAddr[0]}
	}
	return ip, remoteAddr[1]
}

func (ctrl *accessControl) InBlackList(ipGroup []string) (rejectIp string, in bool) {
	_, ok := ctrl.blackList["*"]
	if ok {
		return "*", true
	}
	for _, addr := range ipGroup {
		_, ok := ctrl.blackList[addr]
		if ok {
			return addr, true
		}
	}
	return "", false
}

func (ctrl *accessControl) InWhiteList(ipGroup []string) (ip string, in bool) {
	_, ok := ctrl.whiteList["*"]
	if ok {
		return "*", true
	}

	for _, addr := range ipGroup {
		_, ok := ctrl.whiteList[addr]
		if ok {
			return addr, true
		}
	}
	return "", false
}

func (ctrl *accessControl) AttackCheck(ipGroup []string) int {
	if _, ok := ctrl.InWhiteList(ipGroup); ok {
		return -999
	}
	if _, ok := ctrl.InBlackList(ipGroup); ok {
		return 999
	}

	var cond *AttackCondition
	now := time.Now()

	key := fmt.Sprintf("%v", ipGroup)
	rs := ctrl.attackTimesCacher.Get(key)
	if rs != nil {
		cond = rs.(*AttackCondition)
	}
	if cond != nil {
		// 半秒钟内的请求进行计数
		if now.Sub(cond.LastTime) < ctrl.interval {
			cond.Times++
		}
		cond.LastTime = now
	} else {
		cond = &AttackCondition{1, now}
	}
	ctrl.attackTimesCacher.Put(key, cond, int64(ctrl.cacheTime/1e9))
	return ctrl.limit - cond.Times
}

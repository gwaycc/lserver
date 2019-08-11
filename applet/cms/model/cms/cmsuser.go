package cms

import (
	"encoding/json"
	"strings"

	"github.com/gwaylib/errors"
	"github.com/jameskeane/bcrypt"
)

// 后台用户
type CmsUser struct {
	UserName string
	Passwd   string // 加密的数据
	NickName string // 昵称
	Gid      int
	Status   int // 状态
}

// 密码校验
func (auth *CmsUser) CheckSumPasswd(in string) bool {
	return bcrypt.Match(in, auth.Passwd)
}

// 权限管理
// 以顺序排序的方式进行排列，
// [0][]这一维存一个菜单的数据
// [][0]这一维存储每一个节点的数据
type CmsPriv [][]string

func ParseCmsPriv(in string) (CmsPriv, error) {
	cmsPriv := CmsPriv{}
	if err := json.Unmarshal([]byte(in), &cmsPriv); err != nil {
		return nil, errors.As(err, in)
	}
	return cmsPriv, nil
}

func (cp *CmsPriv) Append(node []string) {
	*cp = append(*cp, node)
}

func (cp CmsPriv) ToJson() string {
	data, err := json.Marshal(&cp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (cp CmsPriv) CheckSum(path string) bool {
	if path == "" || path == "/" {
		return true
	}
	if len(path) > 0 && path[:1] != "/" {
		// 非法的路径
		return false
	}

	node := strings.Split(path, "/")

	newNode := []string{}
	// 移除空节点
	for _, val := range node {
		if len(val) > 0 {
			newNode = append(newNode, val)
		}
	}
	newNodeLen := len(newNode)
	if newNodeLen == 0 {
		// 所有人都可访问根节点
		return true
	}

	// 从请求的末节点向左匹配
	cpLen := len(cp)
	pass := false
	for i := 0; i < cpLen; i++ {
		// 长度小于请求节点长度的，不需要参与校验
		if len(cp[i]) < newNodeLen {
			continue
		}
		pass = true
		// 遍历子节点是否全匹配,如果全匹配，完成校验工作
		for j := newNodeLen - 1; j > -1; j-- {
			if !(cp[i][j] == newNode[j] || cp[i][j] == "*") {
				pass = false
				break
			}
		}
		if pass {
			return true
		}
	}
	return false
}

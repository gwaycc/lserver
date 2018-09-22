package route

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"module/gouuid"
	"module/httptry/ipctrl"
	"module/redis"
	"net/http"
	"time"

	"applet/cms/model/cms"

	"github.com/dchest/captcha"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb"
	"github.com/gwaylib/log"
	"github.com/labstack/echo"
)

const (
	AccSigninViewPath = "/access/signin"
	AccSigninApiPath  = "/access/signin"
	AccPwdApiPath     = "/access/pwd"
	AccSignupApiPath  = "/access/signup"
)

type UserCache struct {
	UserName  string      `json:"username"`
	NickName  string      `json:"nickname"`
	Priv      cms.CmsPriv `json:"priv"`
	Logo      string      `json:"logo"`
	OnlineKey int         `json:"online_key"`
}

func (uc *UserCache) ToJson() []byte {
	data, err := json.Marshal(uc)
	if err != nil {
		panic(err)
	}
	return data
}

const cookieKey = "lserver_cms_session"

func (uc *UserCache) SetToSession(c echo.Context) error {
	sessionName := gouuid.New()
	// 保存到cookies中
	cookie := new(http.Cookie)
	cookie.Name = cookieKey
	cookie.Value = XXTeaEncode(sessionName)
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(1 * 60 * 60 * time.Second)
	// 存储在线授权信息
	if err := redisClient.Set(sessionName, uc, 1*60*60); err != nil {
		return errors.As(err)
	}
	if err := redisClient.Set("lservercms_"+uc.UserName, uc.OnlineKey, 1*60*60); err != nil {
		return errors.As(err)
	}
	c.SetCookie(cookie)
	return nil
}

func (uc *UserCache) CleanSession() error {
	return errors.As(redisClient.Delete("lservercms_" + uc.UserName))
}

func (uc *UserCache) ReAuth(passwd string) bool {
	u, err := cms.NewCmsDB().GetUser(uc.UserName, 1)
	if err != nil {
		log.Warn(errors.As(err))
		return false
	}
	return u.CheckSumPasswd(passwd)
}
func GetUserCache(c echo.Context) *UserCache {
	uCookie, err := c.Cookie(cookieKey)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Warn(errors.As(err))
		}
		return nil
	}

	sessionName := ""
	if err := XXTeaDecode(uCookie.Value, &sessionName); err != nil {
		log.Debug(errors.As(err))
		return nil
	}
	uc := &UserCache{}
	if err := redisClient.Scan(sessionName, uc); err != nil {
		log.Warn(errors.As(err))
		return nil
	}
	onlineKey := 0
	if err := redisClient.Scan("lservercms_"+uc.UserName, &onlineKey); err != nil {
		if !errors.Equal(redis.ErrNil, err) {
			log.Warn(errors.As(err))
		}
		return nil
	}

	// 登录信息已变更，需要重新登录
	if onlineKey != uc.OnlineKey {
		log.Debug(uc.UserName, onlineKey)
		return nil
	}
	// 每次操作延长session时间
	if err := uc.SetToSession(c); err != nil {
		log.Warn(errors.As(err))
		return nil
	}

	return uc
}

func init() {
	r := eweb.Default()

	// view
	r.GET(AccSigninViewPath, AccSigninView)

	// post
	r.POST(AccSigninApiPath, AccSigninApi)
	r.POST(AccPwdApiPath, AccPwdApi)
	r.POST(AccSignupApiPath, AccSignupApi)
}

// 登录页面
func AccSigninView(c echo.Context) error {
	return Index(c)
}

var (
	// 10分钟内连续输错5次
	// TODO:存储到redis中，以便可以在集群中使用
	AttackCtrl = ipctrl.NewAccessController(10*time.Minute, 10*time.Minute, 5)
)

// 登录操作
func AccSigninApi(c echo.Context) error {
	req := c.Request()
	ips, _ := AttackCtrl.ReadIp(req)
	username := FormValue(c, "acc")
	passwd := FormValue(c, "pwd")

	vcodeId := FormValue(c, "vcodeId")
	vcodeData := FormValue(c, "vcodeData")
	if !captcha.VerifyString(vcodeId, vcodeData) {
		return c.String(403, "验证码错误")
	}
	attackKey := []string{fmt.Sprintf("cms_attack_%v_%v", username, ips)}
	if _, in := AttackCtrl.InBlackList(attackKey); in {
		return c.String(403, "因多次错误已禁用该账号，请联系管理员")
	}

	cmsdb := cms.NewCmsDB()
	u, err := cmsdb.GetUser(username, 1)
	if err != nil {
		if errors.ErrNoData.Equal(err) {
			log.Debug(errors.As(err, username))
			// 防暴力猜密码
			if AttackCtrl.AttackCheck(attackKey) < 1 {
				AttackCtrl.AddBlackList(attackKey[0])
			}

			return c.String(403, "账户或密码错误")
		}
		log.Warn(errors.As(err, username))
		return c.String(500, "系统错误")
	}

	if !u.CheckSumPasswd(passwd) {
		// 防暴力猜密码
		if AttackCtrl.AttackCheck(attackKey) < 1 {
			// 记录至数据库
			if err := cms.PutLog(username, "账号攻击", attackKey[0], "受多次输入密码错误，请联系管理员重置密码，并设置高强度密码"); err != nil {
				return c.String(500, "系统错误")
			}
			AttackCtrl.AddBlackList(attackKey[0])
		}

		return c.String(403, "账户或密码错误")
	}

	priv, err := cmsdb.GetPriv(username)
	if err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	respUser := &UserCache{
		UserName:  username,
		NickName:  u.NickName,
		Priv:      priv,
		Logo:      "/img/logo.png",
		OnlineKey: rand.Int(),
	}
	// make session
	if err := respUser.SetToSession(c); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	return c.JSON(200, respUser)
}

// 修改密码
func AccPwdApi(c echo.Context) error {
	oldPwd := FormValue(c, "oldPwd")
	newPwd := FormValue(c, "newPwd")

	uc := GetUserCache(c)
	if uc == nil {
		return c.String(302, "登录已过期, 请重登录")
	}
	if !uc.ReAuth(oldPwd) {
		log.Debug(errors.New("wrong pwd").As(uc.UserName, oldPwd))
		return c.String(403, "原密码错误")
	}

	cmsdb := cms.NewCmsDB()
	if err := cmsdb.ResetPwd(uc.UserName, newPwd); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

// 登出操作
func AccSignupApi(c echo.Context) error {
	uc := GetUserCache(c)
	if uc == nil {
		return nil
	}
	if err := uc.CleanSession(); err != nil {
		log.Warn(errors.As(err))
	}
	return nil
}

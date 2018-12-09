package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/TruthHun/DocHub/helper"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gopkg.in/ldap.v2"
	"time"
)

//用户表
type User struct {
	Id       int    `orm:"column(Id)"`
	Account  string `orm:"size(50);unique;column(Account)"`          //账户
	Email    string `orm:"size(50);unique;column(Email);default();"` //邮箱
	Password string `orm:"size(32);column(Password)"`                //密码
	Username string `orm:"size(16);unique;column(Username)"`         //用户名
	Intro    string `orm:"size(255);default();column(Intro)"`        //个性签名
	Avatar   string `orm:"size(50);default();column(Avatar)"`        //会员头像
}

func NewUser() *User {
	return &User{}
}

func GetTableUser() string {
	return getTable("user")
}

//用户信息表
type UserInfo struct {
	Id         int  `orm:"auto;pk;column(Id)"`                //主键，也就是User表的Id
	Coin       int  `orm:"default(10);index;column(Coin)"`    //金币积分
	Document   int  `orm:"default(0);index;column(Document)"` //文档数量
	Collect    int  `orm:"default(0);column(Collect)"`        //收藏专辑数量，每个收藏专辑下面有文档
	TimeCreate int  `orm:"column(TimeCreate);default(0)"`     //用户注册时间
	Status     bool `orm:"column(Status);default(true)"`      //用户信息状态，false(即0)表示被禁用
	Admin      bool `orm:"column(Admin);default(false)"`      //IsAdmin
}

func NewUserInfo() *UserInfo {
	return &UserInfo{}
}
func GetTableUserInfo() string {
	return getTable("user_info")
}

//根据条件查询用户信息，比如用户登录、用户列表等的获取也可以使用这个函数
//@param            p           int         页码
//@param            listRows    int         每页显示记录数
//@param            orderby     string      排序，如"id desc"
//@param            fields      []string    需要查询的字段
//@param            cond        string      查询条件
//@param            args        ...interface{}  查询条件参数
func (u *User) UserList(p, listRows int, orderby, fields, cond string, args ...interface{}) (params []orm.Params, totalRows int, err error) {
	o := orm.NewOrm()
	if len(orderby) == 0 || orderby == "" {
		orderby = "i.Id desc"
	}
	if len(fields) == 0 {
		fields = "*"
	}
	if len(cond) > 0 {
		cond = "where " + cond
	}

	sqlCount := fmt.Sprintf("select count(i.Id) cnt from %v u left join %v i on u.Id=i.Id %v limit 1",
		GetTableUser(), GetTableUserInfo(), cond,
	)
	var one []orm.Params
	if rows, err := o.Raw(sqlCount, args...).Values(&one); err == nil && rows > 0 {
		totalRows = helper.Interface2Int(one[0]["cnt"])
	}

	sql := fmt.Sprintf("select %v from %v u left join %v i on u.Id=i.Id %v order by %v limit %v offset %v",
		fields, GetTableUser(), GetTableUserInfo(), cond, orderby, listRows, (p-1)*listRows,
	)
	_, err = o.Raw(sql, args...).Values(&params)
	return params, totalRows, err
}

//获取User表的字段
func (u *User) Fields() map[string]string {
	var fields map[string]string
	fields = make(map[string]string)
	v := reflect.ValueOf(u).Elem()
	k := v.Type()
	num := v.NumField()
	for i := 0; i < num; i++ {
		key := k.Field(i)
		fields[key.Name] = key.Name
	}
	return fields
}

//根据用户id获取用户info表的信息
//@param            uid         interface{}         用户UID
//@return           UserInfo                        用户信息
func (u *User) UserInfo(uid interface{}) UserInfo {
	var info UserInfo
	orm.NewOrm().QueryTable(GetTableUserInfo()).Filter("id", uid).One(&info)
	return info
}

//根据查询条件查询User表
//@param            cond            *orm.Condition          查询条件
//@return                           User                    返回查询到的User数据
func (u *User) GetUserField(cond *orm.Condition) (user User) {
	orm.NewOrm().QueryTable(GetTableUser()).SetCond(cond).One(&user)
	return user
}

//用户注册
//@param            account           string        账户
//@param            email           string          邮箱
//@param            username        string          用户名
//@param            password        string          密码
//@param            repassword      string          确认密码
//@param            intro           string          签名
//@return                           error           错误
//@return                           int             注册成功时返回注册id
func (u *User) Reg(account, email, username, password, repassword, intro string) (error, int) {
	var (
		user User
		o    = orm.NewOrm()
		now  = time.Now().Unix()
	)

	l := strings.Count(username, "") - 1
	if l < 2 || l > 16 {
		return errors.New("用户名长度限制在2-16个字符"), 0
	}
	if o.QueryTable(GetTableUser()).Filter("Username", username).One(&user); user.Id > 0 {
		return errors.New("用户名已被注册，请更换新的用户名"), 0
	}
	pwd := helper.MyMD5(password)
	if pwd != helper.MyMD5(repassword) {
		return errors.New("密码和确认密码不一致"), 0
	}
	user = User{Account: account, Email: email, Username: username, Password: pwd, Intro: intro}
	_, err := o.Insert(&user)
	if user.Id > 0 {
		//coin := beego.AppConfig.DefaultInt("coinreg", 10)
		coin := NewSys().GetByField("CoinReg").CoinReg
		var userinfo = UserInfo{Id: user.Id, Status: true, Coin: coin, TimeCreate: int(now)}
		_, err = o.Insert(&userinfo)
	}
	return err, user.Id
}

//获取除了用户密码之外的用户全部信息
//@param                id              用户id
//@return               params          用户信息
//@return               rows            记录数
//@return               err             错误
func (this *User) GetById(id interface{}) (params orm.Params, rows int64, err error) {
	var data []orm.Params
	tables := []string{GetTableUser() + " u", GetTableUserInfo() + " ui"}
	on := []map[string]string{
		{"u.Id": "ui.Id"},
	}
	fields := map[string][]string{
		"u":  GetFields(NewUser()),
		"ui": GetFields(NewUserInfo()),
	}
	if sql, err := LeftJoinSqlBuild(tables, on, fields, 1, 1, nil, nil, "u.Id=?"); err == nil {
		if rows, err = orm.NewOrm().Raw(sql, id).Values(&data); err == nil && len(data) > 0 {
			params = data[0]
		}
	}
	return
}

//
// 使用LDAP进行登录
//
func (u *User) LDAPLogin(account, password string) (params []orm.Params, totalRows int, err error) {

	if beego.AppConfig.DefaultBool("ldap::enable", false) == false {
		return params, totalRows, errors.New("没有启用LDAP服务器")
	}
	addr := fmt.Sprintf("%s:%d", beego.AppConfig.String("ldap::host"), beego.AppConfig.DefaultInt("ldap::port", 3268))
	lc, err := ldap.Dial("tcp", addr)
	if err != nil {
		return params, totalRows, errors.New("无法连接到LDAP服务器")
	}
	defer lc.Close()
	err = lc.Bind(beego.AppConfig.String("ldap::user"), beego.AppConfig.String("ldap::password"))
	if err != nil {
		return params, totalRows, errors.New("第一次LDAP绑定失败")
	}
	searchRequest := ldap.NewSearchRequest(
		beego.AppConfig.String("ldap::base"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//修改objectClass通过配置文件获取值
		fmt.Sprintf("(&(%s)(%s=%s))", beego.AppConfig.String("ldap::filter"), beego.AppConfig.String("ldap::attribute"), account),
		[]string{"dn", "mail", "displayName"},
		nil,
	)
	searchResult, err := lc.Search(searchRequest)
	if err != nil {
		return params, totalRows, errors.New("LDAP搜索失败")
	}
	if len(searchResult.Entries) != 1 {
		return params, totalRows, errors.New("LDAP用户不存在或者多于一个")
	}
	userdn := searchResult.Entries[0].DN
	err = lc.Bind(userdn, password)
	if err != nil {
		return params, totalRows, errors.New("用户密码错误")
	}
	if u.Id == 0 {
		ldapMail := searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap::mail"))
		ldapName := searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap::name"))
		err, u.Id = u.Reg(account, ldapMail, ldapName, password, password, ldapName)
		if err != nil {
			return params, totalRows, errors.New(fmt.Sprint("自动注册LDAP用户错误:%s", err.Error()))
		}
	}
	return u.UserList(1, 1, "", "", "u.`account`=? ", account)
}

func (u *User) CheckLogin(account, password string) (params []orm.Params, totalRows int, err error) {
	usr := u.GetUserField(orm.NewCondition().And("Email", account).Or("Account", account))
	//1.Check Local Login
	if usr.Id > 0 && strings.HasSuffix(strings.ToLower(usr.Account), "@local") {
		return u.UserList(1, 1, "", "", "u.`account`=? and u.`password`=?", usr.Account, helper.MyMD5(password))
	}
	return usr.LDAPLogin(account, password)
}

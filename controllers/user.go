package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"myERP/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	beego.Controller
}

// 随机生成6位数代码
func captcha(n int) string {
	code := ""
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		code = fmt.Sprintf("%s%d", code, rand.Intn(10))
	}
	return code
}

// 密码加密
func hashAndSalt(pwd string) string {
	tmp := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(tmp, bcrypt.MinCost)
	if err != nil {
		log.Println("密码加密失败, err:", err)
	}
	return string(hash)
}

// 密码解密
func comparePassword(hash string, plainPwd []byte) bool {
	bytehash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(bytehash, plainPwd)
	if err != nil {
		log.Println("密码解密失败， err:", err)
		return false
	}
	return true
}

// RegisterCompany 注册公司
func (c *UserController) RegisterCompany() {
	c.TplName = "registerCompany.html"
}

// HandleRegisterCompany 注册公司处理
func (c *UserController) HandleRegisterCompany() {
	// 获取数据
	companyName := c.GetString("companyname")
	// 校验数据
	if companyName == "" {
		log.Println("公司名称不能为空")
		c.Data["errmsg"] = "公司名称不能为空"
		c.TplName = "registerCompany.html"
		return
	}
	// 处理数据
	o := orm.NewOrm()
	var company models.Company
	company.CompanyName = companyName
	err := o.Read(&company, "companyName")
	if err == nil {
		log.Println("公司名称已被注册，请直接登陆")
		c.Data["errmsg"] = "公司名称已被注册，请直接登陆"
		c.Redirect("/register", 302)
		return
	}
	companyCode := captcha(6)
	company.CompanyCode = companyCode
	_, err = o.Insert(&company)
	if err != nil {
		log.Println("注册公司失败，err:", err)
		c.Data["errmsg"] = "注册公司失败"
		c.TplName = "registerCompany.html"
		return
	}
	// 返回数据
	c.Redirect("/register", 302)
}

// RegisterAdmin 注册
func (c *UserController) RegisterAdmin() {
	c.TplName = "registerAdmin.html"
}

// HandleRegisterAdmin 注册业务处理
func (c *UserController) HandleRegisterAdmin() {
	// 获取数据
	companyCode := c.GetString("companycode")
	username := c.GetString("username")
	password := c.GetString("password")
	repassword := c.GetString("repassword")
	// 校验数据
	if companyCode == "" || username == "" || password == "" || repassword == "" {
		log.Println("公司代码或用户名或密码或确认密码不能为空")
		c.Data["errmsg"] = "公司代码或用户名或密码或确认密码不能为空"
		c.TplName = "registerAdmin.html"
		return
	}
	if repassword != password {
		log.Println("密码不一致")
		c.Data["errmsg"] = "密码不一致"
		c.TplName = "registerAdmin.html"
		return
	}
	hashPassword := hashAndSalt(password)
	fmt.Println(hashPassword)
	// 处理数据
	o := orm.NewOrm()
	var company models.Company
	company.CompanyCode = companyCode
	err := o.Read(&company, "companyCode")
	if err != nil {
		log.Println("公司代码错误，err:", err)
		c.Data["errmsg"] = "公司代码错误"
		c.TplName = "registerAdmin.html"
		return
	}
	var admin models.Admin
	admin.Company = &company
	admin.Name = username
	admin.Password = hashPassword
	_, err = o.Insert(&admin)
	if err != nil {
		log.Println("加入用户失败，err:", err)
		c.Data["errmsg"] = "加入用户失败"
		c.TplName = "registerAdmin.html"
		return
	}

	// 注册成功，设置cookie
	c.Ctx.SetCookie("companyCode", company.CompanyCode, 60*10)
	c.Ctx.SetCookie("username", admin.Name, 60*10)

	// 返回登陆页面
	c.Redirect("/login", 302)
}

// Login 登陆
func (c *UserController) Login() {
	companyCode := c.Ctx.GetCookie("companycode")
	username := c.Ctx.GetCookie("username")
	if companyCode == "" && username == "" {
		c.Data["checked"] = "checked"
	} else {
		c.Data["checked"] = ""
	}
	c.Data["companycode"] = companyCode
	c.Data["username"] = username
	c.TplName = "login.html"
}

// HandleLogin 登陆业务处理
func (c *UserController) HandleLogin() {
	// 获取数据
	companyCode := c.GetString("companycode")
	username := c.GetString("username")
	password := c.GetString("password")
	m1, _ := c.GetInt("m1")
	fmt.Println(companyCode, username, password, m1)
	// 校验数据
	if companyCode == "" || username == "" || password == "" {
		log.Println("公司代码或用户名或密码不能为空")
		c.Data["errmsg"] = "公司代码或用户名或密码不能为空"
		c.TplName = "login.html"
		return
	}
	// 处理数据，查询登陆账号是否存在
	o := orm.NewOrm()
	var company models.Company
	company.CompanyCode = companyCode
	err := o.Read(&company, "companyCode")
	if err != nil {
		log.Println("公司代码错误， err:", err)
		c.Data["errmsg"] = "公司代码错误"
		c.TplName = "login.html"
		return
	}

	var admin models.Admin
	qs := o.QueryTable("Admin")
	qs.RelatedSel("Company").Filter("Company__CompanyCode", companyCode).Filter("Name", username).One(&admin)
	if admin.Name == "" {
		log.Println("用户名错误，err:", err)
		c.Data["errmsg"] = "用户名错误"
		c.TplName = "login.html"
		return
	}
	// 校验密码
	pwdmatch := comparePassword(admin.Password, []byte(password))
	if pwdmatch == false {
		log.Println("密码输入错误！")
		c.Data["errmsg"] = "密码输入错误"
		c.TplName = "login.html"
		return
	}

	// 根据m1值，判断是否实现记住公司代码与用户名
	if m1 == 2 {
		c.Ctx.SetCookie("companycode", companyCode, 60*100)
		c.Ctx.SetCookie("username", username, 60*100)
		fmt.Println("setcookie成功")
	} else {
		c.Ctx.SetCookie("companycode", companyCode, -1)
		c.Ctx.SetCookie("username", username, -1)
		fmt.Println("deletecookie成功")
	}
	fmt.Println(admin.Name)
	// 设置 session，用户登陆后页面使用
	c.SetSession("companycode", company.CompanyCode)
	c.SetSession("username", admin.Name)
	fmt.Println("设置session成功")
	c.Redirect("/erpIndex", 302)
}

// ERP首页
func (c *UserController) ErpIndex() {
	companycode := c.GetSession("companycode")
	o := orm.NewOrm()
	var company models.Company
	company.CompanyCode = companycode.(string)
	err := o.Read(&company, "companycode")
	if err != nil {
		log.Println("公司代码错误，无法查找公司， err:", err)
		c.Data["errmsg"] = "公司代码错误，无法查找公司"
		return
	}
	if companycode != nil {
		c.Data["companyName"] = company.CompanyName
		c.Data["companyCode"] = companycode.(string)
	} else {
		c.Data["companyName"] = company.CompanyName
		c.Data["companyCode"] = ""
	}
	c.Layout = "erpIndex.html"
	c.TplName = "companyInfo.html"
}

// 我的企业 -> 企业信息
func (c *UserController) CompanyInfo() {
	company := GetCompany(&c.Controller)
	c.Data["companyName"] = company.CompanyName
	c.Data["companyCode"] = company.CompanyCode
	c.Data["num"] = 1
	c.Layout = "erpIndex.html"
	c.TplName = "companyInfo.html"
}

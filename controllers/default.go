package controllers

import (
	"myERP/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

// GetUser 获取当前登录用户
func GetUser(c *beego.Controller) models.Admin {
	// 根据 session 获取当前登录用户名
	username := c.GetSession("username")
	o := orm.NewOrm()
	var admin models.Admin
	admin.Name = username.(string)
	o.Read(&admin, "name")
	return admin
}

// GetCompany 获取当前企业
func GetCompany(c *beego.Controller) models.Company {
	// 根据 session 获取当前登录所属的企业
	companycode := c.GetSession("companycode")
	o := orm.NewOrm()
	var company models.Company
	company.CompanyCode = companycode.(string)
	o.Read(&company, "companyCode")
	return company
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

// 信息组合
func (c *MainController) InfoGroup() {
	company := GetCompany(&c.Controller)
	admin := GetUser(&c.Controller)
	c.Data["companyName"] = company.CompanyName
	c.Data["adminName"] = admin.Name
	c.Layout = "infoGroup.html"
	c.TplName = "goods.html"
}

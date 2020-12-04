package routers

import (
	"myERP/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	// 主注册页面，注册管理员
	beego.Router("/register", &controllers.UserController{}, "get:RegisterAdmin;post:HandleRegisterAdmin")
	// 注册公司
	beego.Router("/registerCompany", &controllers.UserController{}, "get:RegisterCompany;post:HandleRegisterCompany")
	// 登陆页面
	beego.Router("/login", &controllers.UserController{}, "get:Login;post:HandleLogin")

	// ERP首页
	beego.Router("/erpIndex", &controllers.UserController{}, "get:ErpIndex")

	// 我的企业
	// 企业信息
	beego.Router("/companyInfo", &controllers.UserController{}, "get:CompanyInfo")

	// 信息组合
	// 品目信息
	beego.Router("/goodsInfo", &controllers.GoodsController{}, "get:GoodsInfo")
	// 新增品目
	beego.Router("/addGoods", &controllers.GoodsController{}, "get:AddGoods;post:HandleAddGoods")
	// 用户中心
}

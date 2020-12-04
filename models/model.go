package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// Company 公司结构体
type Company struct {
	Id          int
	CompanyName string      `orm:"size(40);unique"`
	CompanyCode string      `orm:"size(6);unique"`
	Member      []*Admin    `orm:"reverse(many)"`
	GoodsSKU    []*GoodsSKU `orm:"reverse(many)"`
}

// Admin 管理员结构体
type Admin struct {
	Id       int
	Company  *Company `orm:"rel(fk)"`
	Name     string   `orm:"size(40);unique"`
	Password string   `orm:"size(200)"`
}

type GoodsSKU struct {
	Id        int
	Code      string     `orm:"unique"`
	Name      string     `orm:"size(20)"`
	Spec      string     // 规格
	CTN       string     // 包装
	GoodsKind *GoodsKind `orm:"rel(fk)"` // 分类
	Img       *Img       `orm:"rel(fk)"` // 图片
	Company   *Company   `orm:"rel(fk)"`
}

type GoodsKind struct {
	Id        int
	Name      string      `orm:"size(10);unique"`
	GoodsSKUs []*GoodsSKU `orm:"reverse(many)"`
}

type Img struct {
	Id       int
	Name     string
	Path     string
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
}

// init 初始化连接数据库
func init() {
	// 注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/myerp")
	// 注册表结构
	orm.RegisterModel(new(Company), new(Admin), new(GoodsSKU), new(GoodsKind), new(Img))
	// 运行生成表
	orm.RunSyncdb("default", false, true)
}

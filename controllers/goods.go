package controllers

import (
	"fmt"
	"log"
	"myERP/models"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type GoodsController struct {
	beego.Controller
}

// GoodsInfo 品目信息
func (c *GoodsController) GoodsInfo() {
	company := GetCompany(&c.Controller)
	fmt.Println(company.CompanyCode)

	o := orm.NewOrm()
	var goodsSKUs []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU")
	qs.RelatedSel("Company", "GoodsKind", "Img").Filter("Company__CompanyCode", company.CompanyCode).All(&goodsSKUs)
	fmt.Println(goodsSKUs)
	c.Data["goodsSKUs"] = goodsSKUs
	c.Layout = "infoGroup.html"
	c.TplName = "goods.html"
}

// 新增品目
func (c *GoodsController) AddGoods() {
	c.Layout = "infoGroup.html"
	c.TplName = "addGoods.html"
}

// 新增品目业务处理
func (c *GoodsController) HandleAddGoods() {
	company := GetCompany(&c.Controller)
	productcode := c.GetString("productCode")
	productname := c.GetString("productName")
	spec := c.GetString("spec")
	ctn := c.GetString("package")
	kind := c.GetString("kind")
	file, fileheader, err := c.GetFile("img") // 返回文件、文件信息头、错误信息
	if err != nil {
		log.Println("上传图片失败，err:", err)
		c.Layout = "infoGroup.html"
		c.TplName = "addGoods.html"
		return
	}

	defer file.Close() // 关闭上传的文件，否则出现临时文件不清楚的情况

	photoName := fileheader.Filename
	fmt.Println(photoName)
	photo := strings.Split(photoName, ".")
	fmt.Println(photo)
	layout := strings.ToLower(photo[len(photo)-1])

	if layout != "jpg" && layout != "png" && layout != "git" {
		log.Println("请上传符合格式的图片")
		return
	}

	o := orm.NewOrm()
	var goodskind models.GoodsKind
	goodskind.Name = kind
	err = o.Read(&goodskind, "name")
	if err != nil {
		log.Println("品目分类中没有此分类，即将为你新增该品类的分类， err:", err)
		_, err = o.Insert(&goodskind)
		if err != nil {
			log.Println("新增品目分类失败， err:", err)
			return
		}
		fmt.Println("分类增加成功")
	}

	dirstr := fmt.Sprintf("static/upload/%s", kind)

	str := fmt.Sprintf("static/upload/%s/%s", kind, photoName)
	fmt.Printf(str)

	// 创建上图上传的文件夹
	os.Mkdir(dirstr, 0777)

	err = c.SaveToFile("img", str) // 要自己先创建文件夹，它打开的是文件
	if err != nil {
		log.Println("图片上传失败，err:", err)
		c.TplName = "addGoods.html"
		return
	} else {
		fmt.Println("图片上传成功")
	}

	var goodsSKU models.GoodsSKU
	goodsSKU.Code = productcode
	goodsSKU.Name = productname
	goodsSKU.Spec = spec
	goodsSKU.CTN = ctn
	goodsSKU.GoodsKind = &goodskind
	// goodsSKU.Img = photoName
	goodsSKU.Company = &company
	_, err = o.Insert(&goodsSKU)
	if err != nil {
		log.Println("增加品目信息失败，err:", err)
		c.Layout = "infoGroup.html"
		c.TplName = "addGoods.html"
		return
	}

	var img models.Img
	img.Name = photoName
	img.Path = str
	img.GoodsSKU = &goodsSKU

	c.Redirect("/goodsInfo", 302)
}

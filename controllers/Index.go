package controllers

import (
	"github.com/astaxie/beego"
	"go-shorturl/models"
	"github.com/astaxie/beego/orm"
	"github.com/skip2/go-qrcode"
	"github.com/astaxie/beego/cache"
	"log"
)

var (
	urlcache cache.Cache
)

const TIMEOUT = 300

func init() {
	urlcache, _ = cache.NewCache("memory", `{"interval":300}`)
}

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Index() {
	url := this.GetString("url")
	if url != "" {
		urlOne, err := models.UrlGetByUrl(url)
		if err == nil {
			urlcache.Put(urlOne.Short, url, TIMEOUT)
			this.Data["shortUrl"] = beego.AppConfig.String("website") + urlOne.Short
			this.Data["longUrl"] = url
		} else {
			urlModel := new(models.Url)
			urlModel.Url = url
			id, _ := models.UrlAdd(urlModel)
			urlModel.Short = models.Generate(id)
			orm.NewOrm().Update(urlModel)
			urlcache.Put(urlModel.Short, url, TIMEOUT)
			this.Data["shortUrl"] = beego.AppConfig.String("website") + urlModel.Short
			this.Data["longUrl"] = url
		}
	}
	this.TplName = "index/index.tpl";
}
func (this *IndexController) Qrcodeimg() {
	url := this.GetString("url")
	if url != "" {
		this.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
		png, _ := qrcode.Encode(url, qrcode.Medium, 256)
		this.Ctx.ResponseWriter.Write(png)
		this.StopRun()
	}
}

func (this *IndexController) Jump() {
	url := this.Ctx.Input.Param(":url")

	if url == "favicon.ico" {
		return
	}

	//  从 cache 获取不到, 从 db 读取
	//  这里可以优化 cache 采用 lruchache 算法
	redirectUrl := ""
	if longurl := urlcache.Get(url); longurl == nil {
		o := orm.NewOrm()
		urlModel := models.Url{Short: url}
		err := o.Read(&urlModel, "short")

		if err == orm.ErrNoRows {
			log.Println("查询不到")
			this.Ctx.Output.SetStatus(404)
			this.StopRun()
			return
		}
		redirectUrl = urlModel.Url
	} else {
		redirectUrl = longurl.(string)
	}

	log.Println(redirectUrl, "Longurl");
	if redirectUrl != "" {
		//model := new(models.Detail)
		//model.AccessIp = this.Ctx.Input.IP()
		////model.Browser = this.Ctx.Request.UserAgent()
		//model.Source = this.Ctx.Request.Referer()
		//model.AccessTime = time.Now()
		//model.Short = url
		//orm.NewOrm().Insert(model)
		urlcache.Put(url, redirectUrl, TIMEOUT)
		this.Redirect(redirectUrl, 302)
	}
	this.StopRun()
}

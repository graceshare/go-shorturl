package main

import (
	_ "go-shorturl/routers"
	"github.com/astaxie/beego"
	"go-shorturl/models"
)

func main() {
	models.Init()
	beego.Run()
}


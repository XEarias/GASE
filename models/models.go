package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlConnData struct {
	user   string
	pass   string
	ip     string
	dbName string
}

func init() {

	RunMode := beego.BConfig.RunMode

	var mysqlConnData mysqlConnData

	mysqlConnData.user = beego.AppConfig.String(RunMode + "::mysqluser")
	mysqlConnData.pass = beego.AppConfig.String(RunMode + "::mysqlpass")
	mysqlConnData.dbName = beego.AppConfig.String(RunMode + "::mysqldb")

	//fmt.Println(mysqlConnData)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnData.user+":"+mysqlConnData.pass+"@/"+mysqlConnData.dbName+"?charset=utf8")

	orm.RegisterModel(new(Activities), new(Clients), new(Countries), new(Coupons), new(Currencies), new(Gateways), new(Images), new(Locations), new(Orders), new(Portfolios), new(Prices), new(Sectors), new(Services))

	/* 	// Create database from models.
	   	name := "default"

	   	// Drop table and re-create.
	   	force := true

	   	// Print log.
	   	verbose := true

	   	// Error.
	   	err := orm.RunSyncdb(name, force, verbose)
	   	if err != nil {
	   		fmt.Println(err)
	   	} */

	err := AddDefaultDataCurrencies()
	if err != nil {

	}

	err = AddDefaultDataSectors()
	if err != nil {

	}

	errors := addDefaultDataActivities()

	if len(errors) > 0 {
		/* for _, err := range errors {
			println(err.Error())
		} */
	}

	err = AddDefaultDataGateways()
	if err != nil {

	}

	err = AddDefaultDataCurrencies()
	if err != nil {

	}

	errors = addRelationsGatewaysCurrencies()

	if len(errors) > 0 {
		/* for _, err := range errors {
			println(err.Error())
		} */
	}
}

package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

type AccessServer struct {
	gorm.Model
	Name                string
	Endpoint            string
	User                string
	Password            string
	Sni                 string
	SkipVerifySSL       bool
	UseRemoteRouteRules bool
	LocalRouteRules     string
	ProxyDomains        string
}

type Logs struct {
	gorm.Model
	Text string
}

func InitDB(path string) error {

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&AccessServer{})
	Db = db
	return nil
}

func AddAccessServer(server AccessServer) error {
	ret := Db.Create(&server)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func GetAllAccessServer() ([]AccessServer, error) {
	var records []AccessServer
	ret := Db.Find(&records)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return records, nil
}

func DeleteAccessServer(ID uint) error {
	var server AccessServer
	ret := Db.Delete(&server, ID)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func UpdateAccessServer(server AccessServer) error {

	ret := Db.Model(&server).Updates(&map[string]interface{}{
		"ID":                  server.ID,
		"Name":                server.Name,
		"Endpoint":            server.Endpoint,
		"User":                server.User,
		"Password":            server.Password,
		"Sni":                 server.Sni,
		"UseRemoteRouteRules": server.UseRemoteRouteRules,
		"SkipVerifySSL":       server.SkipVerifySSL,
		"LocalRouteRules":     server.LocalRouteRules,
		"ProxyDomains":        server.ProxyDomains,
	})

	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

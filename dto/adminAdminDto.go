package dto

import "go-iptv/models"

type AdminsDto struct {
	LoginUser string             `json:"loginuser"`
	Title     string             `json:"title"`
	Admins    []models.IptvAdmin `json:"admins"`
}

type UdataDto struct {
	LoginUser string `json:"loginuser"`
	Title     string `json:"title"`
	Version   string `json:"version"`
}

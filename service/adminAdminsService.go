package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/url"
	"runtime"
)

func Admins(params url.Values) dto.ReturnJsonDto {
	username := params.Get("username")
	oldPassword := params.Get("oldpassword")
	newpassword := params.Get("newpassword")
	newpassword2 := params.Get("newpassword_2")

	if username == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "用户名不能为空", Type: "danger"}
	}
	if oldPassword == "" && (newpassword != "" || newpassword2 != "") {
		return dto.ReturnJsonDto{Code: 0, Msg: "旧密码不能为空"}
	}

	if newpassword != newpassword2 && newpassword != "" && newpassword2 != "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "两次新密码不一致", Type: "danger"}
	}

	if oldPassword == newpassword && newpassword != "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "新密码不能与旧密码相同", Type: "danger"}
	}

	if !until.IsSafe(username) {
		return dto.ReturnJsonDto{Code: 0, Msg: "用户名不合法", Type: "danger"}
	}

	var adminData models.IptvAdmin
	dao.DB.Model(&models.IptvAdmin{}).Where("id = ?", 1).First(&adminData)
	if adminData.PassWord != until.HashPassword(oldPassword) {
		return dto.ReturnJsonDto{Code: 0, Msg: "旧密码错误", Type: "danger"}
	}

	dao.DB.Model(&models.IptvAdmin{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"password": until.HashPassword(newpassword),
		"username": username,
	})

	// TODO
	return dto.ReturnJsonDto{Code: 1, Msg: "修改成功", Type: "success"}

	// TODO
}

func UpdataCheck() dto.ReturnJsonDto {
	a, v, err := until.CheckNewVer(until.GetVersion())
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "检查更新失败: " + err.Error(), Type: "danger"}
	}
	if a {
		return dto.ReturnJsonDto{Code: 1, Msg: v, Type: "success"}
	}
	return dto.ReturnJsonDto{Code: 2, Msg: "当前已是最新版本", Type: "success"}
}

func UpdataDown() dto.ReturnJsonDto {
	a, v, err := until.DownloadAndVerify(runtime.GOARCH)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "下载失败: " + err.Error(), Type: "danger"}
	}
	if a {
		return dto.ReturnJsonDto{Code: 1, Msg: v, Type: "success"}
	}
	return dto.ReturnJsonDto{Code: 0, Msg: "下载失败: " + v, Type: "danger"}
}

func Updata() dto.ReturnJsonDto {
	err := until.UpdateSignal()
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "触发更新失败: " + err.Error(), Type: "danger"}
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "已触发更新，请稍后...", Type: "success"}
}

package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
	"github.com/astaxie/beego"
)

type File struct {
	Id					int
	File_address		string
	File_date			string
	Isfinance			string
	Describe			string
	Jmuser_id			string	//
	Jm_name				string	//用户名
}

func GetFileById(id int) []File{

	o := orm.NewOrm()
	o.Using("default")
	var file []File
	o.Raw("select *from `file` where id=? ",id).QueryRows(&file)
	for i := 0 ; i < len(file); i++{
		file[i].File_address = beego.AppConfig.String("domain") + file[i].File_address
	}
	return file
}

//获得所有财务账单数据
func GetFile(uid string,pageNum,num int) interface{} {
	var getinfo GetInfo
	o := orm.NewOrm()
	o.Using("default")
	var file []File
	var maps []orm.Params
	var role_id	string
	num2, _ := o.Raw("SELECT role_id from `permission` where jmuser_id=?",uid).Values(&maps)
	if num2 > 0{
		role_id = maps[0]["role_id"].(string)
	}
	/*超级管理员或是财务*/
	if role_id == "1" || role_id == "4" {
		num2, _ = o.Raw("select count(*) as num from `file` where isfinance=1").Values(&maps)
		//页码,页记录数,总记录数
		getinfo.Pager.SumPage = maps[0]["num"].(string)
		if num2 > 0 {
			num3 ,_ := o.Raw("select *from `file` where isfinance=1 order by id desc limit ?,?", (pageNum-1)*num, num).QueryRows(&file)
			if num3 > 0{
				for k,v := range file{
					num3,_:= o.Raw("select jm_name from `jmuser` where id=?",v.Jmuser_id).Values(&maps)
					if num3 > 0 {
						file[k].Jm_name = maps[0]["jm_name"].(string)
					}
				}
			}
			//data
			getinfo.Data = file
			getinfo.Pager.ClientPage = pageNum
			getinfo.Pager.EveryPage = num
		}
	}else {
		num2, _ = o.Raw("select count(*) as num from `file` where isfinance=1 and jmuser_id=?",uid).Values(&maps)
		//页码,页记录数,总记录数
		getinfo.Pager.SumPage = maps[0]["num"].(string)
		if num2 > 0 {
			num3 ,_ := o.Raw("select *from `file` where isfinance=1 and jmuser_id=? order by id desc limit ?,?",uid, (pageNum-1)*num, num).QueryRows(&file)
			if num3 > 0{
				for k,v := range file{
					num3,_:= o.Raw("select jm_name from `jmuser` where id=?",v.Jmuser_id).Values(&maps)
					if num3 > 0 {
						file[k].Jm_name = maps[0]["jm_name"].(string)
					}
				}
			}
			//data
			getinfo.Data = file
			getinfo.Pager.ClientPage = pageNum
			getinfo.Pager.EveryPage = num
		}
	}

	return getinfo
}
//上传文件
func UploadFile(uid,file_address,file_date,describe string,isfinance int) int{
	//判空前端做
	o := orm.NewOrm()
	o.Using("default")
	_,err := o.Raw("insert `file`(jmuser_id,file_address,file_date,`describe`,isfinance) value(?,?,?,?,?)",uid,file_address,file_date,describe,isfinance).Exec()
	if err == nil {
		return 1
	}
	fmt.Println("文件上传失败")
	return 0
}
//获得刚上传的文件id
func GetFilePath(file_address string) string{
	id := "0"
	//判空前端做
	o := orm.NewOrm()
	o.Using("default")
	_,err := o.Raw("insert `file`(file_address) value(?)",file_address).Exec()
	if err == nil {

		var maps []orm.Params
		num, err := o.Raw("SELECT max(id) as id from `file`").Values(&maps)
		if err == nil && num > 0 {
			id = maps[0]["id"].(string)
		}
	}
	fmt.Println("文件上传失败")
	return id
}

//删除文件信息
func DeleteFileById(id string) int64{
	var num int64;//返回影响的行数

	o := orm.NewOrm()
	o.Using("default")
	res,err := o.Raw("delete from `file` where id=?",id).Exec()
	if err == nil {
		//同时删除上传的文件
		//先删除文件
		var maps []orm.Params
		num2, err := o.Raw("SELECT file_address FROM `file` WHERE id = ?", id).Values(&maps)
		if err == nil && num2 > 0 {
			os.Remove("."+maps[0]["file_address"].(string)) // slene
		}
		//影响的行数
		num, _ = res.RowsAffected()
	}
	fmt.Println("文件信息删除失败")
	return num
}

//更新文件,注：提交相同数据返回行数为0
func UpdateFile(id,file_address,file_date,describe string) int64{
	var num int64;//返回影响的行数
	o := orm.NewOrm()
	o.Using("default")

	//先删除文件
	var maps []orm.Params
	num, err := o.Raw("SELECT file_address FROM `file` WHERE id = ?", id).Values(&maps)
	if err == nil && num > 0 {
		os.Remove("."+maps[0]["file_address"].(string)) // slene
	}
	//更新录入新的文件路径
	res,err := o.Raw("update `file` set file_address=?,`describe`=?,file_date=? where id=?",file_address,describe,file_date,id).Exec()
	if err == nil {
		//同时理论上应该删除之前上传的图片
		//......
		num, _ = res.RowsAffected()
	}
	fmt.Println("文件更新失败")
	return num
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(File))
}

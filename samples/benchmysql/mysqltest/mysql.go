package mysqltest 

import (
	"fmt"
	"time"
	"database/sql"
	"github.com/joyous-x/saturn/dbs"
	"github.com/joyous-x/saturn/common/xlog"
)

func GetCacheDB() (*sql.DB, error) {
	tmpConf := dbs.MysqlConf{
		Type: "mysql",
		Host: "154.8.195.109:33",
		User: "test",
		Passwd: "test",
		DbName: "cache_db",
		Debug: false,
	}
	if nil == dbs.MysqlInst(tmpConf) {
		xlog.Error("invalid mysql instance")
	}

	db, err := dbs.MysqlInst().DB("")
	if  err != nil {
		xlog.Error("mysql db:%s error: %v", "", err)
	} else if err := db.Ping(); err != nil {
		xlog.Error("mysql db:%s error: %v", "", err)
	} else {
		xlog.Debug("mysql db:%s ping ok", "")
	}
	return db, err
}

func Insert2Cache(id int, db *sql.DB) {
	rst, err := db.Exec("Insert into t_cache_1(`HIDMD5`,`BrandName`,`Class`,`ExeFilename`,`HardwareID`,`ResultXML`, `LastUpdateTime`) value(?,?,?,?,?,?,?)",
						fmt.Sprintf("HIDMD5-%v", id),fmt.Sprintf("brand-%v", id),fmt.Sprintf("Class-%v", id),fmt.Sprintf("ExeFilename-%v", id),fmt.Sprintf("HardwareID-%v", id),
						`<item vendor="Intel Corporation" pubtime="2018-11-28" file="BA58F2E449E64D18A22EDBF8981EC251.zip" md5="3DD82B97912D26C71C7883CEE6E0E6C4" driverid="0" version="2.2.100.48032" provider="Intel Corporation" class="system" selected="0" id="1933246"  exefile="huawei_Hubble_SGX_win10x64_2_2_100_48032.exe" drivername="HUAWEI电脑荣耀 MagicBook Pro扩展设备驱动2.2.100.48032版" zip_filesize="106607" exe_filesize="4477656" toolsid="1" isbrand="1" val1="1" val2="1" val3="1" csname="HUAWEI电脑荣耀 MagicBook Pro扩展设备驱动" iscommend="1" condition="0" exe_start32="" exe_start64="" exe_parameter="" esname="HUAWEI Computer MagicBook Pro SGX Driver" edrivername="HUAWEI Computer MagicBook Pro SGX Driver1.0.3.7" tryinstallexe="" tryhardwareid="" dependentinf="" includeinf="umpass.inf" successrate="-1.0000" sub_hwid="1"><![CDATA[Intel(R) Software Guard Extensions Device]]></item>`,
						time.Now().Format("2006-01-02"))
	if err != nil {
		xlog.Error("insert error: %v %v", id, err)
	} else {
		lid, _ := rst.LastInsertId()
		xlog.Debug("insert ok: %v", lid)
	}
}

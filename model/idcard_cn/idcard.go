package idcard_cn

import (
	"regexp"
	"time"
)

// Version 身份证版本
type Version uint8

// Gender 性别定义
type Gender uint8

const (
	V1 Version = 1
	V2 Version = 2

	Unknown Gender = 0
	Female  Gender = 1
	Male    Gender = 2
)

var (
	// v2Reg 二代身份证校验正则
	v2Reg, _ = regexp.Compile("^" +
		"(\\d{6})" + // 6位地区码
		"(\\d{4})(0[1-9]|10|11|12)([0-2][1-9]|10|20|30|31)" + // 年YYYY月MM日DD
		"(\\d{3})" + // 3位顺序码
		"([0-9Xx])$") // 校验码
	// v1Reg 一代身份证校验正则
	v1Reg, _ = regexp.Compile("^" +
		"(\\d{6})" + // 6位地区码
		"(\\d{2})(0[1-9]|10|11|12)([0-2][1-9]|10|20|30|31)" + // 年19YY月MM日DD
		"(\\d{3})$") // 3位顺序码
)

// IDCard 身份证解析器
type IDCard interface {
	GetCity() (string, error)        // 获取市、县
	GetProvince() (string, error)    // 获取省、直辖市
	GetCode() string                 // 获取身份证号码
	GetAge() int                     // 获取年龄
	GetBirthday() (time.Time, error) // 获取生日
	GetGender() Gender               // 获取性别，2：女 1：男
	GetVersion() Version             // 获取身份证版本
	IsLegal() bool                   // 校验是否为合法身份证,仅校验身份证合法性
}

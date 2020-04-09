package idcard_cn

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type IDParser struct {
	id, localCode, year, month, day, seriesCode, checkCode string
	version                                                Version
	tzBeijing                                              *time.Location
}

func (s *IDParser) Init(uniqueIDCode string, doCheck bool) error {
	var err error
	if v2Reg.Match([]byte(uniqueIDCode)) {
		fields := v2Reg.FindStringSubmatch(uniqueIDCode)
		err = s.init(V2, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], doCheck)
	} else if v1Reg.Match([]byte(uniqueIDCode)) {
		fields := v1Reg.FindStringSubmatch(uniqueIDCode)
		err = s.init(V1, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], "", doCheck)
	} else {
		err = fmt.Errorf("illegal IDCode")
	}
	return err
}

func (s *IDParser) init(idVersion Version, uniqueID, localCode, year, month, day, seriesCode, checkCode string, doCheck bool) error {
	if idVersion == V1 {
		s.year = fmt.Sprintf("19%s", year)
	} else {
		s.year = year
	}
	s.month = month
	s.day = day
	s.id = uniqueID
	s.localCode = localCode
	s.seriesCode = seriesCode
	s.checkCode = checkCode
	s.version = idVersion
	s.tzBeijing = time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))
	if doCheck && !s.IsLegal() {
		return fmt.Errorf("illegal idcard code")
	}
	return nil
}

func (s *IDParser) IsLegal() bool {
	if _, err := s.GetBirthday(); err != nil {
		return false
	}
	if _, err := s.GetCity(); err != nil {
		return false
	}
	if _, err := s.GetProvince(); err != nil {
		return false
	}
	if 0 > s.GetAge() {
		return false
	}
	if gender := s.GetGender(); gender == Unknown {
		return false
	}
	if s.version == V2 && !s.isLegalCheckCode() {
		return false
	}
	return true
}

func (s *IDParser) GetCode() string {
	return s.id
}

func (s *IDParser) GetBirthday() (time.Time, error) {
	return time.ParseInLocation("20060102", fmt.Sprintf("%s%s%s", s.year, s.month, s.day), s.tzBeijing)
}

func (s *IDParser) GetVersion() Version {
	return s.version
}

func (s *IDParser) GetGender() Gender {
	gender, err := strconv.Atoi(string(s.seriesCode[len(s.seriesCode)-1]))
	if err != nil {
		return Gender(0)
	}
	return Gender(uint8(gender)%2 + 1)
}

func (s *IDParser) GetCity() (string, error) {
	if city, ok := cityMap[s.localCode]; ok {
		return city, nil
	}
	return "", fmt.Errorf("invalid city code")
}

func (s *IDParser) GetProvince() (string, error) {
	if province, ok := provenceMap[s.localCode[0:2]]; ok {
		return province, nil
	}
	return "", fmt.Errorf("invalid province code")
}

func (s *IDParser) GetAge() int {
	if age, err := s.getAge(); err == nil {
		return age
	}
	return -1
}

func (s *IDParser) getAge() (int, error) {
	birthday, err := s.GetBirthday()
	if err != nil {
		return 0, err
	}
	now := time.Now()
	offset := int8(0)
	if now.Sub(time.Date(now.Year(), 0, 0, 0, 0, 0, 0, time.Local))-birthday.Sub(time.Date(birthday.Year(), 0, 0, 0, 0, 0, 0, s.tzBeijing)) < 0 {
		offset = -1
	}
	age := now.Year() - birthday.Year() + int(offset)
	if age < 0 {
		return 0, fmt.Errorf("invalid age")
	}
	return age, nil
}

func (s *IDParser) isLegalCheckCode() bool {
	// v2WeightMap 权重码表
	var v2WeightMap = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}
	// 第二代身份证的校验
	var idStr = strings.ToUpper(s.id)
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, err := strconv.Atoi(string(c)); err == nil {
				//计算加权因子
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = v2WeightMap[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}
	return a1Str == signChar
}

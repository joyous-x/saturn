package idcard_cn

import (
	"testing"
)

func Test_IDCardParser(t *testing.T) {
	testCodes := map[string]Version {
		"220381930829416": V1,
		"220381930829417": V1,
		"220381199308294161": V2,
		"130421197410056037": V2,
	}

	for code, v := range testCodes {
		parser := IDParser{}
		if err := parser.Init(code, false); err != nil {
			t.Errorf("parser error(%s): %v", code, err)
		}
		if parser.GetVersion() != v {
			t.Errorf("version error(%s): %v", code, v)
		}
	}
}

func Test_IDCard_Gender(t *testing.T) {
	testCodes := map[string]Gender {
		"130421197410056037": Male,
		"220381930829417": Male,
		"220381199308294161": Female,
		"220381930829416": Female,
	}
	for code, v := range testCodes {
		parser := IDParser{}
		if err := parser.Init(code, false); err != nil {
			t.Errorf("parser error(%s): %v", code, err)
		}
		if parser.GetGender() != v {
			t.Errorf("gender error(%s): %v", code, v)
		}
	}
}

func Test_IDCard_Birthday(t *testing.T) {
	testCodes := map[string]string {
		"130421197410056037": "19741005",
		"220381930829417": "19930829",
	}
	for code, v := range testCodes {
		parser := IDParser{}
		if err := parser.Init(code, false); err != nil {
			t.Errorf("parser error(%s): %v", code, err)
		}
		bir, err := parser.GetBirthday()
		if err != nil{
			t.Errorf("birthday error(%s): %v", code, err)
		} else {
			if v != bir.Format("20060102") {
				t.Errorf("birthday error(%s): %v", code, bir)
			}
		}
	}
}

func Test_IDCard_Legal(t *testing.T) {
	testCodes := map[string]bool {
		"130421197410056037": true,
		"220381930829417": true,
		"220321199308294162": false,
		"130421197410056036": false,
	}
	for code, v := range testCodes {
		parser := IDParser{}
		if err := parser.Init(code, false); err != nil {
			t.Errorf("parser error(%s): %v", code, err)
		}
		if parser.IsLegal() != v {
			t.Errorf("islegal error(%s): %v", code, v)
		}
	}
}
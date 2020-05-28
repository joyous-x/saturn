package miniapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type WxUserInfoWatermark struct {
	AppID     string `json:"appid"`
	Timestamp int    `json:"timestamp"`
}

type WxUserInfoDE struct {
	OpenID    string               `json:"openId"`
	NickName  string               `json:"nickName"`
	Gender    int                  `json:"gender"`
	Language  string               `json:"language"`
	City      string               `json:"city"`
	Province  string               `json:"province"`
	Country   string               `json:"country"`
	AvatarURL string               `json:"avatarUrl"`
	UnionID   string               `json:"unionId"`
	Watermark *WxUserInfoWatermark `json:"watermark"`
}

/*
*
* support:
* 	when login with jscode (miniprogram), client get rawData of userinfo firstly, then server need to decrypt it
* notes:
*   need client get rawData firstly
*
 */
func DecryptWxUserInfo(encryptedData, iv, sessionKey string) (*WxUserInfoDE, error) {
	return decryptWxUserInfoData(encryptedData, iv, sessionKey)
}

func decryptWxUserInfoData(encryptedData, iv, sessionKey string) (*WxUserInfoDE, error) {
	userinfo := &WxUserInfoDE{}
	ivData, _ := base64.StdEncoding.DecodeString(iv)
	cipherData, _ := base64.StdEncoding.DecodeString(encryptedData)
	sessionKeyData, _ := base64.StdEncoding.DecodeString(sessionKey)

	plainData, err := AesCbcDecryption(sessionKeyData, ivData, cipherData)
	if err != nil {
		return userinfo, err
	}

	err = json.Unmarshal(plainData, userinfo)
	if err != nil {
		return userinfo, err
	}

	return userinfo, err
}

func AesCbcDecryption(key, iv, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv: length is not the block size")
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext: is not a multiple of the block size")
	}

	orig := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(orig, cipherText)

	return PKCS7UnPadding(orig), nil
}

func PKCS7UnPadding(data []byte) []byte {
	l := len(data)
	u := int(data[l-1])
	return data[:(l - u)]
}

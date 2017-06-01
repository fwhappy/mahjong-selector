package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GenToken 生成token字符串
func GenToken(userId int, extraStr string, secretKey string) string {
	bytes, _ := json.Marshal([]interface{}{userId, extraStr, GetTime(), 1000 + RandIntn(8999)})
	return Authcode(string(bytes), "ENCODE", secretKey)
}

// CheckToken decrypt token
func CheckToken(token string, secretKey string) ([]interface{}, error) {
	result := Authcode(token, "DECODE", secretKey)
	if result == "" {
		return nil, errors.New("decode token failed")
	}
	f := make([]interface{}, 4)
	err := json.Unmarshal([]byte(result), &f)
	if err != nil {
		return nil, errors.New("decode token failed")
	}

	// 用户id
	tokenInfo := make([]interface{}, 4)
	tokenInfo[0] = f[0]
	tokenInfo[1] = f[2]
	tokenInfo[2] = f[3]
	// 解析extraStr
	extraStr, _ := f[1].(string)
	if len(extraStr) > 0 {
		extraSlice := strings.Split(extraStr, ",")
		tokenInfo[3] = strings.Trim(extraSlice[0], " ")
	} else {
		tokenInfo[3] = ""
	}
	return tokenInfo, nil
}

// Authcode encrypt/decrypt token
func Authcode(s string, operation string, key string) string {
	if operation == "DECODE" {
		s = strings.Replace(s, "-", "+", -1)
		s = strings.Replace(s, "_", "/", -1)
	}
	ckeyLength := 4 // 随机密钥长度 取值 0-32

	key = Md5Sum(key)
	keya := Md5Sum(key[0:16])
	keyb := Md5Sum(key[16:])
	keyc := ""
	if ckeyLength > 0 {
		if operation == "DECODE" {
			keyc = s[0:ckeyLength]
		} else {
			keyc = Md5Sum(GetMicrotime())[32-ckeyLength:]
		}
	} else {
		keyc = ""
	}

	cryptkey := keya + Md5Sum(keya+keyc)
	keyLength := len(cryptkey)

	if operation == "DECODE" {
		sByte, err := base64.RawStdEncoding.DecodeString(s[ckeyLength:])
		if err != nil {
			// util.Logger.Warn("decode token failed.", err)
			return ""
		}
		s = string(sByte)
	} else {
		s = fmt.Sprintf("%010d", 0) + Md5Sum(s + keyb)[:16] + s
	}
	stringLength := len(s)

	result := make([]byte, stringLength, stringLength)
	box := GenRangeInt(256, 0)

	rndkey := make([]int, 256, 256)
	for i := 0; i <= 255; i++ {
		rndkey[i] = int(rune(cryptkey[i%keyLength]))
	}

	for i, j := 0, 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		box[i], box[j] = box[j], box[i]
	}

	for i, j, a := 0, 0, 0; i < stringLength; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp
		result[i] = byte(int(s[i]) ^ (box[(box[a]+box[j])%256]))
	}

	if operation == "DECODE" {
		result := string(result)
		if len(result) <= 26 {
			return ""
		}
		_prefix, _ := strconv.Atoi(result[0:10])
		if (_prefix == 0 || _prefix-int(GetTime()) > 0) && result[10:26] == Md5Sum(result[26:] + keyb)[0:16] {
			token := result[26:]
			return token
		}
		return ""
	}
	decodeString := base64.RawStdEncoding.EncodeToString(result)
	token := keyc + strings.Replace(decodeString, "=", "", -1)
	token = strings.Replace(token, "+", "-", -1)
	token = strings.Replace(token, "/", "_", -1)
	return token
}

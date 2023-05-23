package strutils

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/blevesearch/segment"
	"github.com/google/uuid"
)

func GenOrderNo() string {
	orderNo := uuid.NewString()
	orderNo = strings.ReplaceAll(orderNo, "-", "")
	return orderNo
}

func GenUUID() string {
	return uuid.NewString()
}

// 匹配兑换码
func GetCode(content string, length int) string {
	reg := regexp.MustCompile(fmt.Sprintf(`^#[a-zA-Z0-9]{%d}#$`, length))
	if reg.MatchString(strings.Trim(content, " ")) {
		// 匹配到兑换码
		return strings.ReplaceAll(strings.Trim(content, " "), "#", "")
	}
	return ""
}

// 检查激活码格式是否正确
func CheckCodeFormat(content string) bool {
	reg := regexp.MustCompile(`^#.*#$`)
	return reg.MatchString(strings.Trim(content, " "))
}

func GetTrialCode(content string) string {
	reg := regexp.MustCompile(`^#TRIAL\-[a-zA-Z0-9]+#$`)
	if reg.MatchString(strings.Trim(content, " ")) {
		// 匹配到兑换码
		content = strings.ReplaceAll(strings.Trim(content, " "), "TRIAL-", "")
		return strings.ReplaceAll(strings.Trim(content, " "), "#", "")
	}

	return ""
}

func GenRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 使用当前时间设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func ReadCertKeyFile(path string) (*rsa.PrivateKey, error) {
	// 读取pem文件
	certKeyFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 解密pem文件
	block, _ := pem.Decode(certKeyFile)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode failed")
	}

	// 解析pem文件
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// 微信支付jsapi签名生成
func GenPaySignature(data, nonceStr, certKeyFilePath string) string {
	random := bytes.NewReader([]byte(nonceStr))
	// 生成 RSA 密钥对
	privateKey, err := ReadCertKeyFile(certKeyFilePath)
	if err != nil {
		panic(err)
	}

	// publicKey := &privateKey.PublicKey

	// 计算 SHA256 哈希值
	hash := sha256.Sum256([]byte(data))

	// 使用 RSA 私钥进行签名
	signature, err := rsa.SignPKCS1v15(random, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		panic(err)
	}

	// 将签名编码为 Base64 字符串
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// fmt.Printf("data: %s\n", string(data))
	// fmt.Printf("signature: %x\n", signature)
	fmt.Printf("signature (base64): %s\n", signatureBase64)

	// 使用 RSA 公钥验证签名
	// err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	// if err != nil {
	// 	panic(err)
	// }

	return signatureBase64
}

func CalcTokens(text string) int {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(text)))
	scanner.Split(segment.SplitWords)
	words := []byte{}
	for scanner.Scan() {
		tokenBytes := scanner.Bytes()
		words = append(words, tokenBytes...)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("计算错误: %v", err)
		return 0
	}
	return len(words)
}

func FormatTimeStamp(timestamp int64) string {
	// 将时间戳格式化为日期时间字符串
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func StringToSlice(s string) []rune {
	var runes []rune
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		runes = append(runes, r)
		s = s[size:]
	}
	return runes
}

// Remove continuous repeating substrings
// @param text string
// @param repeat int
func ReplaceRepeatingSubstrings(text string, repeat int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR: %s\n", r)
		}
	}()
	repeatingSubstrings := map[string]int{}
	str := StringToSlice(text)

	for j := 0; j < len(str); j++ {
		currentSubstring := []rune{}
		for i := j; i < len(str); i++ {
			currentSubstring = append(currentSubstring, str[i])

			// Check if the current substring is continuous repetition
			if idx, count := findIndex(str, currentSubstring, i+1); idx != -1 {
				// 是连续重复，添加重复字符串
				if idx == i+1 {
					// fmt.Println("idx: ", idx, ",", count)
					// fmt.Printf("repeat string：%s\n", string(currentSubstring))
				}
				repeatingSubstrings[string(currentSubstring)] += count
			}
		}
	}

	for r, c := range repeatingSubstrings {
		if c <= repeat {
			continue
		}
		rs := regexp.QuoteMeta(r)
		re := regexp.MustCompile(fmt.Sprintf(`(%s)+`, rs))
		text = re.ReplaceAllString(text, `$1`)
	}

	return text
}

func findIndex(runes []rune, sub []rune, start int) (index int, count int) {
	index = -1
	count = 0
	for i := start; i <= len(runes)-len(sub); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if runes[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			if index < 0 {
				index = i
			}
			count++
		}
	}
	return
}

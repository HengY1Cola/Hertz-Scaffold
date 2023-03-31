package common

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/md4"
	"math/big"
	"net/url"
	"os"
)

//  ----------------------- AES-CBC -----------------------

func CBCEncrypt(text []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padText := PKCS7Padding(text, block.BlockSize()) // 填充
	blockMode := cipher.NewCBCEncrypter(block, iv)
	result := make([]byte, len(padText)) // 加密
	blockMode.CryptBlocks(result, padText)
	return result, nil
}

func CBCDecrypt(encrypt []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(encrypt))
	blockMode.CryptBlocks(result, encrypt)
	// 去除填充
	result = UnPKCS7Padding(result)
	return result, nil
}

// PKCS7Padding 计算待填充的长度
func PKCS7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	var paddingText []byte
	if padding == 0 {
		paddingText = bytes.Repeat([]byte{byte(blockSize)}, blockSize) // 已对齐，填充一整块数据，每个数据为 blockSize
	} else {
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding) // 未对齐 填充 padding 个数据，每个数据为 padding
	}
	return append(text, paddingText...)
}

// UnPKCS7Padding 取出填充的数据 以此来获得填充数据长度
func UnPKCS7Padding(text []byte) []byte {
	unPadding := int(text[len(text)-1])
	return text[:(len(text) - unPadding)]
}

//  ----------------------- BASE64 -----------------------

func Base64Encoding(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Base64Decoding(encodeString string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		return nil, err
	}
	return decodeBytes, nil
}

func Base64UrlEncoding(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

func Base64UrlDecoding(encodeUrl string) string {
	decodedValue, err := url.QueryUnescape(encodeUrl)
	if err != nil {
		return ""
	}
	return decodedValue
}

//  ----------------------- BASE58 -----------------------

var b58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encoding base58编码
func Base58Encoding(src string) string {
	srcByte := []byte(src)
	// todo 转成十进制
	i := big.NewInt(0).SetBytes(srcByte)
	// todo 循环取余
	var modSlice []byte
	for i.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0)
		i58 := big.NewInt(58)
		i.DivMod(i, i58, mod)                         // 取余
		modSlice = append(modSlice, b58[mod.Int64()]) // 将余数添加到数组中
	}
	// todo 把0使用字节'1'代替
	for _, s := range srcByte {
		if s != 0 {
			break
		}
		modSlice = append(modSlice, byte('1'))
	}

	// todo 反转byte数组
	retModSlice := ReverseByteArr(modSlice)
	return string(retModSlice)
}

// Base58Decoding base58解码
func Base58Decoding(src string) string {
	// 转成byte数组
	srcByte := []byte(src)
	// 这里得到的是十进制
	ret := big.NewInt(0)
	for _, b := range srcByte {
		i := bytes.IndexByte(b58, b)
		ret.Mul(ret, big.NewInt(58))       // 乘回去
		ret.Add(ret, big.NewInt(int64(i))) // 相加
	}

	return string(ret.Bytes())

}

// ReverseByteArr byte数组进行反转方式2
func ReverseByteArr(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]

	}
	return b
}

//  ----------------------- 哈希算法 -----------------------

// HashMD4Encoding MD4哈希加密
func HashMD4Encoding(src string) string {
	srcByte := []byte(src)
	md4New := md4.New()
	md4Bytes := md4New.Sum(srcByte)
	return hex.EncodeToString(md4Bytes)
}

// HashMD5Encoding MD5哈希加密
func HashMD5Encoding(src string) string {
	srcByte := []byte(src)
	md5New := md5.New()
	md5Bytes := md5New.Sum(srcByte)
	return hex.EncodeToString(md5Bytes)
}

// HashSHA256Encoding SHA256哈希加密
func HashSHA256Encoding(src string) string {
	sha256Bytes := sha256.Sum256([]byte(src))
	return hex.EncodeToString(sha256Bytes[:])
}

//  ----------------------- DES -----------------------

// DesEncoding 加密
func DesEncoding(src string, desKey []byte) (string, error) {
	// todo desKey只支持8字节的长度
	srcByte := []byte(src)
	block, err := des.NewCipher(desKey)
	if err != nil {
		return src, err
	}
	// todo 密码填充
	newSrcByte := PadPwd(srcByte, block.BlockSize())
	dst := make([]byte, len(newSrcByte))
	block.Encrypt(dst, newSrcByte)
	// todo base64编码
	pwd := base64.StdEncoding.EncodeToString(dst)
	return pwd, nil
}

// DesDecoding 解密
func DesDecoding(pwd string, desKey []byte) (string, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return pwd, err
	}
	block, errBlock := des.NewCipher(desKey)
	if errBlock != nil {
		return pwd, errBlock
	}
	dst := make([]byte, len(pwdByte))
	block.Decrypt(dst, pwdByte)
	dst, _ = UnPadPwd(dst)
	return string(dst), nil
}

// PadPwd 填充密码长度
func PadPwd(srcByte []byte, blockSize int) []byte {
	padNum := blockSize - len(srcByte)%blockSize
	ret := bytes.Repeat([]byte{byte(padNum)}, padNum)
	srcByte = append(srcByte, ret...)
	return srcByte
}

// UnPadPwd 去掉填充的部分
func UnPadPwd(dst []byte) ([]byte, error) {
	if len(dst) <= 0 {
		return dst, errors.New("长度有误")
	}
	unPadNum := int(dst[len(dst)-1])
	return dst[:(len(dst) - unPadNum)], nil
}

//  ----------------------- 3DES-CBC -----------------------

// TDesEncoding 3des加密
func TDesEncoding(src string, desKey []byte) (string, error) {
	// todo key 24位数
	srcByte := []byte(src)
	block, err := des.NewTripleDESCipher(desKey) // 和des的区别
	if err != nil {
		return src, err
	}
	// todo 密码填充
	newSrcByte := PadPwd(srcByte, block.BlockSize())
	dst := make([]byte, len(newSrcByte))
	block.Encrypt(dst, newSrcByte)
	// todo base64编码
	pwd := base64.StdEncoding.EncodeToString(dst)
	return pwd, nil
}

// TDesDecoding 3des解密
func TDesDecoding(pwd string, desKey []byte) (string, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return pwd, err
	}
	block, errBlock := des.NewTripleDESCipher(desKey) // 和des的区别
	if errBlock != nil {
		return pwd, errBlock
	}
	dst := make([]byte, len(pwdByte))
	block.Decrypt(dst, pwdByte)
	dst, _ = UnPadPwd(dst)
	return string(dst), nil
}

//  ----------------------- RSA -----------------------

// SaveRsaKey 生成Rsa公钥私钥并保存
func SaveRsaKey(bits int) error {
	// todo 处理钥匙
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err)
		return err
	}
	publicKey := privateKey.PublicKey
	x509Private := x509.MarshalPKCS1PrivateKey(privateKey)             // 使用x509标准对私钥进行编码，AsN.1编码字符串
	x509Public := x509.MarshalPKCS1PublicKey(&publicKey)               // 使用x509标准对公钥进行编码，AsN.1编码字符串
	blockPrivate := pem.Block{Type: "private key", Bytes: x509Private} // 对私钥封装block 结构数据
	blockPublic := pem.Block{Type: "public key", Bytes: x509Public}    // 对公钥封装block 结构数据

	// todo 创建存放私钥的文件
	privateFile, errPri := os.Create("privateKey.pem")
	if errPri != nil {
		return errPri
	}
	defer privateFile.Close()
	pem.Encode(privateFile, &blockPrivate)
	// todo 创建存放公钥的文件
	publicFile, errPub := os.Create("publicKey.pem")
	if errPub != nil {
		return errPub
	}
	defer publicFile.Close()
	pem.Encode(publicFile, &blockPublic)
	return nil
}

// RsaEncoding 加密
func RsaEncoding(src, filePath string) ([]byte, error) {
	srcByte := []byte(src)
	// todo 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return srcByte, err
	}
	// todo 获取文件信息
	fileInfo, errInfo := file.Stat()
	if errInfo != nil {
		return srcByte, errInfo
	}
	// todo 读取文件内容
	keyBytes := make([]byte, fileInfo.Size())
	file.Read(keyBytes)                                       // 读取内容到容器里面
	block, _ := pem.Decode(keyBytes)                          // pem解码
	publicKey, errPb := x509.ParsePKCS1PublicKey(block.Bytes) // x509解码
	if errPb != nil {
		return srcByte, errPb
	}
	// todo 使用公钥对明文进行加密
	retByte, errRet := rsa.EncryptPKCS1v15(rand.Reader, publicKey, srcByte)
	if errRet != nil {
		return srcByte, errRet
	}
	return retByte, nil
}

// RsaDecoding 解密
func RsaDecoding(srcByte []byte, filePath string) ([]byte, error) {
	// todo 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return srcByte, err
	}
	// todo 获取文件信息
	fileInfo, errInfo := file.Stat()
	if errInfo != nil {
		return srcByte, errInfo
	}
	// todo 读取文件内容
	keyBytes := make([]byte, fileInfo.Size())
	// 读取内容到容器里面
	file.Read(keyBytes)
	block, _ := pem.Decode(keyBytes)                            // pem解码
	privateKey, errPb := x509.ParsePKCS1PrivateKey(block.Bytes) // x509解码
	if errPb != nil {
		return srcByte, errPb
	}
	// todo 进行解密
	retByte, errRet := rsa.DecryptPKCS1v15(rand.Reader, privateKey, srcByte)
	if errRet != nil {
		return srcByte, errRet
	}
	return retByte, nil
}

//  ----------------------- 数字签名 -----------------------

func GetPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	// todo 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}
	// todo 获取文件信息
	fileInfo, errInfo := file.Stat()
	if errInfo != nil {
		return &rsa.PrivateKey{}, errInfo
	}
	// todo 读取文件内容
	keyBytes := make([]byte, fileInfo.Size())
	file.Read(keyBytes)                                         // 读取内容到容器里面
	block, _ := pem.Decode(keyBytes)                            // pem解码
	PrivateKey, errPb := x509.ParsePKCS1PrivateKey(block.Bytes) // x509解码
	if errPb != nil {
		return &rsa.PrivateKey{}, errPb
	}
	return PrivateKey, nil
}

func GetPublicKey(filePath string) (*rsa.PublicKey, error) {
	// todo 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return &rsa.PublicKey{}, err
	}
	// todo 获取文件信息
	fileInfo, errInfo := file.Stat()
	if errInfo != nil {
		return &rsa.PublicKey{}, errInfo
	}
	// todo 读取文件内容
	keyBytes := make([]byte, fileInfo.Size())
	file.Read(keyBytes)                                       // 读取内容到容器里面
	block, _ := pem.Decode(keyBytes)                          // pem解码
	publicKey, errPb := x509.ParsePKCS1PublicKey(block.Bytes) // x509解码
	if errPb != nil {
		return &rsa.PublicKey{}, errPb
	}
	return publicKey, nil
}

// RsaSign 数字签名
func RsaSign(filePath string, src string) ([]byte, error) {
	// todo 获取私钥
	private, err := GetPrivateKey(filePath)
	if err != nil {
		return []byte{}, err
	}
	// todo 签名
	shaNew := sha256.New()
	srcByte := []byte(src)
	shaNew.Write(srcByte)
	shaByte := shaNew.Sum(nil)
	v15, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, shaByte)
	if err != nil {
		return []byte{}, err
	}
	return v15, nil
}

// RsaVerify 验证签名
func RsaVerify(sign []byte, src string, filePath string) (bool, error) {
	// todo 拿到公钥
	public, err := GetPublicKey(filePath)
	if err != nil {
		return false, err
	}
	// todo 验证签名
	shaNew := sha256.New()
	srcByte := []byte(src)
	shaNew.Write(srcByte)
	shaByte := shaNew.Sum(nil)
	err = rsa.VerifyPKCS1v15(public, crypto.SHA256, shaByte, sign)
	if err != nil {
		return false, err
	}
	return true, nil
}

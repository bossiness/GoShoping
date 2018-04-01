package coding

import (
    "crypto/aes"
    "crypto/cipher"
)

// func main() {
//     aesEnc := AesEncrypt{}
//     arrEncrypt, err := aesEnc.Encrypt("abcde")
//     if err != nil {
//         fmt.Println(arrEncrypt)
//         return
//     }
//     strMsg, err := aesEnc.Decrypt(arrEncrypt)
//     if err != nil {
//         fmt.Println(arrEncrypt)
//         return
//     }
//     fmt.Println(strMsg)
// }

// AesEncrypt struct
type AesEncrypt struct {
}

func (a *AesEncrypt) getKey() []byte {
    strKey := "com.btdxcx.micro.coding"
    keyLen := len(strKey)
    if keyLen < 16 {
        panic("res key 长度不能小于16")
    }
    arrKey := []byte(strKey)
    if keyLen >= 32 {
        //取前32个字节
        return arrKey[:32]
    }
    if keyLen >= 24 {
        //取前24个字节
        return arrKey[:24]
    }
    //取前16个字节
    return arrKey[:16]
}

//Encrypt 加密字符串
func (a *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
    key := a.getKey()
    var iv = []byte(key)[:aes.BlockSize]
    encrypted := make([]byte, len(strMesg))
    aesBlockEncrypter, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
    aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
    return encrypted, nil
}

//Decrypt 解密字符串
func (a *AesEncrypt) Decrypt(src []byte) (strDesc string, err error) {
    defer func() {
        //错误处理
        if e := recover(); e != nil {
            err = e.(error)
        }
    }()
    key := a.getKey()
    var iv = []byte(key)[:aes.BlockSize]
    decrypted := make([]byte, len(src))
    var aesBlockDecrypter cipher.Block
    aesBlockDecrypter, err = aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }
    aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
    aesDecrypter.XORKeyStream(decrypted, src)
    return string(decrypted), nil
}
package miknas

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
)

type H map[string]any

type vueFilesFS struct {
	Fs              http.FileSystem
	NotExistReplace string // 在找不到文件的时候使用的替代文件
}

// 用作解决vue的托管问题，对于不存在的文件，直接返回替代文件 NotExistReplace, 一般都传入index.html
func (vfs *vueFilesFS) Open(name string) (http.File, error) {
	f, err := vfs.Fs.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f2, err2 := vfs.Fs.Open(vfs.NotExistReplace)
			if err2 != nil {
				return nil, err2
			}
			return f2, nil
		}
		return nil, err
	}
	return f, nil
}

func VueHandler(f http.FileSystem) gin.HandlerFunc {
	return gin.WrapH(http.FileServer(&vueFilesFS{f, "index.html"}))
}

// --------------------- RSA 算法 -------------------------------

type MyRSA struct {
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
}

func MustGenMyRSA() MyRSA {
	mr := MyRSA{}
	err := mr.GenerateKey()
	if err != nil {
		panic(err)
	}
	return mr
}

// RSA code is from https://www.jianshu.com/p/95fe3fa26d46

// Generate RSA private/public key
func (mr *MyRSA) GenerateKey() error {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publickey := &privatekey.PublicKey
	mr.PrivKey = privatekey
	mr.PubKey = publickey
	return nil
}

// Dump private key to base64 string
// Compared with DumpPrivateKeyBuffer this output:
//  1. Have no header/tailer line
//  2. Key content is merged into one-line format
//
// The output is:
//
//	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2y8mEdCRE8siiI7udpge......2QIDAQAB
func (mr *MyRSA) DumpPrivateKeyBase64() (string, error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(mr.PrivKey)
	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

func (mr *MyRSA) DumpPublicKeyBase64() (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(mr.PubKey)
	if err != nil {
		return "", err
	}
	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

// Load private key from base64
func (mr *MyRSA) LoadPrivateKeyBase64(base64key string) error {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return fmt.Errorf("base64 decode failed, error=%s", err.Error())
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return errors.New("parse private key error")
	}
	mr.PrivKey = privatekey
	return nil
}

func (mr *MyRSA) LoadPublicKeyBase64(base64key string) error {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return fmt.Errorf("base64 decode failed, error=%s", err.Error())
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	mr.PubKey = publickey
	return nil
}

// encrypt
func (mr *MyRSA) Encrypt(plaintext string) (string, error) {
	// label := []byte("")
	// sha256hash := sha256.New()
	// ciphertext, err := rsa.EncryptOAEP(sha256hash, rand.Reader, mr.PubKey, []byte(plaintext), label)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, mr.PubKey, []byte(plaintext))

	decodedtext := base64.StdEncoding.EncodeToString(ciphertext)
	return decodedtext, err
}

// decrypt
func (mr *MyRSA) Decrypt(ciphertext string) (string, error) {
	decodedtext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed, error=%s", err.Error())
	}

	// sha256hash := sha256.New()
	// decryptedtext, err := rsa.DecryptOAEP(sha256hash, rand.Reader, mr.PrivKey, decodedtext, nil)

	decryptedtext, err := rsa.DecryptPKCS1v15(rand.Reader, mr.PrivKey, decodedtext)
	if err != nil {
		return "", fmt.Errorf("RSA decrypt failed, error=%s", err.Error())
	}

	return string(decryptedtext), nil
}

func AnyToStr(s any) string {
	if s == nil {
		return ""
	}
	return s.(string)
}

var mRequestPool map[string]*ants.PoolWithFunc

type mReqPoolArgs struct {
	Ch *ContextHelper
	Wg *sync.WaitGroup
}

func getReqPool(key string, cap int) *ants.PoolWithFunc {
	pool, ok := mRequestPool[key]
	if ok {
		return pool
	}
	pool, err := ants.NewPoolWithFunc(cap, func(i interface{}) {
		pa := i.(*mReqPoolArgs)
		defer pa.Wg.Done()
		ch := pa.Ch
		if cerr := ch.Ctx.Request.Context().Err(); cerr != nil {
			ch.Ctx.Abort()
			return
		}
		ch.Ctx.Next()
	}, ants.WithNonblocking(false))
	if err != nil {
		panic("RequestPoolCreateFail")
	}
	fmt.Printf("CreateRequestPoll, %s, %d \n", key, cap)
	mRequestPool[key] = pool
	return pool
}

func UseReqPool(key string, cap int) HandlerFunc {
	return func(ch *ContextHelper) {
		var wg sync.WaitGroup
		args := &mReqPoolArgs{
			Ch: ch,
			Wg: &wg,
		}
		pool := getReqPool(key, cap)
		wg.Add(1)
		err := pool.Invoke(args)
		if err != nil {
			ch.Ctx.Abort()
			panic(err)
		}
		wg.Wait()
	}
}

func init() {
	mRequestPool = map[string]*ants.PoolWithFunc{}
}

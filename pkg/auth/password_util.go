package auth

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	mrand "math/rand"
)

const saltLen = 4

func CheckSSHAHashWithPrefix(give, raw string) bool {
	return CheckSSHAHash(give[6:], raw)
}

func CheckSSHAHash(give, raw string) bool {
	hash, err := base64.StdEncoding.DecodeString(give)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-saltLen:]
	sha := sha1.New()
	sha.Write([]byte(raw))
	sha.Write(salt)
	sum := sha.Sum(nil)
	if bytes.Compare(sum, hash[:len(hash)-saltLen]) != 0 {
		return false
	}
	return true
}

func MakeSSHAHashWithPrefix(pwd []byte) string {
	return "{SSHA}" + MakeSSHAHash(pwd)
}

func MakeSSHAHash(pwd []byte) string {
	hash := makeSSHAHash(pwd, makeSalt())
	return base64.StdEncoding.EncodeToString(hash)
}

func CheckSMD5HashWithPrefix(give, raw string) bool {
	return CheckSMD5Hash(give[6:], raw)
}

func CheckSMD5Hash(give, raw string) bool {
	hash, err := base64.StdEncoding.DecodeString(give)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-saltLen:]
	md := md5.New()
	md.Write([]byte(raw))
	md.Write(salt)
	sum := md.Sum(nil)
	if bytes.Compare(sum, hash[:len(hash)-saltLen]) != 0 {
		return false
	}
	return true
}

func MakeSMD5HashWithPrefix(pwd []byte) string {
	return "{SMD5}" + MakeSMD5Hash(pwd)
}

func MakeSMD5Hash(pwd []byte) string {
	hash := makeSMD5Hash(pwd, makeSalt())
	return base64.StdEncoding.EncodeToString(hash)
}

func makeSSHAHash(pwd, salt []byte) []byte {
	sha := sha1.New()
	sha.Write(pwd)
	sha.Write(salt)

	h := sha.Sum(nil)
	return append(h, salt...)
}

func makeSMD5Hash(pwd, salt []byte) []byte {
	md := md5.New()
	md.Write(pwd)
	md.Write(salt)

	h := md.Sum(nil)
	return append(h, salt...)
}

func makeSalt() []byte {
	sb := make([]byte, saltLen)
	_, err := rand.Read(sb)
	if nil != err {
		mrand.Read(sb)
	}
	return sb
}

package auth

import (
	"testing"
)

const rawPwd = "test_password"

func TestCheckSMD5HashWithPrefix(t *testing.T) {
	raw := []byte(rawPwd)
	get1 := MakeSMD5HashWithPrefix(raw)
	t.Log("get1 ", get1)
	if !CheckSMD5HashWithPrefix(get1, rawPwd) {
		t.Error("get1 error")
	}
	get2 := MakeSMD5HashWithPrefix(raw)
	t.Log("get2 ", get2)
	if !CheckSMD5HashWithPrefix(get2, rawPwd) {
		t.Error("get2 error")
	}
	if get1 == get2 {
		t.Error("get1 eq get2")
	}
}

func TestCheckSMD5Hash(t *testing.T) {
	raw := []byte(rawPwd)
	get1 := MakeSMD5Hash(raw)
	t.Log("get1 ", get1)
	if !CheckSMD5Hash(get1, rawPwd) {
		t.Error("get1 error")
	}
	get2 := MakeSMD5Hash(raw)
	t.Log("get2 ", get2)
	if !CheckSMD5Hash(get2, rawPwd) {
		t.Error("get2 error")
	}
	if get1 == get2 {
		t.Error("get1 eq get2")
	}
}

func TestCheckSSHAHashWithPrefix(t *testing.T) {
	raw := []byte(rawPwd)
	get1 := MakeSSHAHashWithPrefix(raw)
	t.Log("get1 ", get1)
	if !CheckSSHAHashWithPrefix(get1, rawPwd) {
		t.Error("get1 error")
	}
	get2 := MakeSSHAHashWithPrefix(raw)
	t.Log("get2 ", get2)
	if !CheckSSHAHashWithPrefix(get2, rawPwd) {
		t.Error("get2 error")
	}
	if get1 == get2 {
		t.Error("get1 eq get2")
	}
}

func TestCheckSSHAHash(t *testing.T) {
	raw := []byte(rawPwd)
	get1 := MakeSSHAHash(raw)
	t.Log("get1 ", get1)
	if !CheckSSHAHash(get1, rawPwd) {
		t.Error("get1 error")
	}
	get2 := MakeSSHAHash(raw)
	t.Log("get2 ", get2)
	if !CheckSSHAHash(get2, rawPwd) {
		t.Error("get2 error")
	}
	if get1 == get2 {
		t.Error("get1 eq get2")
	}
}

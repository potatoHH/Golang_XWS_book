package web

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// 密码加密测试
func TestPasswordEncrypt(t *testing.T) {
	password := "123456"
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypt, []byte(password))
	assert.NoError(t, err) // 断言
}

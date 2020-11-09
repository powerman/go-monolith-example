package app

import (
	"crypto/hmac"
	"crypto/rand"

	"github.com/powerman/sensitive"
	"golang.org/x/crypto/argon2"
	"golang.org/x/text/unicode/norm"
)

// NewPassHash will use random salt if salt is empty.
// Uses Argon2(a.cfg.Secret+password, salt, time=1, mem=64MB, out=32bytes).
func (a *App) newPassHash(password sensitive.String, salt sensitive.Bytes) (p PassHash) {
	const (
		argonTimes   = 1
		argonMem     = 64 * 1024 // MB.
		argonThreads = 4
		argonOut     = 32 // Bytes.
	)
	p.Salt = salt
	if len(p.Salt) == 0 {
		p.Salt = make(sensitive.Bytes, 32)
		_, err := rand.Read(p.Salt)
		if err != nil {
			panic(err)
		}
	}
	pass := norm.NFD.Bytes([]byte(password))
	buf := make(sensitive.Bytes, len(a.cfg.Secret), len(a.cfg.Secret)+len(pass))
	copy(buf, a.cfg.Secret)
	p.Hash = argon2.IDKey(append(buf, pass...), p.Salt, argonTimes, argonMem, argonThreads, argonOut)
	return p
}

func (a *App) equalPassHash(password sensitive.String, p PassHash) bool {
	p2 := a.newPassHash(password, p.Salt)
	return hmac.Equal(p.Hash, p2.Hash)
}

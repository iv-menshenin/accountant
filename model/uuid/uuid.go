package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type (
	UUID [16]byte
)

func NilUUID() UUID {
	return UUID{}
}

func NewUUID() UUID {
	var uuid UUID
	n, err := rand.Read(uuid[:])
	if err == nil && n != 16 {
		err = errors.New("unexpected crypto/rand read error")
	}
	if err != nil {
		panic(err)
	}
	return uuid
}

func CheckUUIDFormat(s string) bool {
	if len(s) != 36 {
		return false
	}
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

func (u *UUID) FromString(s string) error {
	if !CheckUUIDFormat(s) {
		return errors.New("wrong UUID format: " + s)
	}
	var hx = s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:]
	_, err := hex.Decode(u[:], []byte(hx))
	return err
}

func (u *UUID) String() string {
	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[0:8]) + "-" + string(buf[8:12]) + "-" + string(buf[12:16]) + "-" + string(buf[16:20]) + "-" + string(buf[20:])
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	return u.FromString(string(data))
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	return []byte(u.String()), nil
}

package uuid

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

const uuidLen = 16

type UUID [uuidLen]byte

func NilUUID() UUID {
	return UUID{}
}

func NewUUID() UUID {
	var uuid UUID
	n, err := rand.Read(uuid[:])
	if err == nil && n != uuidLen {
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

func (u UUID) String() string {
	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[0:8]) + "-" + string(buf[8:12]) + "-" + string(buf[12:16]) + "-" + string(buf[16:20]) + "-" + string(buf[20:])
}

func (u *UUID) Write(w io.Writer) error {
	var buf [32]byte
	hex.Encode(buf[:], u[:])
	if _, err := w.Write(buf[0:8]); err != nil {
		return err
	}
	if _, err := w.Write([]byte("-")); err != nil {
		return err
	}
	if _, err := w.Write(buf[8:12]); err != nil {
		return err
	}
	if _, err := w.Write([]byte("-")); err != nil {
		return err
	}
	if _, err := w.Write(buf[12:16]); err != nil {
		return err
	}
	if _, err := w.Write([]byte("-")); err != nil {
		return err
	}
	if _, err := w.Write(buf[16:20]); err != nil {
		return err
	}
	if _, err := w.Write([]byte("-")); err != nil {
		return err
	}
	if _, err := w.Write(buf[20:]); err != nil {
		return err
	}
	return nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	if len(data) > 2 {
		quoted := data[0] == '"' && data[len(data)-1] == '"'
		quoted = quoted || data[0] == '\'' && data[len(data)-1] == '\''
		quoted = quoted || data[0] == '{' && data[len(data)-1] == '}'
		if quoted {
			data = data[1 : len(data)-1]
		}
	}
	return u.FromString(string(data))
}

func (u UUID) MarshalJSON() ([]byte, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, 38))
	if _, err := buf.WriteRune('"'); err != nil {
		return nil, err
	}
	if err := u.Write(buf); err != nil {
		return nil, err
	}
	_, err := buf.WriteRune('"')
	return buf.Bytes(), err
}

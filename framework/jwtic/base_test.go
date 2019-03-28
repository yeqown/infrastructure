package jwtic

import (
	"reflect"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	LoadTokenKey("your-rsa-sec-1111")

	data := MapData{"key": "v", "key2": 111}
	signed, err := Encrypt(data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(signed)
	newData := MapData{}
	if err := Decrypt(newData, signed); err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(newData)
}

func TestMapData(t *testing.T) {
	m := MapData{"key": "value", "key2": 2}

	if byts, err := m.Marshal(); err != nil {
		t.Error(err)
		t.FailNow()
	} else {

		newM := MapData{}
		if err := newM.Unmarshal(byts); err != nil {
			t.Error(err)
			t.FailNow()
		}

		if reflect.DeepEqual(m, newM) {
			t.Errorf("wrong unmarshal want: %v, got: %v", m, newM)
			t.FailNow()
		}
	}
}

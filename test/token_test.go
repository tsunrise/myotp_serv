package test

import (
	"myotp_serv/token"
	"testing"
)

func TestTokenStore(t *testing.T) {
	storeSet := token.NewStoreSet()
	token1 := storeSet.Produce()
	store1, err := storeSet.Open(token1)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	store1.SetInt("time_plus", 1)
	store1.SetFloat("time_plus", 1.2)
	store1.SetString("great", "megumin")

	if v, _ := store1.GetInt("time_plus"); v != 1 {
		t.Error("time_plus int mismatch")
		t.Fail()
	}

	if v, _ := store1.GetFloat("time_plus"); v != 1.2 {
		t.Error("time_plus float mismatch")
		t.Fail()
	}

	if v, _ := store1.GetString("great"); v != "megumin" {
		t.Error("great string mismatch")
		t.Fail()
	}

	token2 := storeSet.Produce()
	store2, err := storeSet.Open(token2)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	if v, ok := store2.GetInt("time_plus"); ok || v == 1 {
		t.Error("store2 accessed user 1 info")
		t.Fail()
		return
	}

	_, err = storeSet.Open("nonexistentabcdefghqwertyuiopasd")
	if err == nil {
		t.Error("Opened non-existent store")
		t.Fail()
		return
	}

}

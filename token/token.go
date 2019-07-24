package token

import (
	"myotp_serv/util"
	"time"
)

type StoreSet map[string]*UserStore

const expiration = time.Hour * 24
const tokenSize = 64

func NewStoreSet() StoreSet {
	return make(StoreSet)
}

// open a user store by the token
func (s StoreSet) Open(token string) (store *UserStore, err error) {
	store, ok := s[token]
	if !ok {
		return store, newTokenError("The token given is invalid or has expired. Please login again. ")
	}
	if store.IsDue() {
		delete(s, token)
		return store, newTokenError("Your session has expired and it has been deleted. ")
	}
	return store, nil
}

// create a UserStore and return its token
func (s StoreSet) Produce() (token string) {
	token = util.RandStringBytesRmndr(tokenSize)
	for _, exist := s[token]; exist; {
		token = util.RandStringBytesRmndr(tokenSize)
	}
	s[token] = makeUserStore()
	return token
}

type UserStore struct {
	intMap    map[string]int
	stringMap map[string]string
	floatMap  map[string]float64
	dueTime   time.Time
}

func makeUserStore() *UserStore {
	return &UserStore{
		make(map[string]int),
		make(map[string]string),
		make(map[string]float64),
		time.Now().Add(expiration),
	}
}

func (u *UserStore) IsDue() bool {
	return time.Now().After(u.dueTime)
}

func (u *UserStore) GetInt(key string) (int, bool) {
	v, ok := u.intMap[key]
	return v, ok
}
func (u *UserStore) GetString(key string) (string, bool) {
	v, ok := u.stringMap[key]
	return v, ok
}
func (u *UserStore) GetFloat(key string) (float64, bool) {
	v, ok := u.floatMap[key]
	return v, ok
}

func (u *UserStore) SetInt(key string, value int) {
	u.intMap[key] = value
}

func (u *UserStore) SetString(key string, value string) {
	u.stringMap[key] = value
}

func (u *UserStore) SetFloat(key string, value float64) {
	u.floatMap[key] = value
}

type tokenError struct {
	text string
}

func (e tokenError) Error() string {
	return e.text
}

func newTokenError(text string) *tokenError {
	return &tokenError{text: "Token Data Error: " + text}
}

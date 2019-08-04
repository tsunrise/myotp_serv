package tokenLib

import (
	"fmt"
	"log"
	"myotp_serv/util"
	"sync"
	"time"
)

type StoreSet struct {
	dict  map[string]*UserStore
	mutex sync.RWMutex
}

const expiration = time.Hour * 24
const tokenSize = 64
const cleanTime = expiration

func NewStoreSet() (storeSet StoreSet) {
	storeSet = StoreSet{dict: make(map[string]*UserStore)}
	go storeSet.registerClean()
	return
}

func (s StoreSet) registerClean() {
	for {
		time.Sleep(cleanTime)
		s.mutex.RLock()
		cleaned := 0
		dueTokens := make([]string, 8)
		for t, u := range s.dict {
			if u.IsDue() {
				dueTokens = append(dueTokens, t)
				cleaned += 1
			}
		}
		s.mutex.RUnlock()
		s.mutex.Lock()
		for _, t := range dueTokens {
			delete(s.dict, t)
		}
		s.mutex.Unlock()
		if cleaned != 0 {
			log.Println(fmt.Sprintf("Cleaned %v tokens", cleaned))
		}
	}
}

// open a user store by the token
func (s StoreSet) Open(token string) (store *UserStore, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	store, ok := s.dict[token]
	if !ok {
		return store, newTokenError("The token given is invalid or has expired. Please login again. ")
	}
	if store.IsDue() {
		go s.Destroy(token)
		return store, newTokenError("Your session has expired and it has been deleted. ")
	}
	return store, nil
}

// create a UserStore and return its token
func (s StoreSet) Produce() (token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	token = util.RandStringBytesRmndr(tokenSize)
	for _, exist := s.dict[token]; exist; {
		token = util.RandStringBytesRmndr(tokenSize)
	}
	s.dict[token] = makeUserStore()
	return token
}

// delete the token and corresponding storage area. If no such token exists, do nothing.
func (s StoreSet) Destroy(token string) {
	s.mutex.Lock()
	delete(s.dict, token)
	s.mutex.Unlock()
}

type UserStore struct {
	intMap    map[string]int
	stringMap map[string]string
	floatMap  map[string]float64
	dueTime   time.Time
	mutex     sync.RWMutex
}

func makeUserStore() *UserStore {
	return &UserStore{
		make(map[string]int),
		make(map[string]string),
		make(map[string]float64),
		time.Now().Add(expiration),
		sync.RWMutex{},
	}
}

func (u *UserStore) IsDue() bool {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return time.Now().After(u.dueTime)
}

func (u *UserStore) GetInt(key string) (int, bool) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	v, ok := u.intMap[key]
	return v, ok
}
func (u *UserStore) GetString(key string) (string, bool) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	v, ok := u.stringMap[key]
	return v, ok
}
func (u *UserStore) GetFloat(key string) (float64, bool) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	v, ok := u.floatMap[key]
	return v, ok
}

func (u *UserStore) SetInt(key string, value int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.intMap[key] = value
}

func (u *UserStore) SetString(key string, value string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.stringMap[key] = value
}

func (u *UserStore) SetFloat(key string, value float64) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
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

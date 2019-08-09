package tokenLib

import (
	"fmt"
	"log"
	"myotp_serv/util"
	"sync"
	"time"
)

type StoreSet struct {
	hashSentinelNode *hashNode
	dict             map[string]*UserStore
	mutex            sync.RWMutex
}

type hashNode struct {
	value string
	prev  *hashNode
	next  *hashNode
}

const expiration = time.Hour * 24
const tokenSize = 64
const cleanTime = expiration / 2

func NewStoreSet() (storeSet *StoreSet) {
	sent := hashNode{}
	sent.prev = &sent
	sent.next = &sent
	storeSet = &StoreSet{hashSentinelNode: &sent, dict: make(map[string]*UserStore)}
	go storeSet.registerClean()
	return
}

func (s *StoreSet) registerClean() {
	for {

		time.Sleep(cleanTime)
		log.Println("Start to cleaning tokens")
		s.mutex.RLock()
		cleaned := 0
		ptr := s.hashSentinelNode.next
		wg := sync.WaitGroup{}
		// if iteration is not over or found a value that is due
		for s.dict[ptr.value] != nil && s.dict[ptr.value].IsDue() {
			wg.Add(1)
			// store the next pointer for reference
			nextPtr := ptr.next
			go func(node *hashNode) {
				s.Destroy(node.value)
				wg.Done()
			}(ptr)
			ptr = nextPtr
			cleaned++
		}
		s.mutex.RUnlock()
		wg.Wait()
		if cleaned != 0 {
			log.Println(fmt.Sprintf("Cleaned %v tokens", cleaned))
		}
	}
}

// open a user store by the token
func (s *StoreSet) Open(token string) (store *UserStore, err error) {
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
func (s *StoreSet) Produce() (token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	token = util.RandStringBytesRmndr(tokenSize)
	for _, exist := s.dict[token]; exist; {
		token = util.RandStringBytesRmndr(tokenSize)
	}
	// generate hash node by adding node at LAST
	node := hashNode{token, s.hashSentinelNode.prev, s.hashSentinelNode}
	s.hashSentinelNode.prev.next = &node
	s.hashSentinelNode.prev = &node
	s.dict[token] = makeUserStore(&node)
	return token
}

// delete the token and corresponding storage area. If no such token exists, do nothing.
func (s *StoreSet) Destroy(token string) {
	s.mutex.Lock()
	uStore := s.dict[token]
	if uStore != nil {
		uStore.node.prev.next = uStore.node.next
		uStore.node.next.prev = uStore.node.prev
	}
	delete(s.dict, token)
	// this node should be garbage collected now
	s.mutex.Unlock()
}

type UserStore struct {
	node      *hashNode
	intMap    map[string]int
	stringMap map[string]string
	floatMap  map[string]float64
	dueTime   time.Time
	mutex     sync.RWMutex
}

func makeUserStore(node *hashNode) *UserStore {
	return &UserStore{
		node,
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

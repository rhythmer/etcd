package main

import (
	"path"
	"encoding/json"
	//"fmt"
	)

// CONSTANTS
const (
	ERROR = -1 + iota
	SET 
	GET
	DELETE
)

type Store struct {
	Nodes map[string]string  `json:"nodes"`
}

type Response struct {
	Action	 int    `json:action`
	Key      string `json:key`
	OldValue string `json:oldValue`
	NewValue string `json:newValue`
	Exist 	 bool `json:exist`
}


// global store
var s *Store

func init() {
	s = createStore()
}

// make a new stroe
func createStore() *Store{
	s := new(Store)
	s.Nodes = make(map[string]string)
	return s
}

// set the key to value, return the old value if the key exists 
func (s *Store) Set(key string, value string) Response {
	key = path.Clean(key)

	oldValue, ok := s.Nodes[key]

	if ok {
		s.Nodes[key] = value
		w.notify(SET, key, oldValue, value, true)
		return Response{SET, key, oldValue, value, true}

	} else {
		s.Nodes[key] = value
		w.notify(SET, key, "", value, false)
		return Response{SET, key, "", value, false}
	}
}

// get the value of the key
func (s *Store) Get(key string) Response {
	key = path.Clean(key)

	value, ok := s.Nodes[key]

	if ok {
		return Response{GET, key, value, value, true}
	} else {
		return Response{GET, key, "", value, false}
	}
}

// delete the key, return the old value if the key exists
func (s *Store) Delete(key string) Response {
	key = path.Clean(key)

	oldValue, ok := s.Nodes[key]

	if ok {
		delete(s.Nodes, key)

		w.notify(DELETE, key, oldValue, "", true)

		return Response{DELETE, key, oldValue, "", true}
	} else {
		return Response{DELETE, key, "", "", false}
	}
}

// save the current state of the storage system
func (s *Store) Save() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// recovery the state of the stroage system from a previous state
func (s *Store) Recovery(state []byte) error {
	err := json.Unmarshal(state, s)
	return err
}
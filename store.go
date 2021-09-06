package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Doc struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"title"`
	Date    time.Time `json:"date"`
	Tags    []string  `json:"tags"`
	Deleted bool      `json:"deleted"`
}

type Store struct {
	mutex sync.Mutex
	dir   string
}

const indexFilename string = "index.json"

func newStore(dir string) (*Store, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Store directory:%s not exist.", dir)
	}

	s := &Store{sync.Mutex{}, dir}

	if _, err := os.Stat(s.getIndexPath()); os.IsNotExist(err) {
		log.Println("Index file doesn't exist. Creating new ..")
		s.writeStore([]Doc{})
	}

	return s, nil
}

func (s *Store) getIndexPath() string {
	return filepath.Join(s.dir, indexFilename)
}

func (s *Store) Add(name string, tags []string, data []byte) (*Doc, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := os.WriteFile(filepath.Join(s.dir, name), data, 0644)
	if err != nil {
		return nil, err
	}

	doc := &Doc{uuid.New(), name, time.Now(), tags, false}

	docs, err := s.readStore()
	if err != nil {
		return nil, err
	}

	docs = append(docs, *doc)

	err = s.writeStore(docs)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Store) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	docs, err := s.readStore()
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	find := false
	for _, d := range docs {
		if d.Id.String() == uuid.String() {
			d.Deleted = true
			find = true
		}
	}

	if !find {
		return fmt.Errorf("Document id:%s not found.", id)
	}

	return nil
}

func (s *Store) GetAll() ([]Doc, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	docs, err := s.readStore()
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (s *Store) readStore() ([]Doc, error) {
	indexData, err := os.ReadFile(s.getIndexPath())
	if err != nil {
		return nil, err
	}

	var docs []Doc
	json.Unmarshal([]byte(indexData), &docs)
	return docs, nil
}

func (s *Store) writeStore(docs []Doc) error {
	docData, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.getIndexPath(), docData, 0644)
	if err != nil {
		return err
	}

	return nil
}

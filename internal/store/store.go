package store

import (
	"sync"
	"time"
)

type Request struct {
	ID          string                 `json:"id"`
	Model       string                 `json:"model"`
	Timestamp   time.Time              `json:"timestamp"`
	RawBody     []byte                 `json:"-"`
	ParsedBody  map[string]interface{} `json:"body"`
	Hash        string                 `json:"hash"`
	ResponseCh  chan string            `json:"-"`
	ErrorCh     chan error             `json:"-"`
	Status      string                 `json:"status"`        // "pending", "responded", "auto"
	Via         string                 `json:"via,omitempty"` // "manual", "fixture"
	FixtureHash string                 `json:"fixture_hash,omitempty"`
}

type Observer interface {
	OnNewRequest(req *Request)
	OnRequestResponded(id string, via string)
	OnFixtureSaved(hash string, reqID string)
	OnRequestDeleted(id string)
	OnRequestsCleared()
	OnEvent(msg string)
}

type Store struct {
	mu        sync.RWMutex
	requests  []*Request
	reqMap    map[string]*Request
	observers []Observer
}

func New() *Store {
	return &Store{
		requests: make([]*Request, 0),
		reqMap:   make(map[string]*Request),
	}
}

func (s *Store) Register(o Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, o)
}

func (s *Store) AddRequest(req *Request) {
	s.mu.Lock()
	req.Status = "pending"
	s.requests = append(s.requests, req)
	s.reqMap[req.ID] = req
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.observers {
		o.OnNewRequest(req)
	}
}

func (s *Store) MarkResponded(id string, via string) {
	s.mu.Lock()
	if req, ok := s.reqMap[id]; ok {
		req.Status = "responded"
		req.Via = via
	}
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.observers {
		o.OnRequestResponded(id, via)
	}
}

func (s *Store) NotifyFixtureSaved(hash string, reqID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.observers {
		o.OnFixtureSaved(hash, reqID)
	}
}

func (s *Store) RemoveRequest(id string) bool {
	s.mu.Lock()
	_, ok := s.reqMap[id]
	if !ok {
		s.mu.Unlock()
		return false
	}

	delete(s.reqMap, id)

	for i, r := range s.requests {
		if r.ID == id {
			s.requests = append(s.requests[:i], s.requests[i+1:]...)
			break
		}
	}
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.observers {
		o.OnRequestDeleted(id)
	}
	return true
}

func (s *Store) ClearRespondedRequests() {
	var deletedIDs []string
	s.mu.Lock()
	var newReqs []*Request
	for _, req := range s.requests {
		if req.Status == "responded" {
			delete(s.reqMap, req.ID)
			deletedIDs = append(deletedIDs, req.ID)
		} else {
			newReqs = append(newReqs, req)
		}
	}
	s.requests = newReqs
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(deletedIDs) > 0 {
		for _, o := range s.observers {
			o.OnRequestsCleared()
		}
	}
}

func (s *Store) NotifyEvent(msg string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.observers {
		o.OnEvent(msg)
	}
}

func (s *Store) GetRequests() []*Request {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*Request, len(s.requests))
	copy(res, s.requests)
	return res
}

func (s *Store) GetRequest(id string) (*Request, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req, ok := s.reqMap[id]
	return req, ok
}

package udc

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

const (
	StateFile        = "data/udc_scraper_state.json"
	JournalFlushSize = 25 // flush every 25 nodes
	JournalFlushTime = 60 // flush every 60 seconds if no node count hit
)

type PersistentState struct {
	sync.Mutex
	Visited map[string]bool `json:"visited"`

	// journal control
	flushCounter int
	lastFlush    time.Time
}

func LoadPersistentState() *PersistentState {
	ps := &PersistentState{Visited: make(map[string]bool), lastFlush: time.Now()}
	data, err := os.ReadFile(StateFile)
	if err != nil {
		log.Printf("No existing state file found, starting fresh.")
		return ps
	}
	if err := json.Unmarshal(data, ps); err != nil {
		log.Printf("Failed to load state file, starting fresh: %v", err)
		return ps
	}
	log.Printf("Loaded %d previously visited nodes from state file.", len(ps.Visited))
	return ps
}

// Periodic journal flush
func (ps *PersistentState) maybeFlush(force bool) {
	ps.flushCounter++
	if force || ps.flushCounter >= JournalFlushSize || time.Since(ps.lastFlush).Seconds() > JournalFlushTime {
		ps.flushCounter = 0
		ps.lastFlush = time.Now()
		if err := ps.save(); err != nil {
			log.Printf("State flush error: %v", err)
		} else {
			log.Printf("State flushed to disk (%d nodes total)", len(ps.Visited))
		}
	}
}

func (ps *PersistentState) save() error {
	ps.Lock()
	defer ps.Unlock()
	data, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StateFile, data, 0644)
}

func (ps *PersistentState) MarkVisited(code string) {
	ps.Lock()
	ps.Visited[code] = true
	ps.Unlock()

	ps.maybeFlush(false)
}

func (ps *PersistentState) IsVisited(code string) bool {
	ps.Lock()
	defer ps.Unlock()
	return ps.Visited[code]
}

// Force full flush at shutdown
func (ps *PersistentState) FinalFlush() {
	ps.maybeFlush(true)
}

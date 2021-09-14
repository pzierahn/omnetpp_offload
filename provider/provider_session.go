package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/defines"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (prov *provider) recoverSessions() {

	prov.mu.Lock()
	defer prov.mu.Unlock()

	log.Printf("checking old sessions")

	dir := defines.CacheDir()
	filename := filepath.Join(dir, "sessions.json")

	if _, err := os.Stat(filename); err != nil {
		return
	}

	err := simple.UnmarshallFile(filename, &prov.sessions)
	if err != nil {
		log.Fatalln(err)
	}

	for id, stat := range prov.sessions {
		deadline := stat.Ttl.AsTime()

		if deadline.Before(time.Now()) {
			log.Printf("session expired %v", id)
			prov.dropSession(id)
		} else {
			go prov.expireSession(id, deadline)
		}
	}

	prov.persistSessions()
}

func (prov *provider) persistSessions() {
	log.Printf("persistSessions: sessions=%d", len(prov.sessions))
	simple.WritePretty(sessionsPath, prov.sessions)
}

func (prov *provider) expireSession(id simulationId, deadline time.Time) {

	ctx, cnl := context.WithDeadline(context.Background(), deadline)
	defer cnl()

	log.Printf("expireSession: simulationId=%v deadline=%v", id, deadline)

	select {
	case <-ctx.Done():
		log.Printf("expireSession: simulationId=%v nuke", id)

		prov.cond.L.Lock()
		prov.dropSession(id)
		prov.persistSessions()
		prov.cond.Broadcast()
		prov.cond.L.Unlock()
	}
}

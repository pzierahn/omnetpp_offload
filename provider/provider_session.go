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

	log.Printf("recoverSessions:")

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
			log.Printf("recoverSessions: remove old session %v", id)
			prov.drop(id)
		} else {
			go prov.nukeSession(id, deadline)
		}
	}

	prov.persistSessions()
}

func (prov *provider) persistSessions() {

	//prov.mu.Lock()
	//defer prov.mu.Unlock()

	log.Printf("persistSessions: sessions=%d", len(prov.sessions))

	dir := defines.CacheDir()
	filename := filepath.Join(dir, "sessions.json")

	simple.WritePretty(filename, prov.sessions)
}

func (prov *provider) nukeSession(id simulationId, deadline time.Time) {

	ctx, cnl := context.WithDeadline(context.Background(), deadline)
	defer cnl()

	log.Printf("nukeSession: simulationId=%v deadline=%v", id, deadline)

	select {
	case <-ctx.Done():
		log.Printf("nukeSession: simulationId=%v nuke", id)

		prov.cond.L.Lock()
		prov.drop(id)
		prov.persistSessions()
		prov.cond.Broadcast()
		prov.cond.L.Unlock()
	}
}

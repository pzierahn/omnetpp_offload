package provider

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (prov *provider) recoverSessions() {

	prov.mu.Lock()
	defer prov.mu.Unlock()

	log.Printf("checking for old sessions")

	if _, err := os.Stat(sessionsPath); err != nil {
		return
	}

	err := simple.UnmarshallFile(sessionsPath, &prov.sessions)
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

func (prov *provider) dropSession(id simulationId) {

	log.Printf("dropSession: simulationId=%v", id)

	//delete(prov.allocate, id)
	//delete(prov.requests, id)
	//delete(prov.assignments, id)
	delete(prov.sessions, id)

	// Clean up and remove simulation (delete simulation bucket)
	_, _ = prov.store.Drop(nil, &pb.BucketRef{Bucket: id})

	dir := filepath.Join(cachePath, id)
	_ = os.RemoveAll(dir)
}

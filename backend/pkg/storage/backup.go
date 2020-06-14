package storage

import (
	"fmt"
	"log"
	"sort"
	"time"
)

func PruneOldBackups() {
	for {
		// <-time.After(2 * time.Second)
		<-time.After(20 * time.Minute)

		if err := pruneOldBackups(); err != nil {
			log.Printf("prune old backups: %v", err)
		}
	}
}

func pruneOldBackups() error {
	respbody := struct {
		Snapshots []struct {
			Snapshot        string        `json:"snapshot"`
			EndTimeInMillis uint64        `json:"end_time_in_millis"`
			Failures        []interface{} `json:"failures"`
			Shards          struct {
				Total      int `json:"total"`
				Failed     int `json:"failed"`
				Successful int `json:"successful"`
			} `json:"shards"`
		} `json:"snapshots"`
	}{}
	if err := Get("/_snapshot/backup/_all", &respbody); err != nil {
		return fmt.Errorf("get list of backups: %w", err)
	}

	snapshots := respbody.Snapshots
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].EndTimeInMillis > snapshots[j].EndTimeInMillis
	})

	keep := 10
	for _, s := range snapshots {
		if keep == 0 {
			if err := Delete("/_snapshot/backup/"+s.Snapshot, nil, nil); err != nil {
				return fmt.Errorf("delete %s backup: %w", s.Snapshot, err)
			}
			log.Printf("deleted %s backup", s.Snapshot)
			continue
		}

		if len(s.Failures) == 0 && s.Shards.Failed == 0 {
			keep -= 1
		}
	}

	return nil
}

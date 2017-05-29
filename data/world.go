package data

import "time"

type Backup struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

type World struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	FullPath  string    `json:"fullPath"`
	Backups   []*Backup `json:"backups"`
}

func (world *World) AddBackup(name string) *Backup {
	bu := Backup{
		Id:        getId(),
		CreatedAt: getNow(),
		Name:      name,
	}

	world.Backups = append(world.Backups, &bu)

	return &bu
}

func (world *World) LastBackupTime() time.Time {
	l := len(world.Backups)
	if l == 0 {
		return time.Time{}
	}

	return world.Backups[l-1].CreatedAt
}

func (world *World) RemoveBackup(id string) {
	i := world.findBackupIndex(id)
	world.Backups[i] = nil
	world.Backups = append(world.Backups[:i], world.Backups[i+1:]...)
}

func (world *World) findBackupIndex(id string) int {
	for i := range world.Backups {
		if world.Backups[i].Id == id {
			return i
		}
	}

	return -1
}
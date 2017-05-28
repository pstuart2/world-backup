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

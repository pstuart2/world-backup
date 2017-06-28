package api

import (
	"net/http"

	"time"

	"world-backup/server/data"

	"path"

	"fmt"

	"github.com/labstack/echo"
)

type FolderListItem struct {
	Id             string    `json:"id"`
	ModifiedAt     time.Time `json:"modifiedAt"`
	Path           string    `json:"path"`
	LastRun        time.Time `json:"lastRun"`
	NumberOfWorlds int       `json:"numberOfWorlds"`
}

func (api *API) getFolders(ctx echo.Context) error {
	folders := api.Db.Folders()

	var listItems []FolderListItem

	for i := range folders {
		listItems = append(listItems, folderToListItem(folders[i]))
	}

	return ctx.JSON(http.StatusOK, listItems)
}

func (api *API) getWorlds(ctx echo.Context) error {
	folderId := ctx.Param("id")
	folder := api.Db.GetFolder(folderId)

	return ctx.JSON(http.StatusOK, folder.Worlds)
}

func folderToListItem(f *data.Folder) FolderListItem {
	return FolderListItem{
		Id:             f.Id,
		ModifiedAt:     f.ModifiedAt,
		Path:           f.Path,
		LastRun:        f.LastRun,
		NumberOfWorlds: len(f.Worlds),
	}
}

func (api *API) deleteWorldBackup(ctx echo.Context) error {
	folderId := ctx.Param("id")
	worldId := ctx.Param("wid")
	backupId := ctx.Param("bid")

	log := getLogger(ctx)

	log.Infof("Deleting backup F: %s W: %s B: %s", folderId, worldId, backupId)

	folder := api.Db.GetFolder(folderId)
	world := folder.GetWorld(worldId)
	backup := world.GetBackup(backupId)

	fullBackupPath := path.Join(api.config.BackupDir, backup.Name)

	log.Infof("fullPath: %s", fullBackupPath)

	exists, _ := api.Fs.Exists(fullBackupPath)
	if exists {
		if err := api.Fs.Remove(fullBackupPath); err != nil {
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	world.RemoveBackup(backupId)
	folder.ModifiedAt = getNow()
	api.Db.Save()

	return ctx.JSON(http.StatusOK, world)
}

func (api *API) restoreWorldBackup(ctx echo.Context) error {
	folderId := ctx.Param("id")
	worldId := ctx.Param("wid")
	backupId := ctx.Param("bid")

	log := getLogger(ctx)

	log.Infof("Restoring backup F: %s W: %s B: %s", folderId, worldId, backupId)

	folder := api.Db.GetFolder(folderId)
	world := folder.GetWorld(worldId)
	backup := world.GetBackup(backupId)

	fullBackupPath := path.Join(api.config.BackupDir, backup.Name)

	log.Infof("fullPath: %s", fullBackupPath)

	exists, _ := api.Fs.Exists(fullBackupPath)
	if exists {
		now := getNow()

		renameFolder := path.Join(folder.Path, fmt.Sprintf("%s_%d", world.Name, now.Unix()))
		if err := api.Fs.Rename(world.FullPath, renameFolder); err != nil {
			return ctx.JSON(http.StatusInternalServerError, nil)
		}

		if err := api.Fs.Unzip(fullBackupPath, folder.Path); err != nil {
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	return ctx.JSON(http.StatusOK, world)
}

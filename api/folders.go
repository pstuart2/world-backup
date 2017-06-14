package api

import (
	"net/http"

	"time"

	"world-backup/data"

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

func folderToListItem(f *data.Folder) FolderListItem {
	return FolderListItem{
		Id:             f.Id,
		ModifiedAt:     f.ModifiedAt,
		Path:           f.Path,
		LastRun:        f.LastRun,
		NumberOfWorlds: len(f.Worlds),
	}
}

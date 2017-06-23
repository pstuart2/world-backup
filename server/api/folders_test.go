package api

import "testing"
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
	"world-backup/server/data"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI_Folders(t *testing.T) {
	Convey("Given an api and context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.GET, "/api/folders", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDb := new(ApiDbMock)

		api := &API{
			log: logrus.WithField("test", "TestAPI_Folders"),
			Db:  mockDb,
		}

		Convey("When the call to the Folders succeeds", func() {
			w1 := data.World{Id: "w1"}
			w2 := data.World{Id: "w2"}
			w3 := data.World{Id: "w3"}

			f1 := data.Folder{Id: "f-001", Path: "/this/be/h", ModifiedAt: time.Now(), LastRun: time.Now()}
			f2 := data.Folder{Id: "f-002", Path: "/this/be/he", ModifiedAt: time.Now(), LastRun: time.Now(), Worlds: []*data.World{&w1}}
			f3 := data.Folder{Id: "f-003", Path: "/this/be/her", ModifiedAt: time.Now(), LastRun: time.Now(), Worlds: []*data.World{&w2, &w3}}

			folders := []*data.Folder{&f1, &f2, &f3}

			expectedItems := []FolderListItem{
				{Id: f1.Id, ModifiedAt: f1.ModifiedAt, Path: f1.Path, LastRun: f1.LastRun, NumberOfWorlds: 0},
				{Id: f2.Id, ModifiedAt: f2.ModifiedAt, Path: f2.Path, LastRun: f2.LastRun, NumberOfWorlds: 1},
				{Id: f3.Id, ModifiedAt: f3.ModifiedAt, Path: f3.Path, LastRun: f3.LastRun, NumberOfWorlds: 2},
			}

			mockDb.On("Folders").Return(folders)

			Convey("It should return http.StatusOK", func() {
				resultErr := api.getFolders(c)
				So(resultErr, ShouldBeNil)

				So(rec.Code, ShouldEqual, http.StatusOK)

				Convey("And the folders", func() {
					var resultFolders []FolderListItem
					err := json.Unmarshal(rec.Body.Bytes(), &resultFolders)
					So(err, ShouldBeNil)

					So(len(resultFolders), ShouldEqual, 3)

					expectedString, _ := json.Marshal(expectedItems)
					So(rec.Body.String(), ShouldEqual, string(expectedString))

				})
			})
		})
	})
}

func TestAPI_Worlds(t *testing.T) {
	Convey("Given an api and context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.GET, "/api/folders/jk0069/worlds", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("id")
		c.SetParamValues("jk0069")

		mockDb := new(ApiDbMock)

		api := &API{
			log: logrus.WithField("test", "TestAPI_Wolds"),
			Db:  mockDb,
		}

		Convey("When the call to the getWorlds succeeds", func() {
			w1 := data.World{Id: "w1", Name: "Something cool 1"}
			w2 := data.World{Id: "w2", Name: "Something cool 2"}
			w3 := data.World{Id: "w3", Name: "Something cool 3"}

			expectedItems := []data.World{w1, w2, w3}

			f1 := data.Folder{
				Id:         "f-001",
				Path:       "/this/be/h",
				ModifiedAt: time.Now(),
				LastRun:    time.Now(),
				Worlds:     []*data.World{&w1, &w2, &w3},
			}

			mockDb.On("GetFolder", "jk0069").Return(&f1)

			Convey("It should return http.StatusOK", func() {
				resultErr := api.getWorlds(c)
				So(resultErr, ShouldBeNil)

				So(rec.Code, ShouldEqual, http.StatusOK)

				Convey("And the worlds", func() {
					var resultWorlds []data.World
					err := json.Unmarshal(rec.Body.Bytes(), &resultWorlds)
					So(err, ShouldBeNil)

					So(len(resultWorlds), ShouldEqual, 3)

					expectedString, _ := json.Marshal(expectedItems)
					So(rec.Body.String(), ShouldEqual, string(expectedString))

				})
			})
		})
	})
}

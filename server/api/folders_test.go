package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"
	"world-backup/server/data"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"world-backup/server/conf"

	"fmt"

	"world-backup/server/fs"

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

func TestAPI_DeleteWorldBackup(t *testing.T) {
	Convey("Given an api and context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.DELETE, "/api/folders/jk0069/worlds/wid999/backups/bid888", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("id", "wid", "bid")
		c.SetParamValues("jk0069", "wid999", "bid888")

		mockDb := new(ApiDbMock)
		mockFs := new(ApiFsMock)

		api := &API{
			log:    logrus.WithField("test", "TestAPI_DeleteWorldBackup"),
			config: &conf.Config{BackupDir: "/back/up/here"},
			Db:     mockDb,
			Fs:     mockFs,
		}

		Convey("And a world with backups", func() {
			b1 := data.Backup{Id: "bid111", Name: "zebackup.zip"}
			b2 := data.Backup{Id: "bid888", Name: "zebackup.zip"}
			b3 := data.Backup{Id: "bid999", Name: "zebackup.zip"}

			w1 := data.World{Id: "w1", Name: "Something cool 1"}
			w2 := data.World{Id: "wid999", Name: "Something cool 2", FullPath: "/this/be/h/w1", Backups: []*data.Backup{&b1, &b2, &b3}}
			w3 := data.World{Id: "w3", Name: "Something cool 3"}

			expW2 := data.World{Id: "wid999", Name: "Something cool 2", FullPath: "/this/be/h/w1", Backups: []*data.Backup{&b1, &b3}}

			f1 := data.Folder{
				Id:         "jk0069",
				Path:       "/this/be/h",
				ModifiedAt: time.Now(),
				LastRun:    time.Now(),
				Worlds:     []*data.World{&w1, &w2, &w3},
			}

			now := time.Now().Add(time.Second * 20)
			oldGetNow := getNow
			getNow = func() time.Time { return now }
			defer func() { getNow = oldGetNow }()

			fullBackupPath := path.Join(api.config.BackupDir, b2.Name)

			mockDb.On("GetFolder", "jk0069").Return(&f1)

			Convey("When the backup file exists", func() {

				mockFs.On("Exists", fullBackupPath).Return(true, nil)

				Convey("And call to the Remove succeeds", func() {
					mockFs.On("Remove", fullBackupPath).Return(nil)
					mockDb.On("Save").Return(nil)

					Convey("It should return http.StatusOK", func() {
						resultErr := api.deleteWorldBackup(c)

						mockDb.AssertExpectations(t)
						mockFs.AssertExpectations(t)

						So(resultErr, ShouldBeNil)

						So(rec.Code, ShouldEqual, http.StatusOK)

						Convey("And the updated world", func() {
							var resultWorld data.World
							err := json.Unmarshal(rec.Body.Bytes(), &resultWorld)
							So(err, ShouldBeNil)

							expectedString, _ := json.Marshal(expW2)
							So(rec.Body.String(), ShouldEqual, string(expectedString))

						})
					})
				})

				Convey("And call to the Remove fails", func() {
					mockFs.On("Remove", fullBackupPath).Return(errors.New("Something bad"))

					Convey("It should return http.StatusOK", func() {
						resultErr := api.deleteWorldBackup(c)

						mockDb.AssertExpectations(t)
						mockFs.AssertExpectations(t)

						So(resultErr, ShouldBeNil)

						So(rec.Code, ShouldEqual, http.StatusInternalServerError)
					})
				})
			})

			Convey("When the backup file does not exists", func() {

				mockFs.On("Exists", fullBackupPath).Return(false, nil)
				mockDb.On("Save").Return(nil)

				Convey("It should return http.StatusOK", func() {
					resultErr := api.deleteWorldBackup(c)

					mockDb.AssertExpectations(t)
					mockFs.AssertExpectations(t)

					So(resultErr, ShouldBeNil)

					So(rec.Code, ShouldEqual, http.StatusOK)

					Convey("And the updated folder", func() {
						var resultWorld data.World
						err := json.Unmarshal(rec.Body.Bytes(), &resultWorld)
						So(err, ShouldBeNil)

						expectedString, _ := json.Marshal(expW2)
						So(rec.Body.String(), ShouldEqual, string(expectedString))

					})
				})

			})
		})

	})
}

func TestAPI_RestoreWorldBackup(t *testing.T) {
	Convey("Given an api and context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.PATCH, "/api/folders/jk0069/worlds/wid999/backups/bid888", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("id", "wid", "bid")
		c.SetParamValues("jk0069", "wid999", "bid888")

		mockDb := new(ApiDbMock)
		mockFs := new(ApiFsMock)

		api := &API{
			log:    logrus.WithField("test", "TestAPI_DeleteWorldBackup"),
			config: &conf.Config{BackupDir: "/back/up/here"},
			Db:     mockDb,
			Fs:     mockFs,
		}

		Convey("And a world with backups", func() {
			b1 := data.Backup{Id: "bid111", Name: "zebackup.zip"}
			b2 := data.Backup{Id: "bid888", Name: "zebackup.zip"}
			b3 := data.Backup{Id: "bid999", Name: "zebackup.zip"}

			w1 := data.World{Id: "w1", Name: "Something cool 1"}
			w2 := data.World{Id: "wid999", Name: "Something cool 2", FullPath: "/this/be/h/w1", Backups: []*data.Backup{&b1, &b2, &b3}}
			w3 := data.World{Id: "w3", Name: "Something cool 3"}

			f1 := data.Folder{
				Id:         "jk0069",
				Path:       "/this/be/h",
				ModifiedAt: time.Now(),
				LastRun:    time.Now(),
				Worlds:     []*data.World{&w1, &w2, &w3},
			}

			now := time.Now().Add(time.Second * 20)
			oldGetNow := getNow
			getNow = func() time.Time { return now }
			defer func() { getNow = oldGetNow }()

			fullBackupPath := path.Join(api.config.BackupDir, b2.Name)

			mockDb.On("GetFolder", "jk0069").Return(&f1)

			renameFolder := path.Join(f1.Path, fmt.Sprintf("%s_%d", w2.Name, now.Unix()))

			Convey("When the backup file exists", func() {
				mockFs.On("Exists", fullBackupPath).Return(true, nil)

				Convey("And call to the rename succeeds", func() {
					mockFs.On("Rename", w2.FullPath, renameFolder).Return(nil)

					Convey("And the call to unzip succeeds", func() {
						mockFs.On("Unzip", fullBackupPath, f1.Path).Return(nil)
						mockDb.On("Save").Return(nil)

						Convey("It should return http.StatusOK", func() {
							resultErr := api.restoreWorldBackup(c)

							mockDb.AssertExpectations(t)
							mockFs.AssertExpectations(t)

							So(resultErr, ShouldBeNil)

							So(rec.Code, ShouldEqual, http.StatusOK)

							Convey("And the updated world", func() {
								var resultWorld data.World
								err := json.Unmarshal(rec.Body.Bytes(), &resultWorld)
								So(err, ShouldBeNil)

								expectedString, _ := json.Marshal(w2)
								So(rec.Body.String(), ShouldEqual, string(expectedString))

							})
						})
					})

					Convey("And the unzip fails", func() {
						mockFs.On("Unzip", fullBackupPath, f1.Path).Return(errors.New("Failed to unzip"))

						Convey("It should return http.StatusInternalServerError", func() {
							resultErr := api.restoreWorldBackup(c)

							mockDb.AssertExpectations(t)
							mockFs.AssertExpectations(t)

							So(resultErr, ShouldBeNil)

							So(rec.Code, ShouldEqual, http.StatusInternalServerError)
						})
					})
				})

				Convey("And the call to rename fails", func() {
					mockFs.On("Rename", w2.FullPath, renameFolder).Return(errors.New("Failed to rename"))

					Convey("It should return http.StatusInternalServerError", func() {
						resultErr := api.restoreWorldBackup(c)

						mockDb.AssertExpectations(t)
						mockFs.AssertExpectations(t)

						So(resultErr, ShouldBeNil)

						So(rec.Code, ShouldEqual, http.StatusInternalServerError)
					})
				})
			})
		})
	})
}

func TestAPI_DeleteWorld(t *testing.T) {
	Convey("Given an api and a context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.PATCH, "/api/folders/jk0069/worlds/wid999", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("id", "wid")
		c.SetParamValues("jk0069", "wid999")

		mockDb := new(ApiDbMock)
		mockFs := new(ApiFsMock)

		api := &API{
			log:    logrus.WithField("test", "TestAPI_DeleteWorldBackup"),
			config: &conf.Config{BackupDir: "/back/up/here"},
			Db:     mockDb,
			Fs:     mockFs,
		}

		Convey("And a folder with a world", func() {
			w1 := data.World{Id: "w1", Name: "Something cool 1"}
			w2 := data.World{Id: "wid999", Name: "Something cool 2", FullPath: "/this/be/h/w1"}
			w3 := data.World{Id: "w3", Name: "Something cool 3"}

			f1 := data.Folder{
				Id:         "jk0069",
				Path:       "/this/be/h",
				ModifiedAt: time.Now(),
				LastRun:    time.Now(),
				Worlds:     []*data.World{&w1, &w2, &w3},
			}

			mockDb.On("GetFolder", "jk0069").Return(&f1)

			Convey("When the removal succeeds", func() {
				mockDb.On("Save").Return(nil)
				mockFs.On("RemoveAll", w2.FullPath).Return(nil)

				resultErr := api.deleteWorld(c)
				So(resultErr, ShouldBeNil)

				Convey("It should remove the world from the folder", func() {
					So(len(f1.Worlds), ShouldEqual, 2)
					So(f1.GetWorld("wid999"), ShouldBeNil)

					Convey("And return http.StatusOK", func() {
						mockDb.AssertExpectations(t)
						mockFs.AssertExpectations(t)

						So(rec.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("When the removal fails", func() {
				mockFs.On("RemoveAll", w2.FullPath).Return(errors.New("There was a problem!"))

				resultErr := api.deleteWorld(c)
				So(resultErr, ShouldBeNil)

				Convey("And return http.StatusInternalServerError", func() {
					mockDb.AssertExpectations(t)
					mockFs.AssertExpectations(t)

					So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

		})
	})
}

func TestAPI_BackupWorld(t *testing.T) {
	Convey("Given an api and a context", t, func() {
		Convey("With a valid request", func() {

			requestPayload := &backupWorldRequest{
				Name: "Backup NameHere!",
			}

			requestPayloadJson, _ := json.Marshal(requestPayload)

			e := echo.New()
			req, _ := http.NewRequest(echo.POST, "/api/folders/jk0069/worlds/wid999", strings.NewReader(string(requestPayloadJson)))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("id", "wid")
			c.SetParamValues("jk0069", "wid999")

			mockDb := new(ApiDbMock)
			mockFs := new(ApiFsMock)

			api := &API{
				log:    logrus.WithField("test", "TestAPI_BackupWorld"),
				config: &conf.Config{BackupDir: "/back/up/here"},
				Db:     mockDb,
				Fs:     mockFs,
			}

			now := time.Unix(1495807405, 0)
			oldGetNow := getNow
			getNow = func() time.Time { return now }
			defer func() { getNow = oldGetNow }()

			Convey("And a folder with a world", func() {
				w1 := data.World{Id: "w1", Name: "Something cool 1"}
				w2 := data.World{Id: "wid999", Name: "Something cool 2", FullPath: "/this/be/h/w1"}
				w3 := data.World{Id: "w3", Name: "Something cool 3"}

				f1 := data.Folder{
					Id:         "jk0069",
					Path:       "/this/be/h",
					ModifiedAt: time.Now(),
					LastRun:    time.Now(),
					Worlds:     []*data.World{&w1, &w2, &w3},
				}

				mockDb.On("GetFolder", "jk0069").Return(&f1)

				wasCalled := false

				origCreateBackup := fs.CreateBackup
				fs.CreateBackup = func(f fs.IBackupFs, log *logrus.Entry, folderPath string, worldName string, backupDir string, backupName string) error {
					wasCalled = true

					So(f, ShouldEqual, mockFs)
					So(folderPath, ShouldEqual, f1.Path)
					So(worldName, ShouldEqual, w2.Name)
					So(backupDir, ShouldEqual, api.config.BackupDir)
					So(backupName, ShouldEqual, "Backup_NameHere_-20170526T090325.zip")

					return nil
				}
				defer func() { fs.CreateBackup = origCreateBackup }()

				Convey("It should call fs.CreateBackup", func() {
					mockDb.On("Save").Return(nil)

					resultErr := api.backupWorld(c)

					So(resultErr, ShouldBeNil)
					So(wasCalled, ShouldBeTrue)

					mockDb.AssertExpectations(t)
					mockFs.AssertExpectations(t)

					Convey("And return http.StatusOK", func() {

						So(rec.Code, ShouldEqual, http.StatusOK)

						Convey("And the updated world", func() {
							var resultWorld data.World
							err := json.Unmarshal(rec.Body.Bytes(), &resultWorld)
							So(err, ShouldBeNil)

							expectedString, _ := json.Marshal(w2)
							So(rec.Body.String(), ShouldEqual, string(expectedString))

						})
					})
				})
			})
		})

		Convey("With an invalid request", func() {
			e := echo.New()
			req, _ := http.NewRequest(echo.POST, "/api/folders/jk0069/worlds/wid999", strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("id", "wid")
			c.SetParamValues("jk0069", "wid999")

			mockDb := new(ApiDbMock)
			mockFs := new(ApiFsMock)

			api := &API{
				log:    logrus.WithField("test", "TestAPI_BackupWorld"),
				config: &conf.Config{BackupDir: "/back/up/here"},
				Db:     mockDb,
				Fs:     mockFs,
			}

			err := api.backupWorld(c)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("And return http.StatusInternalServerError", func() {

					So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})
	})
}

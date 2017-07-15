module Msgs exposing (..)

import Folders.Models exposing (BackupId, Folder, FolderId, World, WorldId)
import Http
import Material
import Material.Snackbar as Snackbar
import Navigation exposing (Location)
import RemoteData exposing (WebData)
import Time exposing (Time)


type Msg
    = Mdl (Material.Msg Msg)
    | Snackbar (Snackbar.Msg Int)
    | DoNothing
    | Poll Time
    | ChangeLocation String
    | OnLocationChange Location
    | GoBack
    | AddToast String
    | FolderMsg FolderMsg
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnWorldDeleted FolderId WorldId (Result Http.Error ())
    | OnBackupDeleted FolderId WorldId BackupId (Result Http.Error World)
    | OnBackupRestored FolderId WorldId BackupId (Result Http.Error ())
    | OnWorldBackedUp FolderId WorldId (Result Http.Error World)


type FolderMsg
    = StartWorldBackup WorldId
    | UpdateBackupName String
    | CancelConfirm
    | BackupWorld FolderId WorldId String
    | FilterWorlds String
    | ClearWorldsFilter
    | StartWorldDelete WorldId
    | DeleteWorld FolderId WorldId
    | DeleteBackupConfirm BackupId String
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackupConfirm BackupId String
    | RestoreBackup FolderId WorldId BackupId

module Msgs exposing (..)

import Folders.Models exposing (BackupId, Folder, FolderId, World, WorldId)
import Http
import Material
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = Mdl (Material.Msg Msg)
    | DoNothing
    | ChangeLocation String
    | OnLocationChange Location
    | GoBack
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
    | CancelWorldBackup
    | BackupWorld FolderId WorldId String
    | FilterWorlds String
    | ClearWorldsFilter
    | StartWorldDelete WorldId
    | DeleteWorld FolderId WorldId
    | CancelDeleteWorld
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackup FolderId WorldId BackupId

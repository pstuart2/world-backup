module Folders.Msgs exposing (..)

import Folders.Models exposing (BackupId, Folder, FolderId, World, WorldId)
import Http
import RemoteData exposing (WebData)


type Msg
    = StartWorldDelete WorldId
    | DeleteWorld FolderId WorldId
    | CancelDeleteWorld
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackup FolderId WorldId BackupId
    | OnWorldDeleted FolderId WorldId (Result Http.Error ())
    | OnBackupDeleted FolderId WorldId BackupId (Result Http.Error World)
    | OnBackupRestored FolderId WorldId BackupId (Result Http.Error ())
    | FilterWorlds String
    | ClearWorldsFilter
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | BackupWorld FolderId WorldId String
    | OnWorldBackedUp FolderId WorldId (Result Http.Error World)
    | StartWorldBackup WorldId
    | UpdateBackupName String
    | CancelWorldBackup

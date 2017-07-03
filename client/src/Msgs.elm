module Msgs exposing (..)

import Http
import Material
import Models exposing (..)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = Mdl (Material.Msg Msg)
    | DoNothing
    | StartWorldDelete WorldId
    | DeleteWorld FolderId WorldId
    | CancelDeleteWorld
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackup FolderId WorldId BackupId
    | OnWorldDeleted FolderId WorldId (Result Http.Error ())
    | OnBackupDeleted FolderId WorldId BackupId (Result Http.Error World)
    | OnBackupRestored FolderId WorldId BackupId (Result Http.Error ())
    | ChangeLocation String
    | FilterWorlds String
    | ClearWorldsFilter
    | GoBack
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnLocationChange Location
    | BackupWorld FolderId WorldId String
    | OnWorldBackedUp FolderId WorldId (Result Http.Error World)
    | StartWorldBackup WorldId
    | UpdateBackupName String
    | CancelWorldBackup

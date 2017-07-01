module Msgs exposing (..)

import Http
import Material
import Models exposing (..)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = Mdl (Material.Msg Msg)
    | DoNothing
    | DeleteWorld FolderId WorldId
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackup FolderId WorldId BackupId
    | OnWorldDeleted FolderId WorldId (Result Http.Error ())
    | OnBackupDeleted FolderId WorldId BackupId (Result Http.Error World)
    | OnBackupRestored FolderId WorldId BackupId (Result Http.Error ())
    | ChangeLocation String
    | FilterWorlds String
    | GoBack
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnLocationChange Location

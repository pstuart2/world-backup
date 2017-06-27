module Msgs exposing (..)

import Http
import Material
import Models exposing (..)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = Mdl (Material.Msg Msg)
    | DoNothing
    | DeleteBackup FolderId WorldId BackupId
    | RestoreBackup FolderId WorldId BackupId
    | OnBackupDeleted FolderId WorldId BackupId (Result Http.Error ())
    | OnBackupRestored FolderId WorldId BackupId (Result Http.Error ())
    | ChangeLocation String
    | GoBack
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnLocationChange Location

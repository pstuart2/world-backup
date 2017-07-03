module Folders.Models exposing (..)

import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)


type alias FolderModel =
    { createBackupId : Maybe WorldId
    , deleteWorldId : Maybe WorldId
    , backupName : String
    , worldFilter : String
    , folders : WebData (List Folder)
    }


initialModel : FolderModel
initialModel =
    { createBackupId = Nothing
    , deleteWorldId = Nothing
    , backupName = ""
    , worldFilter = ""
    , folders = RemoteData.Loading
    }


type alias FolderId =
    String


type alias WorldId =
    String


type alias BackupId =
    String


type alias Folder =
    { id : FolderId
    , path : String
    , modifiedAt : DateTime
    , lastRun : DateTime
    , numberOfWorlds : Int
    , worlds : WebData (List World)
    }


type alias World =
    { id : WorldId
    , name : String
    , backups : List Backup
    }


type alias Backup =
    { id : String
    , name : String
    , createdAt : DateTime
    }


type alias BackupRequest =
    { name : String }

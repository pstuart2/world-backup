module Models exposing (..)

import Material
import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)


type Route
    = FoldersRoute
    | FolderRoute FolderId
    | NotFoundRoute


type alias Flags =
    { apiUrl : String
    }


type alias Model =
    { flags : Flags
    , mdl : Material.Model
    , folders : WebData (List Folder)
    , route : Route
    , folderView : FolderView
    }


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
    , mdl = Material.model
    , folders = RemoteData.Loading
    , route = route
    , folderView = FolderView Nothing Nothing "" ""
    }


type alias FolderView =
    { createBackupId : Maybe WorldId
    , deleteWorldId : Maybe WorldId
    , backupName : String
    , worldFilter : String
    }


type alias FolderId =
    String


type alias WorldId =
    String


type alias BackupId =
    String


type alias IconClass =
    String


type alias Message =
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

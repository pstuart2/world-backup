module Models exposing (..)

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
    , folders : WebData (List Folder)
    , route : Route
    }


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
    , folders = RemoteData.Loading
    , route = route
    }


type alias FolderId =
    String


type alias Folder =
    { id : FolderId
    , path : String
    , modifiedAt : DateTime
    , lastRun : DateTime
    , numberOfWorlds : Int
    }

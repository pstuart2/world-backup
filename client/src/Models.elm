module Models exposing (..)

import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)


type Route
    = FoldersRoute
    | FolderRoute FolderId
    | NotFoundRoute


type alias Options =
    { baseUrl : String
    }


type alias Model =
    { options : Options
    , folders : WebData (List Folder)
    , route : Route
    }


initialModel : Options -> Route -> Model
initialModel flags route =
    { options = flags
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
    }

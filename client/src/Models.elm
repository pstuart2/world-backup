module Models exposing (..)

import Folders.Models as Folders exposing (BackupId, Folder, FolderId, WorldId)
import Material


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
    , route : Route
    , folders : Folders.FolderModel
    }


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
    , mdl = Material.model
    , route = route
    , folders = Folders.initialModel
    }


type alias IconClass =
    String


type alias Message =
    String

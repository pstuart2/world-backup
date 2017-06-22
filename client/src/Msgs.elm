module Msgs exposing (..)

import Models exposing (Folder, FolderId, World)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = ChangeLocation String
    | GoBack
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnLocationChange Location

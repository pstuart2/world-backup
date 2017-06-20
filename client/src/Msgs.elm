module Msgs exposing (..)

import Models exposing (Folder, World)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = ChangeLocation String
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds (WebData (List World))
    | OnLocationChange Location

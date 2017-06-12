module Msgs exposing (..)

import Models exposing (Folder)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = ChangeLocation String
    | OnFetchFolders (WebData (List Folder))
    | OnLocationChange Location

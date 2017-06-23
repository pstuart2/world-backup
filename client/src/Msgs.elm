module Msgs exposing (..)

import Material
import Models exposing (Folder, FolderId, World)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = Mdl (Material.Msg Msg)
    | DoNothing
    | ChangeLocation String
    | GoBack
    | OnFetchFolders (WebData (List Folder))
    | OnFetchWorlds FolderId (WebData (List World))
    | OnLocationChange Location

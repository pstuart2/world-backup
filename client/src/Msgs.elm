module Msgs exposing (..)

import Folders.Msgs as FoldersMsgs
import Material
import Navigation exposing (Location)


type Msg
    = Mdl (Material.Msg Msg)
    | FolderMsg FoldersMsgs.Msg
    | DoNothing
    | ChangeLocation String
    | GoBack
    | OnLocationChange Location

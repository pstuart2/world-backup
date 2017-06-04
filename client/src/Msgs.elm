module Msgs exposing (..)

import Models exposing (Player)
import Navigation exposing (Location)
import RemoteData exposing (WebData)


type Msg
    = ChangeLocation String
    | OnFetchPlayers (WebData (List Player))
    | OnLocationChange Location

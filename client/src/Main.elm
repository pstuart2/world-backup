module Main exposing (..)

import Commands exposing (fetchFolders)
import Models exposing (Model, Options, initialModel)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Routing
import Update exposing (update)
import View exposing (view)


init : Options -> Location -> ( Model, Cmd Msg )
init options location =
    let
        currentRoute =
            Routing.parseLocation location
    in
    ( initialModel options currentRoute, fetchFolders )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- MAIN


main : Program Options Model Msg
main =
    Navigation.programWithFlags Msgs.OnLocationChange
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }

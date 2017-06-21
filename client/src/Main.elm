module Main exposing (..)

import Commands exposing (fetchFolders)
import Models exposing (Flags, Model, Route, initialModel)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Routing
import Update exposing (update)
import View exposing (view)


initialCommands : String -> Route -> Cmd Msg
initialCommands apiUrl currentRoute =
    let
        folderCommand =
            Routing.getLocationCommand apiUrl Models.FoldersRoute

        routeCommand =
            if currentRoute == Models.FoldersRoute then
                folderCommand
            else
                Cmd.batch [ Routing.getLocationCommand apiUrl currentRoute, folderCommand ]
    in
    routeCommand


init : Flags -> Location -> ( Model, Cmd Msg )
init flags location =
    let
        currentRoute =
            Routing.parseLocation location

        cmds =
            initialCommands flags.apiUrl currentRoute
    in
    ( initialModel flags currentRoute, cmds )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- MAIN


main : Program Flags Model Msg
main =
    Navigation.programWithFlags Msgs.OnLocationChange
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }

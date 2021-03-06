module Main exposing (..)

import Material
import Models exposing (Flags, Model, Route, initialModel)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Routing
import Time exposing (Time, second)
import Update exposing (update)
import View exposing (view)


initialCommands : String -> Route -> Cmd Msg
initialCommands apiUrl currentRoute =
    let
        folderCommand =
            Routing.getLocationCommand apiUrl Models.FoldersRoute

        routeCommand =
            Routing.getLocationCommand apiUrl currentRoute
    in
    if currentRoute == Models.FoldersRoute then
        folderCommand
    else
        Cmd.batch [ Material.init Msgs.Mdl, routeCommand, folderCommand ]


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
    Sub.batch
        [ Material.subscriptions Msgs.Mdl model
        , Time.every (20 * second) Msgs.Poll
        ]



-- MAIN


main : Program Flags Model Msg
main =
    Navigation.programWithFlags Msgs.OnLocationChange
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }

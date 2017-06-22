module Routing exposing (..)

import Commands exposing (fetchFolderWorlds, fetchFolders)
import Html exposing (Attribute)
import Html.Events exposing (onWithOptions)
import Json.Decode as Decode
import Models exposing (FolderId, Model, Route(..))
import Msgs exposing (Msg)
import Navigation exposing (Location)
import UrlParser exposing (..)


matchers : Parser (Route -> a) a
matchers =
    oneOf
        [ map FoldersRoute top
        , map FolderRoute (s "folders" </> string)
        , map FoldersRoute (s "folders")
        ]


onLinkClick : msg -> Attribute msg
onLinkClick message =
    let
        options =
            { stopPropagation = False
            , preventDefault = True
            }
    in
    onWithOptions "click" options (Decode.succeed message)


parseLocation : Location -> Route
parseLocation location =
    case parsePath matchers location of
        Just route ->
            route

        Nothing ->
            NotFoundRoute


getLocationCommand : String -> Route -> Cmd Msg
getLocationCommand apiUrl route =
    case route of
        Models.FoldersRoute ->
            fetchFolders apiUrl

        Models.FolderRoute id ->
            fetchFolderWorlds apiUrl id

        _ ->
            Cmd.none


homePath : String
homePath =
    "/"


foldersPath : String
foldersPath =
    "/folders"


folderPath : FolderId -> String
folderPath id =
    "/folders/" ++ id

module Folders.List exposing (view)

import Folders.Models exposing (Folder)
import Folders.Msgs as FolderMsgs
import Html exposing (..)
import Html.Attributes exposing (class, href)
import Material.Button as Button
import Material.Grid as Grid exposing (Device(..), cell, grid, size)
import Material.Options as Options
import Models exposing (Model)
import Msgs exposing (Msg)
import Numeral exposing (format)
import RemoteData exposing (WebData)
import Routing exposing (folderPath, onLinkClick)
import Time.DateTime as DateTime exposing (DateTime)


view : (FolderMsgs.Msg -> Msg) -> Model -> Html Msg
view pMsg model =
    maybeList pMsg model


maybeList : (FolderMsgs.Msg -> Msg) -> Model -> Html Msg
maybeList pMsg model =
    case model.folders.folders of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success folders ->
            list pMsg model folders

        RemoteData.Failure error ->
            text (toString error)


list : (FolderMsgs.Msg -> Msg) -> Model -> List Folder -> Html Msg
list pMsg model folders =
    div [ class "grid-outer" ]
        [ headerRow
        , folderBody pMsg model folders
        ]


headerRow : Html Msg
headerRow =
    grid [ Options.cs "headers" ]
        [ cell [ size All 2 ] [ text "" ]
        , cell [ size All 2 ] [ text "Last Run" ]
        , cell [ size All 6 ] [ text "Path" ]
        , cell [ size All 2 ] [ text "# of Worlds" ]
        ]


folderBody : (FolderMsgs.Msg -> Msg) -> Model -> List Folder -> Html Msg
folderBody pMsg model folders =
    div [] (List.map (folderRow pMsg model) folders)


folderRow : (FolderMsgs.Msg -> Msg) -> Model -> Folder -> Html Msg
folderRow pMsg model folder =
    grid []
        [ cell [ size All 2 ] [ viewFolderButton model folder ]
        , cell [ size All 2 ] [ text (DateTime.toISO8601 folder.lastRun) ]
        , cell [ size All 6 ] [ text folder.path ]
        , cell [ size All 2 ] [ text (format "0,0" (toFloat folder.numberOfWorlds)) ]
        ]


viewFolderButton : Model -> Folder -> Html.Html Msg
viewFolderButton model folder =
    Button.render Msgs.Mdl
        [ 0 ]
        model.mdl
        [ Button.raised
        , Button.colored
        , Options.onClick (Msgs.ChangeLocation (folderPath folder.id))
        ]
        [ text "View" ]

module Folders.List exposing (view)

import Folders.Models exposing (Folder)
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


view : Model -> Html Msg
view model =
    maybeList model


maybeList : Model -> Html Msg
maybeList model =
    case model.folders.folders of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success folders ->
            list model folders

        RemoteData.Failure error ->
            text (toString error)


list : Model -> List Folder -> Html Msg
list model folders =
    case folders of
        [] ->
            h4 [] [ text "No results..." ]

        _ ->
            populatedList model folders


populatedList : Model -> List Folder -> Html Msg
populatedList model folders =
    div [ class "grid-outer" ]
        [ headerRow
        , folderBody model folders
        ]


headerRow : Html Msg
headerRow =
    grid [ Options.cs "headers" ]
        [ cell [ size All 2 ] [ text "" ]
        , cell [ size All 2 ] [ text "Last Run" ]
        , cell [ size All 6 ] [ text "Path" ]
        , cell [ size All 2 ] [ text "# of Worlds" ]
        ]


folderBody : Model -> List Folder -> Html Msg
folderBody model folders =
    div [] (List.map (folderRow model) folders)


folderRow : Model -> Folder -> Html Msg
folderRow model folder =
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

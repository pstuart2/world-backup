module Folders.List exposing (..)

import Html exposing (..)
import Html.Attributes exposing (class, href)
import Models exposing (Folder)
import Msgs exposing (Msg)
import Numeral exposing (format)
import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)


view : WebData (List Folder) -> Html Msg
view response =
    div []
        [ maybeList response
        ]


maybeList : WebData (List Folder) -> Html Msg
maybeList response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success folders ->
            list folders

        RemoteData.Failure error ->
            text (toString error)


list : List Folder -> Html Msg
list folders =
    div [ class "grid-outer" ]
        [ headerRow
        , folderBody folders
        ]


headerRow : Html Msg
headerRow =
    div [ class "mdl-grid headers" ]
        [ div [ class "mdl-cell mdl-cell--2-col" ] [ text "Id" ]
        , div [ class "mdl-cell mdl-cell--2-col" ] [ text "Last Run" ]
        , div [ class "mdl-cell mdl-cell--6-col" ] [ text "Path" ]
        , div [ class "mdl-cell mdl-cell--2-col" ] [ text "Number of Worlds" ]
        ]


folderBody : List Folder -> Html Msg
folderBody folders =
    div [ class "grid-body" ] (List.map folderRow folders)


folderRow : Folder -> Html Msg
folderRow folder =
    div [ class "mdl-grid" ]
        [ div [ class "mdl-cell mdl-cell--2-col" ] [ text folder.id ]
        , div [ class "mdl-cell mdl-cell--2-col" ] [ text (DateTime.toISO8601 folder.lastRun) ]
        , div [ class "mdl-cell mdl-cell--6-col" ] [ text folder.path ]
        , div [ class "mdl-cell mdl-cell--2-col" ] [ text (format "0,0" (toFloat folder.numberOfWorlds)) ]
        ]

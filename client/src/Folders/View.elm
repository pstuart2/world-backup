module Folders.View exposing (view)

import Html exposing (..)
import Html.Attributes exposing (class, href, value)
import Models exposing (Backup, Folder, World)
import Msgs exposing (Msg)
import RemoteData
import Time.DateTime as DateTime exposing (DateTime)


view : Folder -> Html Msg
view folder =
    div []
        [ maybeList folder.worlds
        ]


maybeList : RemoteData.WebData (List World) -> Html Msg
maybeList response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success worlds ->
            list worlds

        RemoteData.Failure error ->
            text (toString error)


list : List World -> Html Msg
list worlds =
    div [ class "grid-outer" ]
        (List.map worldSection worlds)


worldSection : World -> Html Msg
worldSection world =
    div []
        [ h2 [] [ text world.name ]
        , worldBackups world.backups
        ]


worldBackups : List Backup -> Html Msg
worldBackups backups =
    div [ class "grid-body" ] (List.map worldBackup backups)


worldBackup : Backup -> Html Msg
worldBackup backup =
    div [ class "mdl-grid" ]
        [ div [ class "mdl-cell mdl-cell--2-col" ] [ text (DateTime.toISO8601 backup.createdAt) ]
        , div [ class "mdl-cell mdl-cell--10-col" ] [ text backup.name ]
        ]

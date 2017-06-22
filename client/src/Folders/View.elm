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
        , backupsTable world.backups
        ]


backupsTable : List Backup -> Html Msg
backupsTable backups =
    table [ class "mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp" ]
        [ thead []
            [ tr []
                [ th [ class "mdl-data-table__cell--non-numeric" ] [ text "Name" ]
                , th [ class "mdl-data-table__cell--non-numeric" ] [ text "Created At" ]
                ]
            ]
        , tbody [] (List.map backupRow backups)
        ]


backupRow : Backup -> Html Msg
backupRow backup =
    tr []
        [ td [ class "mdl-data-table__cell--non-numeric" ] [ text backup.name ]
        , td [ class "mdl-data-table__cell--non-numeric" ] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]

module Folders.View exposing (view)

import Html exposing (Html, div, h2, text)
import Html.Attributes exposing (class, href, value)
import Material.Options as Options
import Material.Table as Table exposing (table, tbody, td, th, thead, tr)
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
    table [ Options.css "width" "100%" ]
        [ thead []
            [ tr []
                [ th [ Table.numeric ] [ text "Name" ]
                , th [] [ text "Created At" ]
                ]
            ]
        , tbody [] (List.map backupRow backups)
        ]


backupRow : Backup -> Html Msg
backupRow backup =
    tr []
        [ td [ Table.numeric ] [ text backup.name ]
        , td [] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]

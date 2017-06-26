module Folders.View exposing (view)

import Html exposing (Html, div, h2, i, text)
import Html.Attributes exposing (class, href, value)
import Material.Button as Button
import Material.Color as Color
import Material.Options as Options
import Material.Table as Table exposing (table, tbody, td, th, thead, tr)
import Models exposing (..)
import Msgs exposing (Msg)
import RemoteData
import Time.DateTime as DateTime exposing (DateTime)


view : Model -> Folder -> Html Msg
view model folder =
    div []
        [ maybeList model folder.id folder.worlds
        ]


maybeList : Model -> FolderId -> RemoteData.WebData (List World) -> Html Msg
maybeList model folderId response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success worlds ->
            list model folderId worlds

        RemoteData.Failure error ->
            text (toString error)


list : Model -> FolderId -> List World -> Html Msg
list model folderId worlds =
    div [ class "grid-outer" ]
        (List.map (worldSection model folderId) worlds)


worldSection : Model -> FolderId -> World -> Html Msg
worldSection model folderId world =
    div []
        [ h2 [] [ text world.name ]
        , backupsTable model folderId world.id world.backups
        ]


backupsTable : Model -> FolderId -> WorldId -> List Backup -> Html Msg
backupsTable model folderId worldId backups =
    table [ Options.css "width" "100%" ]
        [ thead []
            [ tr []
                [ th [ Table.numeric ] [ text "Actions" ]
                , th [ Table.numeric ] [ text "Name" ]
                , th [] [ text "Created At" ]
                ]
            ]
        , tbody [] (List.map (backupRow model folderId worldId) backups)
        ]


backupRow : Model -> FolderId -> WorldId -> Backup -> Html Msg
backupRow model folderId worldId backup =
    tr []
        [ td [ Table.numeric ]
            [ iconButton model "fa fa-remove" (Color.color Color.Red Color.S900) (Msgs.DeleteBackup folderId worldId backup.id)
            , iconButton model "fa fa-check" (Color.color Color.Green Color.S900) Msgs.DoNothing
            ]
        , td [ Table.numeric ] [ text backup.name ]
        , td [] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]


iconButton : Model -> IconClass -> Color.Color -> Msg -> Html.Html Msg
iconButton model icon color clickMsg =
    Button.render Msgs.Mdl
        [ 0 ]
        model.mdl
        [ Button.icon
        , Color.text color
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] [] ]

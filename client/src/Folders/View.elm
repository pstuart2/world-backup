module Folders.View exposing (view)

import Folders.Buttons exposing (..)
import Folders.Confirms exposing (..)
import Html exposing (Html, div, h2, h4, i, text)
import Html.Attributes exposing (class, href, style, value)
import Material.Color as Color
import Material.Grid as Grid exposing (Align(..), Device(..), align, cell, grid, size)
import Material.Options as Options
import Material.Table as Table exposing (table, tbody, td, th, thead, tr)
import Material.Textfield as Textfield
import Models exposing (..)
import Msgs exposing (Msg)
import RemoteData
import Time.DateTime as DateTime exposing (DateTime)


view : Model -> Folder -> Html Msg
view model folder =
    div [ class "folder-worlds" ]
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
    let
        viewWorld id world =
            worldSection id model folderId world

        worldFilter world =
            String.contains model.folderView.worldFilter world.name

        filteredWorlds =
            List.filter worldFilter worlds
    in
    grid [ Grid.noSpacing ]
        [ cell [ size All 12 ] [ filter model ]
        , cell [ size All 12 ]
            (List.indexedMap viewWorld filteredWorlds)
        ]


filter : Model -> Html Msg
filter model =
    grid []
        [ cell [ size All 12 ]
            [ searchField model Msgs.FilterWorlds
            , cancelIconButton [ 0 ] model Msgs.ClearWorldsFilter
            ]
        ]


searchField : Model -> (String -> Msg) -> Html.Html Msg
searchField model msg =
    Textfield.render Msgs.Mdl
        [ 7 ]
        model.mdl
        [ Textfield.label "Filter worlds"
        , Textfield.floatingLabel
        , Textfield.value model.folderView.worldFilter
        , Options.css "width" "calc(100%  - 32px)"
        , Options.onInput msg
        ]
        []


worldSection : Int -> Model -> FolderId -> World -> Html Msg
worldSection iWorld model folderId world =
    let
        confirmContent =
            if model.folderView.createBackupId == Just world.id then
                backupConfirm [ iWorld ] model folderId world.id
            else if model.folderView.deleteWorldId == Just world.id then
                deleteConfirm [ iWorld ] model folderId world.id
            else
                text ""

        reverseCreatedAt a b =
            DateTime.compare b.createdAt a.createdAt

        sortedBackups =
            List.sortWith reverseCreatedAt world.backups
    in
    grid [ Options.cs "world" ]
        [ cell [ size Desktop 9, size Tablet 8, size Phone 4, Options.cs "world-title", align Middle ]
            [ h2 []
                [ i [ class "fa fa-globe" ] []
                , text world.name
                ]
            ]
        , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "world-buttons" ]
            [ backupButton [ iWorld ] model (Msgs.StartWorldBackup world.id)
            , deleteButton [ iWorld ] model (Msgs.StartWorldDelete world.id)
            ]
        , cell [ size All 12 ]
            [ confirmContent
            ]
        , cell [ size All 12 ]
            [ backupsTable iWorld model folderId world.id sortedBackups
            ]
        ]


backupsTable : Int -> Model -> FolderId -> WorldId -> List Backup -> Html Msg
backupsTable iWorld model folderId worldId backups =
    let
        viewBackup id backup =
            backupRow [ iWorld, id ] model folderId worldId backup
    in
    table [ Options.css "width" "100%" ]
        [ thead []
            [ tr []
                [ th [ Table.numeric ] [ text "Actions" ]
                , th [ Table.numeric ] [ text "Name" ]
                , th [ Table.descending ] [ text "Created At" ]
                ]
            ]
        , tbody [] (List.indexedMap viewBackup backups)
        ]


backupRow : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
backupRow idx model folderId worldId backup =
    tr []
        [ td [ Table.numeric, Options.cs "button-group" ]
            [ iconButton idx model "fa fa-trash-o" (Color.color Color.Red Color.S900) (Msgs.DeleteBackup folderId worldId backup.id)
            , iconButton idx model "fa fa-recycle" (Color.color Color.Green Color.S900) (Msgs.RestoreBackup folderId worldId backup.id)
            ]
        , td [ Table.numeric ] [ text backup.name ]
        , td [] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]

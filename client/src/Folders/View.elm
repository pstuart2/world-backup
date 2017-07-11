module Folders.View exposing (view)

import Folders.Buttons exposing (..)
import Folders.Confirms exposing (..)
import Folders.Models exposing (Backup, Folder, FolderId, World, WorldId)
import Folders.Msgs as FolderMsgs
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


view : (FolderMsgs.Msg -> Msg) -> Model -> Folder -> Html Msg
view pMsg model folder =
    div [ class "folder-worlds" ]
        [ maybeList pMsg model folder.id folder.worlds
        ]


maybeList : (FolderMsgs.Msg -> Msg) -> Model -> FolderId -> RemoteData.WebData (List World) -> Html Msg
maybeList pMsg model folderId response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success worlds ->
            list pMsg model folderId worlds

        RemoteData.Failure error ->
            text (toString error)


list : (FolderMsgs.Msg -> Msg) -> Model -> FolderId -> List World -> Html Msg
list pMsg model folderId worlds =
    let
        viewWorld id world =
            worldSection pMsg id model folderId world

        worldFilter world =
            String.contains model.folders.worldFilter world.name

        filteredWorlds =
            List.filter worldFilter worlds
    in
    grid [ Grid.noSpacing ]
        [ cell [ size All 12 ] [ filter pMsg model ]
        , cell [ size All 12 ]
            (List.indexedMap viewWorld filteredWorlds)
        ]


filter : (FolderMsgs.Msg -> Msg) -> Model -> Html Msg
filter pMsg model =
    grid []
        [ cell [ size All 12 ]
            [ searchField model (\inp -> pMsg (FolderMsgs.FilterWorlds inp))
            , cancelIconButton [ 0 ] model (pMsg FolderMsgs.ClearWorldsFilter)
            ]
        ]


searchField : Model -> (String -> Msg) -> Html.Html Msg
searchField model msg =
    Textfield.render Msgs.Mdl
        [ 7 ]
        model.mdl
        [ Textfield.label "Filter worlds"
        , Textfield.floatingLabel
        , Textfield.value model.folders.worldFilter
        , Options.css "width" "calc(100%  - 32px)"
        , Options.onInput msg
        ]
        []


worldSection : (FolderMsgs.Msg -> Msg) -> Int -> Model -> FolderId -> World -> Html Msg
worldSection pMsg iWorld model folderId world =
    let
        confirmContent =
            if model.folders.createBackupId == Just world.id then
                backupConfirm pMsg [ iWorld ] model folderId world.id
            else if model.folders.deleteWorldId == Just world.id then
                deleteConfirm pMsg [ iWorld ] model folderId world.id
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
            [ backupButton [ iWorld ] model (pMsg (FolderMsgs.StartWorldBackup world.id))
            , deleteButton [ iWorld ] model (pMsg (FolderMsgs.StartWorldDelete world.id))
            ]
        , cell [ size All 12 ]
            [ confirmContent
            ]
        , cell [ size All 12 ]
            [ backupsTable pMsg iWorld model folderId world.id sortedBackups
            ]
        ]


backupsTable : (FolderMsgs.Msg -> Msg) -> Int -> Model -> FolderId -> WorldId -> List Backup -> Html Msg
backupsTable pMsg iWorld model folderId worldId backups =
    let
        viewBackup id backup =
            backupRow pMsg [ iWorld, id ] model folderId worldId backup
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


backupRow : (FolderMsgs.Msg -> Msg) -> List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
backupRow pMsg idx model folderId worldId backup =
    tr []
        [ td [ Table.numeric, Options.cs "button-group" ]
            [ iconButton idx model "fa fa-trash-o" (Color.color Color.Red Color.S900) (pMsg (FolderMsgs.DeleteBackup folderId worldId backup.id))
            , iconButton idx model "fa fa-recycle" (Color.color Color.Green Color.S900) (pMsg (FolderMsgs.RestoreBackup folderId worldId backup.id))
            ]
        , td [ Table.numeric ] [ text backup.name ]
        , td [] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]

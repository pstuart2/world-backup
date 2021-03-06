module Folders.View exposing (view)

import Folders.Buttons exposing (..)
import Folders.Confirms exposing (..)
import Folders.Models exposing (Backup, Folder, FolderId, World, WorldId)
import Html exposing (Html, div, h2, h4, i, text)
import Html.Attributes exposing (class, href, style, value)
import Material.Grid as Grid exposing (Align(..), Device(..), align, cell, grid, size)
import Material.Options as Options
import Material.Table as Table exposing (table, tbody, td, th, thead, tr)
import Material.Textfield as Textfield
import Models exposing (..)
import Msgs exposing (..)
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
            String.contains (String.toLower model.folders.worldFilter) (String.toLower world.name)

        filteredWorlds =
            List.filter worldFilter worlds

        lowerCompare a b =
            compare (String.toLower a.name) (String.toLower b.name)

        sortedWorlds =
            List.sortWith lowerCompare filteredWorlds

        content =
            case filteredWorlds of
                [] ->
                    [ emptyList ]

                _ ->
                    List.indexedMap viewWorld sortedWorlds
    in
    grid [ Grid.noSpacing ]
        [ cell [ size All 12 ] [ filter model ]
        , cell [ size All 12 ]
            content
        ]


emptyList : Html Msg
emptyList =
    h4 [] [ text "No results..." ]


filter : Model -> Html Msg
filter model =
    grid []
        [ cell [ size All 12 ]
            [ searchField model (FilterWorlds >> FolderMsg)
            , cancelIconButton [ 0 ] model (ClearWorldsFilter |> FolderMsg)
            ]
        ]


searchField : Model -> (String -> Msg) -> Html.Html Msg
searchField model msg =
    Textfield.render Mdl
        [ 7 ]
        model.mdl
        [ Textfield.label "Filter worlds"
        , Textfield.floatingLabel
        , Textfield.value model.folders.worldFilter
        , Options.css "width" "calc(100%  - 32px)"
        , Options.onInput msg
        ]
        []


worldSection : Int -> Model -> FolderId -> World -> Html Msg
worldSection iWorld model folderId world =
    let
        confirmContent =
            if model.folders.createBackupId == Just world.id then
                backupConfirm [ iWorld ] model folderId world.id
            else if model.folders.deleteWorldId == Just world.id then
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
            [ backupButton [ iWorld ] model (StartWorldBackup world.id |> FolderMsg)
            , deleteButton [ iWorld ] model (StartWorldDelete world.id |> FolderMsg)
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
            if model.folders.restoreBackupId == Just backup.id then
                backupRestoreConfirmRow [ iWorld, id ] model folderId worldId backup
            else if model.folders.deleteBackupId == Just backup.id then
                backupDeleteConfirmRow [ iWorld, id ] model folderId worldId backup
            else
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
            [ trashButton idx model (DeleteBackupConfirm backup.id backup.name |> FolderMsg)
            , recycleButton idx model (RestoreBackupConfirm backup.id backup.name |> FolderMsg)
            ]
        , td [ Table.numeric ] [ text backup.name ]
        , td [] [ text (DateTime.toISO8601 backup.createdAt) ]
        ]


backupRestoreConfirmRow : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
backupRestoreConfirmRow idx model folderId worldId backup =
    tr []
        [ td
            [ Options.cs "row-confirm-msg"
            , Options.attribute (Html.Attributes.colspan 3)
            ]
            [ restoreBackupConfirm idx model folderId worldId backup
            ]
        ]


backupDeleteConfirmRow : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
backupDeleteConfirmRow idx model folderId worldId backup =
    tr []
        [ td
            [ Options.cs "row-confirm-msg"
            , Options.attribute (Html.Attributes.colspan 3)
            ]
            [ deleteBackupConfirm idx model folderId worldId backup
            ]
        ]

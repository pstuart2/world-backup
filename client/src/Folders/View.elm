module Folders.View exposing (view)

import Html exposing (Html, div, h2, h4, i, text)
import Html.Attributes exposing (class, href, style, value)
import Material.Button as Button
import Material.Color as Color
import Material.Elevation as Elevation
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
            , iconButton [ 9 ] model "fa fa-times-circle" (Color.color Color.Grey Color.S400) Msgs.ClearWorldsFilter
            ]
        ]


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


iconButton : List Int -> Model -> IconClass -> Color.Color -> Msg -> Html.Html Msg
iconButton idx model icon color clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 0 ] idx)
        model.mdl
        [ Button.icon
        , Color.text color
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] [] ]


deleteButton : List Int -> Model -> Msg -> Html.Html Msg
deleteButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 1 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Color.text Color.white
        , Color.background (Color.color Color.Red Color.S300)
        , Options.onClick clickMsg
        ]
        [ i [ class "fa fa-trash-o" ] []
        , text "Delete"
        ]


backupButton : List Int -> Model -> Msg -> Html.Html Msg
backupButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 2 ] idx)
        model.mdl
        [ Button.ripple
        , Button.colored
        , Button.raised
        , Options.onClick clickMsg
        ]
        [ i [ class "fa fa-clone" ] []
        , text "Backup"
        ]


searchField : Model -> (String -> Msg) -> Html.Html Msg
searchField model msg =
    Textfield.render Msgs.Mdl
        [ 7 ]
        model.mdl
        [ Textfield.label "Filter"
        , Textfield.floatingLabel
        , Textfield.value model.folderView.worldFilter
        , Options.css "width" "calc(100%  - 32px)"
        , Options.onInput msg
        ]
        []


backupNameField : List Int -> Model -> Html.Html Msg
backupNameField idx model =
    Textfield.render Msgs.Mdl
        (List.append [ 8 ] idx)
        model.mdl
        [ Textfield.label "Backup name"
        , Textfield.floatingLabel
        , Textfield.value model.folderView.backupName
        , Options.css "width" "100%"
        , Options.onInput Msgs.UpdateBackupName
        ]
        []


inlineInfo : Html.Html Msg -> Html.Html Msg
inlineInfo content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Blue Color.S50)
        , Options.cs "confirm"
        ]
        [ content ]


backupConfirm : List Int -> Model -> FolderId -> WorldId -> Html.Html Msg
backupConfirm idx model folderId worldId =
    inlineInfo
        (grid []
            [ cell [ size All 12 ] [ h4 [] [ text "Enter a name for your backup and confirm." ] ]
            , cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ backupNameField idx model ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model Msgs.CancelWorldBackup
                , confirmButton idx "Backup" "fa fa-clone" model (Msgs.BackupWorld folderId worldId model.folderView.backupName)
                ]
            ]
        )


confirmButton : List Int -> Message -> IconClass -> Model -> Msg -> Html.Html Msg
confirmButton idx buttonText icon model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 10 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Button.colored
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] []
        , text buttonText
        ]


destructiveConfirmButton : List Int -> Message -> IconClass -> Model -> Msg -> Html.Html Msg
destructiveConfirmButton idx buttonText icon model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 10 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Color.text Color.white
        , Color.background (Color.color Color.Red Color.S900)
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] []
        , text buttonText
        ]


cancelButton : List Int -> Model -> Msg -> Html.Html Msg
cancelButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 10 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Color.text Color.white
        , Color.background (Color.color Color.Grey Color.S500)
        , Options.onClick clickMsg
        ]
        [ i [ class "fa fa-times" ] []
        , text "Cancel"
        ]


inlineWarning : Html.Html Msg -> Html.Html Msg
inlineWarning content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Red Color.S50)
        , Color.text (Color.color Color.Red Color.S900)
        , Options.cs "confirm"
        ]
        [ content ]


deleteConfirm : List Int -> Model -> FolderId -> WorldId -> Html.Html Msg
deleteConfirm idx model folderId worldId =
    inlineWarning
        (grid
            []
            [ cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ h4 [] [ text "Are you sure you want to delete this world?" ] ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model Msgs.CancelDeleteWorld
                , destructiveConfirmButton idx "Delete" "fa fa-trash-o" model (Msgs.DeleteWorld folderId worldId)
                ]
            ]
        )

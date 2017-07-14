module Folders.Confirms exposing (backupConfirm, deleteBackupConfirm, deleteConfirm, restoreBackupConfirm)

import Folders.Buttons exposing (..)
import Folders.Models exposing (Backup, BackupId, FolderId, WorldId)
import Html exposing (Html, div, h2, h4, i, span, text)
import Html.Attributes exposing (class)
import Material.Color as Color
import Material.Elevation as Elevation
import Material.Grid as Grid exposing (Align(..), Device(..), align, cell, grid, size)
import Material.Options as Options
import Material.Textfield as Textfield
import Models exposing (..)
import Msgs exposing (..)


deleteConfirm : List Int -> Model -> FolderId -> WorldId -> Html Msg
deleteConfirm idx model folderId worldId =
    inlineWarning
        (grid
            []
            [ cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ h4 [] [ text "Are you sure you want to delete this world?" ] ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (CancelConfirm |> FolderMsg)
                , destructiveConfirmButton "fa fa-trash-o" "Delete" idx model (DeleteWorld folderId worldId |> FolderMsg)
                ]
            ]
        )


restoreBackupConfirm : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
restoreBackupConfirm idx model folderId worldId backup =
    inlineInfo
        (grid
            []
            [ cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ span []
                    [ text "Are you sure you want to restore backup "
                    , span [ class "confirm-target" ] [ text backup.name ]
                    , text "?"
                    ]
                ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (CancelConfirm |> FolderMsg)
                , primaryButton "fa fa-recycle" "Restore" idx model (RestoreBackup folderId worldId backup.id |> FolderMsg)
                ]
            ]
        )


deleteBackupConfirm : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
deleteBackupConfirm idx model folderId worldId backup =
    inlineWarning
        (grid
            []
            [ cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ span []
                    [ text "Are you sure you want to delete backup "
                    , span [ class "confirm-target" ] [ text backup.name ]
                    , text "?"
                    ]
                ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (CancelConfirm |> FolderMsg)
                , destructiveConfirmButton "fa fa-trash-o" "Delete" idx model (DeleteBackup folderId worldId backup.id |> FolderMsg)
                ]
            ]
        )


backupConfirm : List Int -> Model -> FolderId -> WorldId -> Html Msg
backupConfirm idx model folderId worldId =
    inlineInfo
        (grid []
            [ cell [ size All 12 ] [ h4 [] [ text "Enter a name for your backup and confirm." ] ]
            , cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ backupNameField idx model ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (CancelConfirm |> FolderMsg)
                , primaryButton "fa fa-clone" "Backup" idx model (BackupWorld folderId worldId model.folders.backupName |> FolderMsg)
                ]
            ]
        )


backupNameField : List Int -> Model -> Html Msg
backupNameField idx model =
    Textfield.render Mdl
        (List.append [ 8 ] idx)
        model.mdl
        [ Textfield.label "Backup name"
        , Textfield.floatingLabel
        , Textfield.value model.folders.backupName
        , Options.css "width" "100%"
        , Options.onInput (UpdateBackupName >> FolderMsg)
        ]
        []


inlineInfo : Html Msg -> Html Msg
inlineInfo content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Blue Color.S50)
        , Options.cs "confirm"
        ]
        [ content ]


inlineWarning : Html Msg -> Html Msg
inlineWarning content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Red Color.S50)
        , Color.text (Color.color Color.Red Color.S900)
        , Options.cs "confirm"
        ]
        [ content ]

module Folders.Confirms exposing (backupConfirm, deleteConfirm)

import Folders.Buttons exposing (..)
import Folders.Models exposing (FolderId, WorldId)
import Folders.Msgs as FolderMsgs
import Html exposing (Html, div, h2, h4, i, text)
import Material.Color as Color
import Material.Elevation as Elevation
import Material.Grid as Grid exposing (Align(..), Device(..), align, cell, grid, size)
import Material.Options as Options
import Material.Textfield as Textfield
import Models exposing (..)
import Msgs exposing (Msg)


deleteConfirm : List Int -> Model -> FolderId -> WorldId -> Html.Html Msg
deleteConfirm idx model folderId worldId =
    inlineWarning
        (grid
            []
            [ cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ h4 [] [ text "Are you sure you want to delete this world?" ] ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (Msgs.FolderMsg FolderMsgs.CancelDeleteWorld)
                , destructiveConfirmButton idx "Delete" "fa fa-trash-o" model (Msgs.FolderMsg (FolderMsgs.DeleteWorld folderId worldId))
                ]
            ]
        )



-- TODO: Msgs.FolderMsg should be passed in from the top level


backupConfirm : List Int -> Model -> FolderId -> WorldId -> Html.Html Msg
backupConfirm idx model folderId worldId =
    inlineInfo
        (grid []
            [ cell [ size All 12 ] [ h4 [] [ text "Enter a name for your backup and confirm." ] ]
            , cell [ size Desktop 9, size Tablet 8, size Phone 4 ]
                [ backupNameField idx model ]
            , cell [ size Desktop 3, size Tablet 8, size Phone 4, align Middle, Options.cs "button-group" ]
                [ cancelButton idx model (Msgs.FolderMsg FolderMsgs.CancelWorldBackup)
                , confirmButton idx "Backup" "fa fa-clone" model (Msgs.FolderMsg (FolderMsgs.BackupWorld folderId worldId model.folders.backupName))
                ]
            ]
        )


backupNameField : List Int -> Model -> Html.Html Msg
backupNameField idx model =
    Textfield.render Msgs.Mdl
        (List.append [ 8 ] idx)
        model.mdl
        [ Textfield.label "Backup name"
        , Textfield.floatingLabel
        , Textfield.value model.folders.backupName
        , Options.css "width" "100%"
        , Options.onInput (\inp -> Msgs.FolderMsg (FolderMsgs.UpdateBackupName inp))
        ]
        []



-- TODO: These could be in a more generic place


inlineInfo : Html.Html Msg -> Html.Html Msg
inlineInfo content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Blue Color.S50)
        , Options.cs "confirm"
        ]
        [ content ]


inlineWarning : Html.Html Msg -> Html.Html Msg
inlineWarning content =
    Options.div
        [ Elevation.e2
        , Color.background (Color.color Color.Red Color.S50)
        , Color.text (Color.color Color.Red Color.S900)
        , Options.cs "confirm"
        ]
        [ content ]

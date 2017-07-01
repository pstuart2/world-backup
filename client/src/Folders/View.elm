module Folders.View exposing (view)

import Html exposing (Html, div, h2, i, text)
import Html.Attributes exposing (class, href, value)
import Material.Button as Button
import Material.Color as Color
import Material.Options as Options
import Material.Table as Table exposing (table, tbody, td, th, thead, tr)
import Material.Textfield as Textfield
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
    let
        viewWorld id world =
            worldSection id model folderId world

        worldFilter world =
            String.contains model.worldFilter world.name

        filteredWorlds =
            List.filter worldFilter worlds
    in
    div []
        [ div [ class "search-box" ] [ searchField model Msgs.FilterWorlds ]
        , div [ class "grid-outer" ]
            (List.indexedMap viewWorld filteredWorlds)
        ]


worldSection : Int -> Model -> FolderId -> World -> Html Msg
worldSection iWorld model folderId world =
    div []
        [ div [ class "world-buttons" ]
            [ backupButton [ iWorld ] model Msgs.DoNothing
            , deleteButton [ iWorld ] model (Msgs.DeleteWorld folderId world.id)
            ]
        , h2 [] [ text world.name ]
        , backupsTable iWorld model folderId world.id world.backups
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
                , th [] [ text "Created At" ]
                ]
            ]
        , tbody [] (List.indexedMap viewBackup backups)
        ]


backupRow : List Int -> Model -> FolderId -> WorldId -> Backup -> Html Msg
backupRow idx model folderId worldId backup =
    tr []
        [ td [ Table.numeric ]
            [ iconButton idx model "fa fa-remove" (Color.color Color.Red Color.S900) (Msgs.DeleteBackup folderId worldId backup.id)
            , iconButton idx model "fa fa-check" (Color.color Color.Green Color.S900) (Msgs.RestoreBackup folderId worldId backup.id)
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
        [ text "Delete" ]


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
        [ text "Backup" ]


searchField : Model -> (String -> Msg) -> Html.Html Msg
searchField model msg =
    Textfield.render Msgs.Mdl
        [ 7 ]
        model.mdl
        [ Textfield.label "Search"
        , Textfield.floatingLabel
        , Textfield.expandable "id-of-expandable-1"
        , Textfield.expandableIcon "search"
        , Options.onInput msg
        ]
        []

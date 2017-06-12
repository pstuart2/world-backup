module View exposing (..)

import Folders.Edit
import Folders.List
import Html exposing (Html, div, text)
import Models exposing (FolderId, Model)
import Msgs exposing (Msg)
import RemoteData


view : Model -> Html Msg
view model =
    div []
        [ page model ]


page : Model -> Html Msg
page model =
    case model.route of
        Models.FoldersRoute ->
            Folders.List.view model.folders

        Models.FolderRoute id ->
            folderEditPage model id

        Models.NotFoundRoute ->
            notFoundView


folderEditPage : Model -> FolderId -> Html Msg
folderEditPage model folderId =
    case model.folders of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading ..."

        RemoteData.Success folders ->
            let
                maybeFolder =
                    folders
                        |> List.filter (\folder -> folder.id == folderId)
                        |> List.head
            in
            case maybeFolder of
                Just folder ->
                    Folders.Edit.view folder

                Nothing ->
                    notFoundView

        RemoteData.Failure err ->
            text (toString err)


notFoundView : Html msg
notFoundView =
    div []
        [ text "Not Found" ]

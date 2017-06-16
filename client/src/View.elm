module View exposing (..)

import Folders.List
import Folders.View
import Html exposing (Html, div, header, main_, span, text)
import Html.Attributes exposing (class, href)
import Models exposing (FolderId, Model)
import Msgs exposing (Msg)
import RemoteData


view : Model -> Html Msg
view model =
    div [ class "mdl-layout mdl-js-layout mdl-layout--fixed-header" ]
        [ pageHeader model
        , pageContent model
        ]


pageHeader : Model -> Html Msg
pageHeader model =
    header [ class "mdl-layout__header" ]
        [ div [ class "mdl-layout__header-row" ]
            [ span [ class "mdl-layout-title" ] [ text "World Backup" ] ]
        ]


pageContent : Model -> Html Msg
pageContent model =
    main_ [ class "mdl-layout__content" ]
        [ div [ class "page-content" ]
            [ page model ]
        ]


page : Model -> Html Msg
page model =
    case model.route of
        Models.FoldersRoute ->
            Folders.List.view model.folders

        Models.FolderRoute id ->
            folderViewPage model id

        Models.NotFoundRoute ->
            notFoundView


folderViewPage : Model -> FolderId -> Html Msg
folderViewPage model folderId =
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
                    Folders.View.view folder

                Nothing ->
                    notFoundView

        RemoteData.Failure err ->
            text (toString err)


notFoundView : Html msg
notFoundView =
    div []
        [ text "Not Found" ]

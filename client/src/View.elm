module View exposing (..)

import Folders.List
import Folders.Models exposing (FolderId)
import Folders.View
import Html exposing (..)
import Html.Attributes exposing (class, href, style)
import Material.Button as Button
import Material.Layout as Layout
import Material.Options as Options exposing (Style)
import Models exposing (Model)
import Msgs exposing (Msg)
import RemoteData
import Routing exposing (homePath, onLinkClick)


view : Model -> Html Msg
view model =
    Layout.render Msgs.Mdl
        model.mdl
        [ Layout.fixedHeader
        ]
        { header = [ pageHeader model ]
        , drawer = []
        , tabs = ( [], [] )
        , main = [ page model ]
        }


pageHeader : Model -> Html Msg
pageHeader model =
    h3 [ style [ ( "padding-left", "10px" ) ] ]
        [ homeButton model
        , text "World Backup"
        ]


homeButton : Model -> Html.Html Msg
homeButton model =
    Button.render Msgs.Mdl
        [ 0 ]
        model.mdl
        [ Button.icon
        , Options.onClick (Msgs.ChangeLocation homePath)
        , Options.css "margin-right" "20px"
        ]
        [ i [ class "fa fa-home" ] [] ]


page : Model -> Html Msgs.Msg
page model =
    case model.route of
        Models.FoldersRoute ->
            Folders.List.view model

        Models.FolderRoute id ->
            folderViewPage model id

        Models.NotFoundRoute ->
            notFoundView


folderViewPage : Model -> FolderId -> Html Msg
folderViewPage model folderId =
    case model.folders.folders of
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
                    Folders.View.view model folder

                Nothing ->
                    notFoundView

        RemoteData.Failure err ->
            text (toString err)


notFoundView : Html msg
notFoundView =
    div []
        [ text "Not Found" ]

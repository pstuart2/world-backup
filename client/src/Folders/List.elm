module Folders.List exposing (..)

import Html exposing (..)
import Html.Attributes exposing (class, href)
import Models exposing (Folder)
import Msgs exposing (Msg)
import RemoteData exposing (WebData)
import Routing exposing (folderPath, onLinkClick)


view : WebData (List Folder) -> Html Msg
view response =
    div []
        [ nav
        , maybeList response
        ]


maybeList : WebData (List Folder) -> Html Msg
maybeList response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success folders ->
            list folders

        RemoteData.Failure error ->
            text (toString error)


nav : Html Msg
nav =
    div [ class "clearfix mb2 white bg-black" ]
        [ div [ class "left p2" ]
            [ text "Folders" ]
        ]


list : List Folder -> Html Msg
list folders =
    div [ class "p2" ]
        [ table []
            [ thead []
                [ tr []
                    [ th [] [ text "Id" ]
                    , th [] [ text "Path" ]
                    ]
                ]
            , tbody [] (List.map folderRow folders)
            ]
        ]


folderRow : Folder -> Html Msg
folderRow folder =
    tr []
        [ td [] [ text folder.id ]
        , td [] [ text folder.path ]
        ]


editBtn : Folder -> Html.Html Msg
editBtn folder =
    let
        path =
            folderPath folder.id
    in
    a [ class "btn regular", href path, onLinkClick (Msgs.ChangeLocation path) ]
        [ i [ class "fa fa-pencil mr1" ] []
        , text "Edit"
        ]

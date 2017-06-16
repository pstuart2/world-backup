module Folders.View exposing (..)

import Html exposing (..)
import Html.Attributes exposing (class, href, value)
import Models exposing (Folder)
import Msgs exposing (Msg)
import Routing exposing (foldersPath, onLinkClick)


view : Folder -> Html Msg
view model =
    div []
        [ form model
        ]


form : Folder -> Html Msg
form folder =
    div [ class "m3" ]
        [ h1 [] [ text folder.path ]
        , formLevel folder
        ]


formLevel : Folder -> Html Msg
formLevel folder =
    div [ class "clearfix py1" ]
        [ div [ class "col col-5" ] [ text "Level" ]
        , div [ class "col col-7" ]
            [ span [ class "h2 bold" ] [ text folder.id ]
            ]
        ]

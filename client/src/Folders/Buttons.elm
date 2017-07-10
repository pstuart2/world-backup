module Folders.Buttons exposing (..)

import Html exposing (Html, div, h2, h4, i, text)
import Html.Attributes exposing (class, href, style, value)
import Material.Button as Button
import Material.Color as Color
import Material.Options as Options
import Models exposing (..)
import Msgs exposing (Msg)


iconButton : List Int -> Model -> IconClass -> Color.Color -> Msg -> Html Msg
iconButton idx model icon color clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 0 ] idx)
        model.mdl
        [ Button.icon
        , Color.text color
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] [] ]


cancelIconButton : List Int -> Model -> Msg -> Html Msg
cancelIconButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 1 ] idx)
        model.mdl
        [ Button.icon
        , Color.text (Color.color Color.Grey Color.S400)
        , Options.onClick clickMsg
        ]
        [ i [ class "fa fa-times-circle" ] [] ]


deleteButton : List Int -> Model -> Msg -> Html Msg
deleteButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 2 ] idx)
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


backupButton : List Int -> Model -> Msg -> Html Msg
backupButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 3 ] idx)
        model.mdl
        [ Button.ripple
        , Button.colored
        , Button.raised
        , Options.onClick clickMsg
        ]
        [ i [ class "fa fa-clone" ] []
        , text "Backup"
        ]


destructiveConfirmButton : List Int -> Message -> IconClass -> Model -> Msg -> Html Msg
destructiveConfirmButton idx buttonText icon model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 4 ] idx)
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


cancelButton : List Int -> Model -> Msg -> Html Msg
cancelButton idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 5 ] idx)
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


confirmButton : List Int -> Message -> IconClass -> Model -> Msg -> Html Msg
confirmButton idx buttonText icon model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 6 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Button.colored
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] []
        , text buttonText
        ]

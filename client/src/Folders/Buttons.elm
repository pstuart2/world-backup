module Folders.Buttons exposing (..)

import Html exposing (Html, div, h2, h4, i, text)
import Html.Attributes exposing (class, href, style, value)
import Material.Button as Button
import Material.Color as Color exposing (Color)
import Material.Options as Options
import Models exposing (..)
import Msgs exposing (Msg)


-- Base Buttons


iconButton : IconClass -> Color -> List Int -> Model -> Msg -> Html Msg
iconButton icon color idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 0 ] idx)
        model.mdl
        [ Button.icon
        , Color.text color
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] [] ]


primaryButton : IconClass -> Message -> List Int -> Model -> Msg -> Html Msg
primaryButton icon buttonText idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 1 ] idx)
        model.mdl
        [ Button.ripple
        , Button.colored
        , Button.raised
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] []
        , text buttonText
        ]


whiteTextButton : IconClass -> Message -> Color -> List Int -> Model -> Msg -> Html Msg
whiteTextButton icon buttonText backgroundColor idx model clickMsg =
    Button.render Msgs.Mdl
        (List.append [ 2 ] idx)
        model.mdl
        [ Button.raised
        , Button.ripple
        , Color.text Color.white
        , Color.background backgroundColor
        , Options.onClick clickMsg
        ]
        [ i [ class icon ] []
        , text buttonText
        ]



-- Icon Buttons


cancelIconButton : List Int -> Model -> Msg -> Html Msg
cancelIconButton idx model clickMsg =
    iconButton "fa fa-times-circle" (Color.color Color.Grey Color.S400) (List.append [ 0 ] idx) model clickMsg


trashButton : List Int -> Model -> Msg -> Html Msg
trashButton idx model clickMsg =
    iconButton "fa fa-trash-o" (Color.color Color.Red Color.S900) (List.append [ 1 ] idx) model clickMsg


recycleButton : List Int -> Model -> Msg -> Html Msg
recycleButton idx model clickMsg =
    iconButton "fa fa-recycle" (Color.color Color.Green Color.S900) (List.append [ 2 ] idx) model clickMsg



-- Primary Buttons


backupButton : List Int -> Model -> Msg -> Html Msg
backupButton idx model clickMsg =
    primaryButton "fa fa-clone" "Backup" (List.append [ 0 ] idx) model clickMsg



-- White Text Buttons


deleteButton : List Int -> Model -> Msg -> Html Msg
deleteButton idx model clickMsg =
    whiteTextButton "fa fa-trash-o" "Delete" (Color.color Color.Red Color.S300) (List.append [ 0 ] idx) model clickMsg


destructiveConfirmButton : IconClass -> Message -> List Int -> Model -> Msg -> Html Msg
destructiveConfirmButton icon buttonText idx model clickMsg =
    whiteTextButton icon buttonText (Color.color Color.Red Color.S900) (List.append [ 1 ] idx) model clickMsg


cancelButton : List Int -> Model -> Msg -> Html Msg
cancelButton idx model clickMsg =
    whiteTextButton "fa fa-times" "Cancel" (Color.color Color.Grey Color.S500) (List.append [ 2 ] idx) model clickMsg

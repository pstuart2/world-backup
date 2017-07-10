module Update exposing (..)

import Folders.Update as FoldersUpdate
import Material
import Models exposing (..)
import Msgs exposing (Msg)
import Navigation exposing (back, newUrl)
import Routing exposing (getLocationCommand, parseLocation)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msgs.DoNothing ->
            ( model, Cmd.none )

        Msgs.Mdl msg_ ->
            Material.update Msgs.Mdl msg_ model

        Msgs.ChangeLocation path ->
            ( model, newUrl path )

        Msgs.GoBack ->
            ( model, back 1 )

        Msgs.OnLocationChange location ->
            let
                newRoute =
                    parseLocation location

                newCommand =
                    getLocationCommand model.flags.apiUrl newRoute
            in
            ( { model | route = newRoute }, newCommand )

        Msgs.FolderMsg msg_ ->
            FoldersUpdate.update Msgs.FolderMsg msg_ model

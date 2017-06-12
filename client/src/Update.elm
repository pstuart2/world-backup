module Update exposing (..)

import Models exposing (Folder, Model)
import Msgs exposing (Msg)
import Navigation exposing (newUrl)
import RemoteData
import Routing exposing (parseLocation)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msgs.ChangeLocation path ->
            ( model, newUrl path )

        Msgs.OnFetchFolders response ->
            ( { model | folders = response }, Cmd.none )

        Msgs.OnLocationChange location ->
            let
                newRoute =
                    parseLocation location
            in
            ( { model | route = newRoute }, Cmd.none )


updateFolder : Model -> Folder -> Model
updateFolder model updatedFolder =
    let
        pick currentFolder =
            if updatedFolder.id == currentFolder.id then
                updatedFolder
            else
                currentFolder

        updateFolderList folders =
            List.map pick folders

        updatedFolders =
            RemoteData.map updateFolderList model.folders
    in
    { model | folders = updatedFolders }

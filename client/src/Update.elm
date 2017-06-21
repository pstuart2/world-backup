module Update exposing (..)

import Models exposing (Folder, FolderId, Model, World)
import Msgs exposing (Msg)
import Navigation exposing (newUrl)
import RemoteData
import Routing exposing (getLocationCommand, parseLocation)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msgs.ChangeLocation path ->
            ( model, newUrl path )

        Msgs.OnFetchFolders response ->
            ( { model | folders = response }, Cmd.none )

        Msgs.OnFetchWorlds response ->
            ( updateWorlds model "" response, Cmd.none )

        Msgs.OnLocationChange location ->
            let
                newRoute =
                    parseLocation location

                newCommand =
                    getLocationCommand model newRoute
            in
            ( { model | route = newRoute }, newCommand )


updateWorlds : Model -> FolderId -> RemoteData.WebData (List World) -> Model
updateWorlds model folderId updatedWorlds =
    let
        pick currentFolder =
            if folderId == currentFolder.id then
                { currentFolder | worlds = RemoteData.toMaybe updatedWorlds }
            else
                currentFolder

        updateFolderList folders =
            List.map pick folders

        updatedFolders =
            RemoteData.map updateFolderList model.folders
    in
    { model | folders = updatedFolders }

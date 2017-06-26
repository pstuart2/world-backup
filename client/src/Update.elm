module Update exposing (..)

import Api
import Material
import Models exposing (..)
import Msgs exposing (Msg)
import Navigation exposing (back, newUrl)
import RemoteData
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

        Msgs.OnFetchFolders response ->
            ( { model | folders = response }, Cmd.none )

        Msgs.OnFetchWorlds folderId response ->
            ( updateWorlds model folderId response, Cmd.none )

        Msgs.OnLocationChange location ->
            let
                newRoute =
                    parseLocation location

                newCommand =
                    getLocationCommand model.flags.apiUrl newRoute
            in
            ( { model | route = newRoute }, newCommand )

        Msgs.DeleteBackup folderId worldId backupId ->
            ( model, Api.deleteBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.OnBackupDeleted folderId worldId backupId result ->
            case result of
                Ok _ ->
                    ( deleteBackup model folderId worldId backupId, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )


updateWorlds : Model -> FolderId -> RemoteData.WebData (List World) -> Model
updateWorlds model folderId updatedWorlds =
    let
        pick currentFolder =
            if folderId == currentFolder.id then
                { currentFolder | worlds = updatedWorlds }
            else
                currentFolder

        updateFolderList folders =
            List.map pick folders

        updatedFolders =
            RemoteData.map updateFolderList model.folders
    in
    { model | folders = updatedFolders }


deleteBackup : Model -> FolderId -> WorldId -> BackupId -> Model
deleteBackup model folderId worldId backupId =
    let
        isNotDeleted currentBackup =
            backupId /= currentBackup.id

        findWorld currentWorld =
            if worldId == currentWorld.id then
                { currentWorld | backups = List.filter isNotDeleted currentWorld.backups }
            else
                currentWorld

        updateWorldsList worlds =
            List.map findWorld worlds

        findFolder currentFolder =
            if folderId == currentFolder.id then
                { currentFolder | worlds = RemoteData.map updateWorldsList currentFolder.worlds }
            else
                currentFolder

        updateFolderList folders =
            List.map findFolder folders

        updatedFolders =
            RemoteData.map updateFolderList model.folders
    in
    { model | folders = updatedFolders }

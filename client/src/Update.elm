module Update exposing (..)

import Api
import Debug exposing (log)
import Folders.Update exposing (..)
import Material
import Models exposing (..)
import Msgs exposing (Msg)
import Navigation exposing (back, newUrl)
import Routing exposing (getLocationCommand, parseLocation)
import Task


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
            ( { model | folders = updateFolders model.folders response }, Cmd.none )

        Msgs.OnFetchWorlds folderId response ->
            ( { model | folders = updateWorlds model.folders folderId response }, Cmd.none )

        Msgs.OnLocationChange location ->
            let
                newRoute =
                    parseLocation location

                newCommand =
                    getLocationCommand model.flags.apiUrl newRoute
            in
            ( { model | route = newRoute }, newCommand )

        Msgs.StartWorldDelete worldId ->
            ( { model | folders = setDeleteWorldId (Just worldId) model.folders }, Cmd.none )

        Msgs.DeleteWorld folderId worldId ->
            ( model, Api.deleteWorld model.flags.apiUrl folderId worldId )

        Msgs.CancelDeleteWorld ->
            ( { model | folders = setDeleteWorldId Nothing model.folders }, Cmd.none )

        Msgs.OnWorldDeleted folderId worldId result ->
            case result of
                Ok _ ->
                    ( { model | folders = deleteWorld model.folders folderId worldId }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        Msgs.DeleteBackup folderId worldId backupId ->
            ( model, Api.deleteBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.OnBackupDeleted folderId worldId backupId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, Cmd.none )

                Err world ->
                    let
                        x =
                            Debug.log "error deleting world" world
                    in
                    ( model, Cmd.none )

        Msgs.RestoreBackup folderId worldId backupId ->
            ( model, Api.restoreBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.OnBackupRestored folderId worldId backupId result ->
            case result of
                Ok _ ->
                    ( model, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        Msgs.BackupWorld folderId worldId name ->
            ( model, Api.backupWorld model.flags.apiUrl folderId worldId name )

        Msgs.OnWorldBackedUp folderId worldId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, sendMessage Msgs.CancelWorldBackup )

                Err world ->
                    let
                        x =
                            Debug.log "error backing up world" world
                    in
                    ( model, Cmd.none )

        Msgs.FilterWorlds filter ->
            ( { model | folders = setWorldFilter filter model.folders }, Cmd.none )

        Msgs.ClearWorldsFilter ->
            ( { model | folders = setWorldFilter "" model.folders }, Cmd.none )

        Msgs.StartWorldBackup worldId ->
            ( { model | folders = setCreateBackupId (Just worldId) model.folders }, Cmd.none )

        Msgs.UpdateBackupName name ->
            ( { model | folders = setCreateBackupName name model.folders }, Cmd.none )

        Msgs.CancelWorldBackup ->
            ( { model | folders = setCreateBackupId Nothing model.folders }, Cmd.none )


sendMessage : msg -> Cmd msg
sendMessage msg =
    Task.succeed msg
        |> Task.perform identity

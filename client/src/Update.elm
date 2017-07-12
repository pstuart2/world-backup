module Update exposing (..)

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

        Msgs.FolderMsg msg_ ->
            updateFolder msg_ model

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

        Msgs.OnWorldDeleted folderId worldId result ->
            case result of
                Ok _ ->
                    ( { model | folders = deleteWorld model.folders folderId worldId }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

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

        Msgs.OnBackupRestored folderId worldId backupId result ->
            case result of
                Ok _ ->
                    ( model, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        Msgs.OnWorldBackedUp folderId worldId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, sendMessage (Msgs.FolderMsg Msgs.CancelWorldBackup) )

                Err world ->
                    let
                        x =
                            Debug.log "error backing up world" world
                    in
                    ( model, Cmd.none )


sendMessage : msg -> Cmd msg
sendMessage msg =
    Task.succeed msg
        |> Task.perform identity

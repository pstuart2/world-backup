module Update exposing (..)

import Debug exposing (log)
import Folders.Update exposing (..)
import Material
import Material.Helpers exposing (cssTransitionStep, delay, map1st, map2nd, pure)
import Material.Snackbar as Snackbar
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

        Msgs.AddToast message ->
            addMessage (\k -> Snackbar.toast 0 message) model

        Msgs.Snackbar (Snackbar.Begin k) ->
            ( model, Cmd.none )

        Msgs.Snackbar (Snackbar.End k) ->
            ( model, Cmd.none )

        Msgs.Snackbar (Snackbar.Click k) ->
            ( model, Cmd.none )

        Msgs.Snackbar msg_ ->
            Snackbar.update msg_ model.snackbar
                |> map1st (\s -> { model | snackbar = s })
                |> map2nd (Cmd.map Msgs.Snackbar)

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
                    ( { model | folders = deleteWorld model.folders folderId worldId }, createToast "World has been deleted" )

                Err _ ->
                    ( model, Cmd.none )

        Msgs.OnBackupDeleted folderId worldId backupId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, createToast "Backup has been deleted" )

                Err world ->
                    let
                        x =
                            Debug.log "error deleting world" world
                    in
                    ( model, Cmd.none )

        Msgs.OnBackupRestored folderId worldId backupId result ->
            case result of
                Ok _ ->
                    ( model, createToast "Backup has been restored" )

                Err _ ->
                    ( model, createToast "There was an error restoring the backup" )

        Msgs.OnWorldBackedUp folderId worldId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }
                    , Cmd.batch
                        [ createCommand (Msgs.FolderMsg Msgs.CancelWorldBackup)
                        , createToast "World has been deleted"
                        ]
                    )

                Err world ->
                    let
                        x =
                            Debug.log "error backing up world" world
                    in
                    ( model, Cmd.none )


createCommand : msg -> Cmd msg
createCommand msg =
    Task.succeed msg
        |> Task.perform identity


createToast : String -> Cmd Msg
createToast message =
    createCommand (Msgs.AddToast message)


addMessage : (Int -> Snackbar.Contents Int) -> Model -> ( Model, Cmd Msg )
addMessage f model =
    let
        ( snackbar_, effect ) =
            Snackbar.add (f 0) model.snackbar
                |> map2nd (Cmd.map Msgs.Snackbar)

        model_ =
            { model | snackbar = snackbar_ }
    in
    ( model_
    , effect
    )

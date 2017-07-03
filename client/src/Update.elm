module Update exposing (..)

import Api
import Debug exposing (log)
import Material
import Models exposing (..)
import Msgs exposing (Msg)
import Navigation exposing (back, newUrl)
import RemoteData
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

        Msgs.StartWorldDelete worldId ->
            ( { model | folderView = setDeleteWorldId (Just worldId) model.folderView }, Cmd.none )

        Msgs.DeleteWorld folderId worldId ->
            ( model, Api.deleteWorld model.flags.apiUrl folderId worldId )

        Msgs.CancelDeleteWorld ->
            ( { model | folderView = setDeleteWorldId Nothing model.folderView }, Cmd.none )

        Msgs.OnWorldDeleted folderId worldId result ->
            case result of
                Ok _ ->
                    ( deleteWorld model folderId worldId, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        Msgs.DeleteBackup folderId worldId backupId ->
            ( model, Api.deleteBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.OnBackupDeleted folderId worldId backupId result ->
            case result of
                Ok world ->
                    ( updateWorld model folderId world, Cmd.none )

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
                    ( updateWorld model folderId world, sendMessage Msgs.CancelWorldBackup )

                Err world ->
                    let
                        x =
                            Debug.log "error backing up world" world
                    in
                    ( model, Cmd.none )

        Msgs.FilterWorlds filter ->
            ( { model | folderView = setWorldFilter filter model.folderView }, Cmd.none )

        Msgs.ClearWorldsFilter ->
            ( { model | folderView = setWorldFilter "" model.folderView }, Cmd.none )

        Msgs.StartWorldBackup worldId ->
            ( { model | folderView = setCreateBackupId (Just worldId) model.folderView }, Cmd.none )

        Msgs.UpdateBackupName name ->
            ( { model | folderView = setCreateBackupName name model.folderView }, Cmd.none )

        Msgs.CancelWorldBackup ->
            ( { model | folderView = setCreateBackupId Nothing model.folderView }, Cmd.none )


sendMessage : msg -> Cmd msg
sendMessage msg =
    Task.succeed msg
        |> Task.perform identity


setDeleteWorldId : Maybe WorldId -> FolderView -> FolderView
setDeleteWorldId worldId oldFv =
    { oldFv | deleteWorldId = worldId, backupName = "" }


setCreateBackupId : Maybe WorldId -> FolderView -> FolderView
setCreateBackupId worldId oldFv =
    { oldFv | createBackupId = worldId, backupName = "" }


setCreateBackupName : String -> FolderView -> FolderView
setCreateBackupName name oldFv =
    { oldFv | backupName = name }


setWorldFilter : String -> FolderView -> FolderView
setWorldFilter filter oldFv =
    { oldFv | worldFilter = filter }


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


updateWorld : Model -> FolderId -> World -> Model
updateWorld model folderId updatedWorld =
    let
        findWorld currentWorld =
            if updatedWorld.id == currentWorld.id then
                updatedWorld
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


deleteWorld : Model -> FolderId -> WorldId -> Model
deleteWorld model folderId worldId =
    let
        isNotDeleted currentWorld =
            worldId /= currentWorld.id

        updateWorldsList worlds =
            List.filter isNotDeleted worlds

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

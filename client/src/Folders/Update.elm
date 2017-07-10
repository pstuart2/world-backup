module Folders.Update exposing (..)

import Folders.Api as Api
import Folders.Models exposing (..)
import Folders.Msgs as FolderMsgs
import Models exposing (..)
import Msgs exposing (Msg)
import RemoteData
import Task


-- TODO: This should only care about FolderModel


update : (FolderMsgs.Msg -> Msg) -> FolderMsgs.Msg -> Model -> ( Model, Cmd Msg )
update pMsg msg model =
    case msg of
        FolderMsgs.OnFetchFolders response ->
            ( { model | folders = updateFolders model.folders response }, Cmd.none )

        FolderMsgs.OnFetchWorlds folderId response ->
            ( { model | folders = updateWorlds model.folders folderId response }, Cmd.none )

        FolderMsgs.StartWorldDelete worldId ->
            ( { model | folders = setDeleteWorldId (Just worldId) model.folders }, Cmd.none )

        FolderMsgs.DeleteWorld folderId worldId ->
            ( model, Api.deleteWorld pMsg model.flags.apiUrl folderId worldId )

        FolderMsgs.CancelDeleteWorld ->
            ( { model | folders = setDeleteWorldId Nothing model.folders }, Cmd.none )

        FolderMsgs.OnWorldDeleted folderId worldId result ->
            case result of
                Ok _ ->
                    ( { model | folders = deleteWorld model.folders folderId worldId }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        FolderMsgs.DeleteBackup folderId worldId backupId ->
            ( model, Api.deleteBackup pMsg model.flags.apiUrl folderId worldId backupId )

        FolderMsgs.OnBackupDeleted folderId worldId backupId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, Cmd.none )

                Err world ->
                    let
                        x =
                            Debug.log "error deleting world" world
                    in
                    ( model, Cmd.none )

        FolderMsgs.RestoreBackup folderId worldId backupId ->
            ( model, Api.restoreBackup pMsg model.flags.apiUrl folderId worldId backupId )

        FolderMsgs.OnBackupRestored folderId worldId backupId result ->
            case result of
                Ok _ ->
                    ( model, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        FolderMsgs.BackupWorld folderId worldId name ->
            ( model, Api.backupWorld pMsg model.flags.apiUrl folderId worldId name )

        FolderMsgs.OnWorldBackedUp folderId worldId result ->
            case result of
                Ok world ->
                    ( { model | folders = updateWorld model.folders folderId world }, sendMessage (pMsg FolderMsgs.CancelWorldBackup) )

                Err world ->
                    let
                        x =
                            Debug.log "error backing up world" world
                    in
                    ( model, Cmd.none )

        FolderMsgs.FilterWorlds filter ->
            ( { model | folders = setWorldFilter filter model.folders }, Cmd.none )

        FolderMsgs.ClearWorldsFilter ->
            ( { model | folders = setWorldFilter "" model.folders }, Cmd.none )

        FolderMsgs.StartWorldBackup worldId ->
            ( { model | folders = setCreateBackupId (Just worldId) model.folders }, Cmd.none )

        FolderMsgs.UpdateBackupName name ->
            ( { model | folders = setCreateBackupName name model.folders }, Cmd.none )

        FolderMsgs.CancelWorldBackup ->
            ( { model | folders = setCreateBackupId Nothing model.folders }, Cmd.none )


sendMessage : msg -> Cmd msg
sendMessage msg =
    Task.succeed msg
        |> Task.perform identity


setDeleteWorldId : Maybe WorldId -> FolderModel -> FolderModel
setDeleteWorldId worldId oldFv =
    { oldFv | deleteWorldId = worldId, createBackupId = Nothing, backupName = "" }


setCreateBackupId : Maybe WorldId -> FolderModel -> FolderModel
setCreateBackupId worldId oldFv =
    { oldFv | createBackupId = worldId, deleteWorldId = Nothing, backupName = "" }


setCreateBackupName : String -> FolderModel -> FolderModel
setCreateBackupName name oldFv =
    { oldFv | backupName = name }


setWorldFilter : String -> FolderModel -> FolderModel
setWorldFilter filter oldFv =
    { oldFv | worldFilter = filter }


updateFolders : FolderModel -> RemoteData.WebData (List Folder) -> FolderModel
updateFolders model updatedFolders =
    { model | folders = updatedFolders }


updateWorlds : FolderModel -> FolderId -> RemoteData.WebData (List World) -> FolderModel
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


updateWorld : FolderModel -> FolderId -> World -> FolderModel
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


deleteWorld : FolderModel -> FolderId -> WorldId -> FolderModel
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

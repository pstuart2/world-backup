module Folders.Update exposing (..)

import Api
import Folders.Models exposing (..)
import Models exposing (..)
import Msgs exposing (FolderMsg, Msg)
import RemoteData


updateFolder : FolderMsg -> Model -> ( Model, Cmd Msg )
updateFolder msg model =
    case msg of
        Msgs.StartWorldDelete worldId ->
            ( { model | folders = setDeleteWorldId (Just worldId) model.folders }, Cmd.none )

        Msgs.DeleteWorld folderId worldId ->
            ( model, Api.deleteWorld model.flags.apiUrl folderId worldId )

        Msgs.CancelDeleteWorld ->
            ( { model | folders = setDeleteWorldId Nothing model.folders }, Cmd.none )

        Msgs.DeleteBackup folderId worldId backupId ->
            ( model, Api.deleteBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.RestoreBackup folderId worldId backupId ->
            ( model, Api.restoreBackup model.flags.apiUrl folderId worldId backupId )

        Msgs.BackupWorld folderId worldId name ->
            ( model, Api.backupWorld model.flags.apiUrl folderId worldId name )

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

module Folders.Update exposing (..)

import Folders.Models exposing (..)
import Models exposing (..)
import RemoteData


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

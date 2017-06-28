module Api exposing (..)

import Commands exposing (foldersDecoder, worldDecoder, worldsDecoder)
import Http
import Json.Decode as Decode exposing (Decoder)
import Models exposing (BackupId, Folder, FolderId, WorldId)
import Msgs exposing (Msg)
import RemoteData exposing (WebData)
import Urls


delete : String -> Decode.Decoder a -> Http.Request a
delete url decoder =
    Http.request
        { method = "DELETE"
        , headers = []
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


patch : String -> Http.Body -> Http.Request ()
patch url body =
    Http.request
        { method = "PATCH"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
        }


fetchFolders : String -> Cmd Msg
fetchFolders baseApiUrl =
    Http.get (Urls.folders baseApiUrl) foldersDecoder
        |> RemoteData.sendRequest
        |> Cmd.map Msgs.OnFetchFolders


fetchFolderWorlds : String -> FolderId -> Cmd Msg
fetchFolderWorlds baseApiUrl folderId =
    Http.get (Urls.worlds baseApiUrl folderId) worldsDecoder
        |> RemoteData.sendRequest
        |> Cmd.map (Msgs.OnFetchWorlds folderId)


deleteBackup : String -> FolderId -> WorldId -> BackupId -> Cmd Msg
deleteBackup baseApiUrl folderId worldId backupId =
    let
        request =
            delete (Urls.backup baseApiUrl folderId worldId backupId) worldDecoder
    in
    Http.send (Msgs.OnBackupDeleted folderId worldId backupId) request


restoreBackup : String -> FolderId -> WorldId -> BackupId -> Cmd Msg
restoreBackup baseApiUrl folderId worldId backupId =
    let
        request =
            patch (Urls.backup baseApiUrl folderId worldId backupId) Http.emptyBody
    in
    Http.send (Msgs.OnBackupRestored folderId worldId backupId) request

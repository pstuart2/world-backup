module Api exposing (..)

import Folders.Commands exposing (backupRequestEncoder, foldersDecoder, worldDecoder, worldsDecoder)
import Folders.Models exposing (BackupId, FolderId, WorldId)
import Http
import Json.Decode as Decode exposing (Decoder)
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


post : String -> Http.Body -> Decode.Decoder a -> Http.Request a
post url body decoder =
    Http.request
        { method = "POST"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


deleteNoResponse : String -> Http.Request ()
deleteNoResponse url =
    Http.request
        { method = "DELETE"
        , headers = []
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectStringResponse (\_ -> Ok ())
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


deleteWorld : String -> FolderId -> WorldId -> Cmd Msg
deleteWorld baseApiUrl folderId worldId =
    let
        request =
            deleteNoResponse (Urls.world baseApiUrl folderId worldId)
    in
    Http.send (Msgs.OnWorldDeleted folderId worldId) request


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


backupWorld : String -> FolderId -> WorldId -> String -> Cmd Msg
backupWorld baseApiUrl folderId worldId backupName =
    let
        backupRequestBody =
            backupRequestEncoder { name = backupName }
                |> Http.jsonBody

        request =
            post (Urls.backups baseApiUrl folderId worldId) backupRequestBody worldDecoder
    in
    Http.send (Msgs.OnWorldBackedUp folderId worldId) request

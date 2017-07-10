module Folders.Api exposing (..)

import Folders.Commands exposing (backupRequestEncoder, foldersDecoder, worldDecoder, worldsDecoder)
import Folders.Models exposing (BackupId, FolderId, WorldId)
import Folders.Msgs as FolderMsgs
import Http
import Json.Decode as Decode exposing (Decoder)
import Msgs exposing (Msg)
import RemoteData exposing (WebData)
import Urls


get : String -> Decode.Decoder a -> Http.Request a
get url decoder =
    Http.request
        { method = "GET"
        , headers = []
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


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


fetchFolders : (FolderMsgs.Msg -> Msg) -> String -> Cmd Msg
fetchFolders pMsg baseApiUrl =
    Http.get (Urls.folders baseApiUrl) foldersDecoder
        |> RemoteData.sendRequest
        |> Cmd.map (\a -> pMsg (FolderMsgs.OnFetchFolders a))


fetchFolderWorlds : (FolderMsgs.Msg -> Msg) -> String -> FolderId -> Cmd Msg
fetchFolderWorlds pMsg baseApiUrl folderId =
    Http.get (Urls.worlds baseApiUrl folderId) worldsDecoder
        |> RemoteData.sendRequest
        |> Cmd.map (\a -> pMsg (FolderMsgs.OnFetchWorlds folderId a))


deleteWorld : (FolderMsgs.Msg -> Msg) -> String -> FolderId -> WorldId -> Cmd Msg
deleteWorld pMsg baseApiUrl folderId worldId =
    let
        request =
            deleteNoResponse (Urls.world baseApiUrl folderId worldId)
    in
    Http.send (FolderMsgs.OnWorldDeleted folderId worldId) request
        |> Cmd.map pMsg


deleteBackup : (FolderMsgs.Msg -> Msg) -> String -> FolderId -> WorldId -> BackupId -> Cmd Msg
deleteBackup pMsg baseApiUrl folderId worldId backupId =
    let
        request =
            delete (Urls.backup baseApiUrl folderId worldId backupId) worldDecoder
    in
    Http.send (FolderMsgs.OnBackupDeleted folderId worldId backupId) request
        |> Cmd.map pMsg


restoreBackup : (FolderMsgs.Msg -> Msg) -> String -> FolderId -> WorldId -> BackupId -> Cmd Msg
restoreBackup pMsg baseApiUrl folderId worldId backupId =
    let
        request =
            patch (Urls.backup baseApiUrl folderId worldId backupId) Http.emptyBody
    in
    Http.send (FolderMsgs.OnBackupRestored folderId worldId backupId) request
        |> Cmd.map pMsg


backupWorld : (FolderMsgs.Msg -> Msg) -> String -> FolderId -> WorldId -> String -> Cmd Msg
backupWorld pMsg baseApiUrl folderId worldId backupName =
    let
        backupRequestBody =
            backupRequestEncoder { name = backupName }
                |> Http.jsonBody

        request =
            post (Urls.backups baseApiUrl folderId worldId) backupRequestBody worldDecoder
    in
    Http.send (FolderMsgs.OnWorldBackedUp folderId worldId) request
        |> Cmd.map pMsg

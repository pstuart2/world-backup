module Api exposing (..)

import Commands exposing (foldersDecoder, worldsDecoder)
import Http
import Models exposing (BackupId, FolderId, WorldId)
import Msgs exposing (Msg)
import RemoteData exposing (WebData)
import Urls


delete : String -> Http.Request ()
delete url =
    Http.request
        { method = "DELETE"
        , headers = []
        , url = url
        , body = Http.emptyBody
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
            delete (Urls.backup baseApiUrl folderId worldId backupId)
    in
    Http.send (Msgs.OnBackupDeleted folderId worldId backupId) request

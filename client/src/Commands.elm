module Commands exposing (..)

import Http
import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, required)
import Json.Encode as Encode
import Models exposing (Flags, Folder, FolderId)
import Msgs exposing (Msg)
import RemoteData
import Time.DateTime as DateTime exposing (DateTime)


fetchFolders : String -> Cmd Msg
fetchFolders baseApiUrl =
    Http.get (fetchFoldersUrl baseApiUrl) foldersDecoder
        |> RemoteData.sendRequest
        |> Cmd.map Msgs.OnFetchFolders


fetchFoldersUrl : String -> String
fetchFoldersUrl baseApiUrl =
    baseApiUrl ++ "/folders"


dateTimeDecoder : Decoder DateTime
dateTimeDecoder =
    let
        convert : String -> Decoder DateTime
        convert raw =
            case DateTime.fromISO8601 raw of
                Ok date ->
                    succeed date

                Err error ->
                    fail error
    in
    string |> andThen convert


foldersDecoder : Decode.Decoder (List Folder)
foldersDecoder =
    Decode.list folderDecoder


folderDecoder : Decode.Decoder Folder
folderDecoder =
    decode Folder
        |> required "id" Decode.string
        |> required "path" Decode.string
        |> required "modifiedAt" dateTimeDecoder
        |> required "lastRun" dateTimeDecoder
        |> required "numberOfWorlds" Decode.int


saveFolderUrl : FolderId -> String
saveFolderUrl folderId =
    "/folders/" ++ folderId


saveFolderRequest : Folder -> Http.Request Folder
saveFolderRequest folder =
    Http.request
        { body = folderEncoder folder |> Http.jsonBody
        , expect = Http.expectJson folderDecoder
        , headers = []
        , method = "PATCH"
        , timeout = Nothing
        , url = saveFolderUrl folder.id
        , withCredentials = False
        }



--
-- saveFolderCmd : Folder -> Cmd Msg
-- saveFolderCmd folder =
--     saveFolderRequest folder
--         |> Http.send Msgs.OnFolderSave


folderEncoder : Folder -> Encode.Value
folderEncoder folder =
    let
        attributes =
            [ ( "id", Encode.string folder.id )
            , ( "path", Encode.string folder.path )
            ]
    in
    Encode.object attributes

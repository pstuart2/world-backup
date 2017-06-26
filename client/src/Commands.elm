module Commands exposing (..)

import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Models exposing (Backup, Flags, Folder, FolderId, World)
import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)


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
        |> hardcoded RemoteData.Loading


worldsDecoder : Decode.Decoder (List World)
worldsDecoder =
    Decode.list worldDecoder


worldDecoder : Decode.Decoder World
worldDecoder =
    decode World
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "backups" backupsDecoder


backupsDecoder : Decode.Decoder (List Backup)
backupsDecoder =
    Decode.list backupDecoder


backupDecoder : Decode.Decoder Backup
backupDecoder =
    decode Backup
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "createdAt" dateTimeDecoder

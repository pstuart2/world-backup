module Folders.Commands exposing (..)

import Commands exposing (dateTimeDecoder)
import Folders.Models exposing (..)
import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Json.Encode as Encode
import RemoteData exposing (WebData)


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


backupRequestEncoder : BackupRequest -> Encode.Value
backupRequestEncoder request =
    let
        attributes =
            [ ( "name", Encode.string request.name ) ]
    in
    Encode.object attributes

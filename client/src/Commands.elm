module Commands exposing (..)

import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
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

module Urls exposing (..)

import Models exposing (BackupId, FolderId, WorldId)


folders : String -> String
folders baseApiUrl =
    baseApiUrl ++ "/folders"


folder : String -> FolderId -> String
folder baseApiUrl folderId =
    folders baseApiUrl ++ "/" ++ folderId


worlds : String -> FolderId -> String
worlds baseApiUrl folderId =
    folder baseApiUrl folderId ++ "/worlds"


world : String -> FolderId -> WorldId -> String
world baseApiUrl folderId worldId =
    worlds baseApiUrl folderId ++ "/" ++ worldId


backups : String -> FolderId -> WorldId -> String
backups baseApiUrl folderId worldId =
    world baseApiUrl folderId worldId ++ "/backups"


backup : String -> FolderId -> WorldId -> BackupId -> String
backup baseApiUrl folderId worldId backupId =
    backups baseApiUrl folderId worldId ++ "/" ++ backupId

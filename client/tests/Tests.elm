module Tests exposing (all)

import ElmTest.Extra exposing (..)
import Example exposing (suite)


all : Test
all =
    describe "Sample Test Suite"
        [ suite ]

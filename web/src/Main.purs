module Main where

import Prelude

import Effect (Effect)

foreign import alertMsg :: String -> Effect Unit

main :: Effect Unit
main = do
  alertMsg "hello"

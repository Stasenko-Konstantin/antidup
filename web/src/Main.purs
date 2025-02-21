module Main where

import Prelude

import Effect (Effect)
import Effect.Exception (throw)
import Data.Maybe

import Web.HTML (window)
import Web.HTML.Window (document)
import Web.HTML.HTMLDocument as HD
import Web.HTML.HTMLElement as HE

import Web.DOM.Document as DD
import Web.DOM.Element as DE
import Web.DOM.Node (setTextContent, appendChild)

inBody (Just b) f = f b
inBody Nothing  _ = throw "cannot get body of document"

main :: Effect Unit
main = do
  let titleStr = "antidup"
  
  doc <- window >>= document
  HD.setTitle titleStr doc

  body <- HD.body doc
                              
  title <- DD.createElement "h1" $ HD.toDocument doc
  setTextContent titleStr (DE.toNode title)
  inBody body (\b -> appendChild (DE.toNode title) (HE.toNode b))

  parahraph <- DD.createElement "p" $ HD.toDocument doc
  setTextContent
    "load .zip archive of your pictures and I say what duplicates it have."
    (DE.toNode parahraph)
  inBody body (\b -> appendChild (DE.toNode parahraph) (HE.toNode b))

  

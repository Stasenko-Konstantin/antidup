module Main where

import Prelude

import Effect (Effect)
import Effect.Exception (throw)
import Effect.Console (log)
import Data.Maybe
import Unsafe.Coerce (unsafeCoerce)

import Web.HTML (window)
import Web.HTML.Window (document)
import Web.HTML.HTMLDocument as HD
import Web.HTML.HTMLElement as HE
import Web.HTML.HTMLBaseElement as HBE
import Web.HTML.HTMLInputElement as HIE
import Web.HTML.HTMLLinkElement as HLE

import Web.DOM.Document as DD
import Web.DOM.Element as DE
import Web.DOM.Node (setTextContent, appendChild)
import Web.DOM.Text as DT

inMaybe :: forall a. Maybe a -> (a -> Effect Unit) -> Effect Unit
inMaybe (Just x) f = f x
inMaybe Nothing  _ = throw "cannot get body of document"

main :: Effect Unit
main = do
  let titleStr = "antidup"
  
  doc <- window >>= document
  HD.setTitle titleStr doc

  let createElement s = DD.createElement s (HD.toDocument doc)
  let createTextNode s = DD.createTextNode s (HD.toDocument doc)

  body <- HD.body doc
                              
  title <- createElement "h1"
  setTextContent titleStr (DE.toNode title)
  inMaybe body (\b -> appendChild (DE.toNode title) (HE.toNode b))

  parahraph <- createElement "p"
  setTextContent
    "load .zip archive of your pictures and I say what duplicates it have."
    (DE.toNode parahraph)
  inMaybe body (\b -> appendChild (DE.toNode parahraph) (HE.toNode b))

  needRemoveDuplicatesCheckbox <- createElement "input"
  inMaybe (HIE.fromElement needRemoveDuplicatesCheckbox) (\c -> HIE.setType "checkbox" c) 
  DE.setId "toggleSwitch" needRemoveDuplicatesCheckbox

  needRemoveDuplicatesLabel <- createElement "label"
  appendChild (DE.toNode needRemoveDuplicatesCheckbox) (DE.toNode needRemoveDuplicatesLabel)
  rmDupTextNode <- createTextNode "remove duplicates?"
  appendChild (DT.toNode rmDupTextNode) (DE.toNode needRemoveDuplicatesLabel)
  inMaybe body (\b -> appendChild (DE.toNode needRemoveDuplicatesLabel) (HE.toNode b))

  splitLine <- createElement "hr"
  inMaybe body (\b -> appendChild (DE.toNode splitLine) (HE.toNode b))
  
  author <- createElement "a" -- WARN: forced to use unsafeCoerce because standard conversions doesnt work 
  let linkE = unsafeCoerce author :: HLE.HTMLLinkElement
  HLE.setHref "https://github.com/Stasenko-Konstantin" linkE
  authorText <- createTextNode "author github page"
  appendChild (DT.toNode authorText) (DE.toNode author)
  let baseE = unsafeCoerce author :: HBE.HTMLBaseElement
  HBE.setTarget "_blank" baseE
  inMaybe body (\b -> appendChild (DE.toNode author) (HE.toNode b))
  
  log "run"

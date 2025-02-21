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

main :: Effect Unit
main = do
  let titleStr = "antidup"
  
  doc <- window >>= document
  HD.setTitle titleStr doc

  body <- HD.body doc
                        
  title <- DD.createElement "h1" $ HD.toDocument doc
  setTextContent titleStr (DE.toNode title)
  -- bodyNode <- HE.toNode body
  case body of
    Just b -> do
      appendChild (DE.toNode title) (HE.toNode b)
    Nothing -> pure unit

{-

const title = document.createElement('h1');
title.textContent = `antidup`;
document.body.appendChild(title);

-}

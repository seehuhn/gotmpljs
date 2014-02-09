// JavaScript helper functions for github.com/seehuhn/gotmpljs.
// Copyright (C) 2014  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

goog.provide('seehuhn.gotmpl');

goog.require('goog.string');


/**
 * Escape a string to make it safe for use in an HTML attribute.
 * @param {string} str The string to escape.
 * @return {string} The escaped string.
 */
image.escape.attrescaper = function(str) {
  return goog.string.htmlEscape(str);
};


/**
 * Escape a string to make is safe for use inside an HTML <p> element.
 * @param {string} str The string to escape.
 * @return {string} The escaped string.
 */
image.escape.htmlescaper = function(str) {
  return goog.string.htmlEscape(str);
};


/**
 * Escape a string to make is safe for use inside a URL.
 * @param {string} str The string to escape.
 * @return {string} The escaped string.
 */
image.escape.urlnormalizer = function(str) {
  return str;               // TODO(voss): what needs to be done here?
};

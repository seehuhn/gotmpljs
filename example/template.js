// template.js - generated by github.com/seehuhn/gotmpljs, do not edit

goog.provide('template');

goog.require('seehuhn.gotmpl');


/**
 * Execute the "test" template.
 * @param {*} data The data to apply the template to.
 * @return {string} The template output.
 */
template.test = function(data) {
  var res = new Array();
  res.push('<p><b>');
  res.push(seehuhn.gotmpl.htmlescaper(data['A']));
  res.push('</b> &mdash; ');
  res.push(seehuhn.gotmpl.htmlescaper(data['B']));
  res.push('\n');
  return res.join('');
};
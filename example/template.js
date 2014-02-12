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
  res.push('</b> &mdash; this was the first part of the test\n<ul>');
  for (var i = 0; i < data['B'].length; i++) {
    res.push('\n<li>');
    res.push(seehuhn.gotmpl.htmlescaper(data['B'][i]));
  }
  res.push('\n</ul>\n');
  return res.join('');
};

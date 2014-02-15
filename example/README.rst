GoTmplJS Example Code
=====================

The example code in this directory illustrates use of the GoTmplJs
compiler for compiling Go templates into JavaScript.

1. Install the GoTmplJs compiler::

    go get -u github.com/seehuhn/gotmpljs

2. You will need a copy of the `Google Closure Library`_ source
   installed to run this example.  If you don't have a copy already,
   you can check out a new copy using :code:`git`::

    git clone https://code.google.com/p/closure-library/

   .. _Google Closure Library: https://developers.google.com/closure/library/

3. Compile the test template provided in the file
   :code:`exmaple.html`::

    gotmpljs -n test -o example.js example.html

4. Run the test program::

    go run main.go -c ./closure-library/closure/goog/base.js

   If you are using a pre-installed copy of the Google Closure
   Library, you will need to adjust to Closure Library's 'base.js'
   file in the command above.

5. Check the test output by visiting `<http://localhost:8080/>`_ with
   a web browser.  The output shows two copies of the rendered
   template.  The left-hand column is generated on the server, the
   right-hand column is generated using JavaScript in the browser.

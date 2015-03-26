stdio
=====

Small simple library to adopt ``FILE*`` from libc as standard Go ``io.Reader`` 
and ``io.Writer`` and create ``FILE*`` from Go readers/writers.

This code is adopted from abandoned branch of https://github.com/cookieo9/go-misc/blob/a44c038110b949d742b6ecb8df9c85b6735f5890/cgo/ by Carlos Castillo (licensed under MIT license), and rewritten by Alexander V. Nikolaev for Sendgrid, Inc.

Rewriten version also licensed under MIT license, and distributed as part of 
``go-gmime`` library.

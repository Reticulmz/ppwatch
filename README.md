ppwatch
=======

A simple command-line application to give you updated PP counts as you
play osu!.

Usage
-----

1. Register an API key on the [osu! API page][api], and copy it.
2. Download a ppwatch binary from [the Releases page][releases], and run it.
3. Edit the `.ppwatch.yml` file in your user directory, changing the `EDITME`
   after `api_key:` to the API key you got from the osu! website, and doing the
   same for your username. Leave all the other values as they are.
4. Run the ppwatch binary again, and you will see your PP come in when you finish
   a map!

[api]: https://osu.ppy.sh/p/api
[releases]: https://github.com/txanatan/ppwatch/releases

Compiling
---------

`go get github.com/txanatan/ppwatch`

License
-------

The MIT License (MIT)

Copyright (c) 2015 Alice Jenkinson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

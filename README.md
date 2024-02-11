# getNovel - download novel to txt
[![GoDoc](https://godoc.org/github.com/z-Wind/getNovel?status.png)](http://godoc.org/github.com/z-Wind/getNovel)

## Table of Contents

* [Installation](#installation)
* [Examples](#examples)
* [Support](#support)

## Installation

Please note that because of the goquery dependency, getNovel requires Go1.1+.

    $ go get github.com/z-Wind/getNovel

To build with two ways

    $ cd $GOPATH/src/github.com/z-Wind/getNovel
    $ make

(optional) To run unit tests:

    $ cd $GOPATH/src/github.com/z-Wind/getNovel
    $ make test

(optional) To clean all except source code:

    $ cd $GOPATH/src/github.com/z-Wind/getNovel
    $ make clean

## Examples

    $ cd $GOPATH/src/github.com/z-Wind/getNovel
    $ ./getNovel -url url_novel_contents

## Anti Cloudflare

Currently, only one website has cloudflare installed. You will need to do something extra to crawl from it.

First, visit the site (e.g. https://czbooks.net/) and pass the cloudflare check. After that, open the browser developer tool and check the http request's **Request Headers**.

You will need to copy the following two values represented by `<xxx>`.

```
Cookie: cf_clearance=<base64 encoded string>; ...
User-Agent: <user agent>
```

Then, run the program with new arguments.

```
./getNovel -url URL -cf <base64 encoded string> -ua <user agent>
```

## Support
- [黃金屋](https://tw.hjwzw.com/)
- [飄天文學](https://www.ptwxz.com/)
- [UU看書網](https://www.uukanshu.com/)
- [小說狂人](https://czbooks.net/) !cloudflare!
- ~~[完本神站](https://www.wanbentxt.com/)~~ dead

## Adding new files to test_dataset

	$ cd test_dataset; wget2 -kx URL

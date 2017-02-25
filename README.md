# bypass-webproxy [![Build Status](https://travis-ci.org/pietroglyph/bypass-webproxy.svg?branch=master)](https://travis-ci.org/pietroglyph/bypass-webproxy) [![License](https://img.shields.io/badge/license-MPL--2.0-orange.svg)](https://github.com/pietroglyph/bypass-webproxy/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/pietroglyph/bypass-webproxy)](https://goreportcard.com/report/github.com/pietroglyph/bypass-webproxy)

A simple webproxy written in Go that uses Goquery to parse and modify proxied HTML pages so that links, images, and other resources are fed back through the proxy. Bypass also serves static files.

## Dependencies

+ [goquery](https://github.com/PuerkitoBio/goquery)
+ [osext](https://github.com/kardianos/osext)
+ [iconv-go](https://github.com/djimenez/iconv-go)
+ [go-encoding](https://github.com/pietroglyph/go-encoding) (Forked from [mattn/go-encoding](https://github.com/mattn/go-encoding))

## Building

_TODO: Building Bypass from source_

istools
=======

islint, an intermediate schema linter. [What is linting?](http://stackoverflow.com/questions/8503559/what-is-linting)

Documentation on [godoc.org](https://godoc.org/github.com/miku/istools).

Install current version:

    $ go get github.com/miku/istools/cmd/...

Outdated precompiled Linux 64-bit toy: [islint](https://github.com/miku/istools/releases/download/v0.1.4/islint)

Usage
-----

```sh
$ islint -h
Usage of islint:
  -ls
        list tests
  -sample float
        ratio of records to test (default 1)
  -v    show version and exit
  -verbose
        show every error

$ islint -ls
CurrencyInTitle
EndPageBeforeStartPage
EtAlAuthorName
ExcessivePunctuation
HTMLEntityInAuthorName
InvalidCollection
InvalidEndPage
InvalidStartPage
InvalidURL
KeyTooLong
NAInAuthorName
NonCanonicalISSN
NoPublisher
NoURL
PublicationDateTooEarly
PublicationDateTooLate
RepeatedSlash
RepeatedSubtitle
ShortAuthorName
SuspiciousPageCount
WhitespaceAuthor

$ islint < file.is | jq
{
  "damaged": 53262,
  "dist": {
    "CurrencyInTitle": 2177,
    "EndPageBeforeStartPage": 352,
    "EtAlAuthorName": 29,
    "ExcessivePunctuation": 8,
    "InvalidCollection": 6006,
    "InvalidStartPage": 220,
    "PublicationDateTooEarly": 3680,
    "RepeatedSlash": 13,
    "RepeatedSubtitle": 37501,
    "ShortAuthorName": 4352
  },
  "elapsed": 47.49654878,
  "errcount": {
    "0": 946738,
    "1": 52188,
    "2": 1072,
    "3": 2
  },
  "ratio": "5.326",
  "start": "2015-12-07T18:41:06.50489407+01:00",
  "total": 1000000
}
...
{
  "damaged": 1994583,
  "dist": {
    "CurrencyInTitle": 33179,
    "EndPageBeforeStartPage": 8391,
    "EtAlAuthorName": 1363,
    "ExcessivePunctuation": 337,
    "InvalidCollection": 206737,
    "InvalidEndPage": 387,
    "InvalidStartPage": 9087,
    "InvalidURL": 1,
    "NoPublisher": 1393457,
    "PublicationDateTooEarly": 58379,
    "RepeatedSlash": 6717,
    "RepeatedSubtitle": 242985,
    "ShortAuthorName": 97244,
    "SuspiciousPageCount": 5
  },
  "elapsed": 509.939547991,
  "errcount": {
    "0": 6913680,
    "1": 1931478,
    "2": 62524,
    "3": 581
  },
  "ratio": "22.390",
  "start": "2015-12-07T18:41:06.50489407+01:00",
  "total": 8908263
}
```

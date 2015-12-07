islint
======

Intermediate Schema linter. [What is linting?](http://stackoverflow.com/questions/8503559/what-is-linting)

Documentation on [godoc.org](https://godoc.org/github.com/miku/islint).

Linux 64-bit toy: [islint](https://github.com/miku/islint/releases/download/v0.1.1/islint)

Usage
-----

```sh
$ islint -ls
CurrencyInTitle
EndPageBeforeStartPage
EtAlAuthorName
ExcessivePunctuation
InvalidCollection
InvalidEndPage
InvalidStartPage
InvalidURL
KeyTooLong
NAInAuthorName
NoPublisher
PublicationDateTooEarly
PublicationDateTooLate
RepeatedSlash
RepeatedSubtitle
ShortAuthorName
SuspiciousPageCount
WhitespaceAuthor

$ islint < file.is
2015/12/03 14:45:55 1000000
2015/12/03 14:45:55 1000000 total, 911306 ok, 88694 or 9.733% with issues
2015/12/03 14:45:55 map[SuspiciousPageCount:5 ExcessivePuctuation:5
                        CurrencyInTitle:1007 PublicationDateTooEarly:52361
                        RepeatedSubtitle:18294 EndPageBeforeStartPage:390
                        InvalidStartPage:231 InvalidCollection:16782]
2015/12/03 14:46:47 2000000
2015/12/03 14:46:47 2000000 total, 1786939 ok, 213061 or 11.923% with issues
2015/12/03 14:46:47 map[CurrencyInTitle:5668 InvalidStartPage:685
                        SuspiciousPageCount:5 PublicationDateTooEarly:146781
                        RepeatedSubtitle:34849 EndPageBeforeStartPage:5146
                        InvalidCollection:20572 ExcessivePuctuation:381
                        InvalidEndPage:7]
2015/12/03 14:47:37 3000000
2015/12/03 14:47:37 3000000 total, 2651675 ok, 348325 or 13.136% with issues
2015/12/03 14:47:37 map[PublicationDateTooEarly:195313 RepeatedSubtitle:118735
                        EndPageBeforeStartPage:5712 InvalidCollection:21339
                        ExcessivePuctuation:388 InvalidEndPage:7
                        CurrencyInTitle:7511 InvalidStartPage:731
                        SuspiciousPageCount:5]
...
2015/12/03 16:01:27 88521109 total, 83013026 ok, 5508083 or 6.635% with issues
2015/12/03 16:01:27 map[CurrencyInTitle:330554 InvalidStartPage:90924
                        SuspiciousPageCount:63 InvalidURL:37
                        PublicationDateTooLate:4 PublicationDateTooEarly:582577
                        RepeatedSubtitle:2402953 EndPageBeforeStartPage:81716
                        InvalidCollection:2060252 ExcessivePuctuation:3169
                        InvalidEndPage:3856]
```

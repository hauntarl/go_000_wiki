# Writing Web Applications

Implementation of [wiki tutorial](https://golang.org/doc/articles/wiki/), along with some bonus tasks

**Logging:**

``` terminal
D:\godemo\wiki>go run wiki.go
2021/04/01 18:17:48 /                              : redirected to front page
2021/04/01 18:17:49 /view/FrontPage                : redirected to edit page
2021/04/01 18:17:49 FrontPage                      : file opened in edit mode
2021/04/01 18:19:13 FrontPage                      : file saved succesfully
2021/04/01 18:19:13 FrontPage                      : file displayed
2021/04/01 18:19:29 /view/sample                   : redirected to edit page
2021/04/01 18:19:29 sample                         : file opened in edit mode
2021/04/01 18:21:20 sample                         : file saved succesfully
2021/04/01 18:21:20 sample                         : file displayed
2021/04/01 18:21:29 FrontPage                      : file displayed
2021/04/01 18:21:54 /view/invalid_name             : path not found
2021/04/01 18:22:31 ANewPage                       : file opened in edit mode
2021/04/01 18:23:52 ANewPage                       : file saved succesfully
2021/04/01 18:23:52 ANewPage                       : file displayed
2021/04/01 18:24:10 FrontPage                      : file displayed
2021/04/01 18:25:21 FrontPage                      : file opened in edit mode
2021/04/01 18:25:54 FrontPage                      : file saved succesfully
2021/04/01 18:25:54 FrontPage                      : file displayed
exit status 3221225786
```

Note: the generated data files are not pushed on Github

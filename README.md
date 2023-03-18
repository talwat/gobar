# gobar

[![godocs](https://godoc.org/github.com/talwat/gobar?status.svg)](https://godoc.org/github.com/talwat/gobar)

A minimal and extremely hackable progressbar for golang.

## Install

```sh
go get -u github.com/talwat/gobar
```

## Why?

Originally, I was using [progressbar](https://github.com/schollz/progressbar) by [schollz](https://github.com/schollz) in my programs.

However, that library has (mostly) been unmaintained, only with PR's being merged and dependencies being updated from time to time.

And this would have been fine, except that [I was having a bug](https://github.com/schollz/progressbar/issues/155) which meant that it
put my other program to a halt because I couldn't do anything else.

So, after waiting a while, I decided to just make my own progressbar with only **~70** sloc!
(In comparison to **~800** of [progressbar](https://github.com/schollz/progressbar))

Of course this comes at the cost of customizability, but I think it's worth it in my opinion to at least have something that works for me.

## Examples

### Basic increment

```go
bar := gobar.NewBar(0, 100, "basic", "done!")

for i := 0; i < 100; i++ {
    time.Sleep(10 * time.Millisecond)
    bar.Increment(1)
}
```

### IO operations

```go
req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
resp, _ := http.DefaultClient.Do(req)
defer resp.Body.Close()

f, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
defer f.Close()

bar := gobar.NewBar(0, resp.ContentLength, "io", "done!")
io.Copy(io.MultiWriter(f, bar), resp.Body)
```

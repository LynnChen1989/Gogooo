# Gogooo
golang

# 交叉编译

```
GOOS=linux GOARCH=amd64 go build -o tt t1.go

# 编译到 linux 64bit
$ GOOS=linux GOARCH=amd64 go build

# 或者可以使用 -o 选项指定生成二进制文件名字
$ GOOS=linux GOARCH=amd64 go build -o app.linux

# 编译到 linux 32bit
$ GOOS=linux GOARCH=386 go build

# 编译到 windows 64bit
$ GOOS=windows GOARCH=amd64 go build

# 编译到 windows 32bit
$ GOOS=windows GOARCH=386 go build

# 编译到 Mac OS X 64bit
$ GOOS=darwin GOARCH=amd64 go build

```
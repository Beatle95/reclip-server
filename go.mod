module reclip/main

go 1.23.2

require internal v1.0.0

replace internal => ./internal

require communication v1.0.0

require github.com/gammazero/deque v1.0.0 // indirect

replace communication => ./communication

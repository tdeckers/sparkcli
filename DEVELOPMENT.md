# Development

Build and deploy using goxc: [See here.](https://github.com/laher/goxc/blob/master/README.md)

## Install go from source

    git clone https://go.googlesource.com/go
    git branch go1.5
    cd src
    ./all.bash
    go get golang.org/x/tools/cmd/...

## Install goxc`

    go install github.com/laher/goxc

## Run goxc

Inside the project directory, run:

    goxc

## bintray uploads

Add API key to .goxc.local.json

    goxc bintray

Configuration for bintray plugin is in .goxc.yml.  API key is in 
.goxc.local.yml (not checked in!).  Format:

    {
	"ConfigVersion": "0.9",
	"TaskSettings": {
		"bintray": {
                "apikey": "5d1f300712a5da07b2f64109921cc0346622e14c"
            }
	}
    }

#  TODO

* travis builds
* gocover.io
* godoc creation
* unit testing - [gotests](https://github.com/cweill/gotests)

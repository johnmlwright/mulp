#!/usr/bin/env bash

package=$1
package_name=${package%????}
platforms=("linux/386"	"linux/amd64"	"linux/arm"	"linux/arm64"	"linux/ppc64"	"linux/ppc64le"	"linux/mips"	"linux/mipsle" "linux/mips64"	"linux/mips64le"	"solaris/amd64"	"windows/386"	"windows/amd64"	)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    env GOOS=$GOOS GOARCH=$GOARCH go build -o exec\/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
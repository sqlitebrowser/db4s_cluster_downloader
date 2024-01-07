#!/bin/sh

# Useful variables
DEST=${PWD}/local
export PKG_CONFIG_PATH=${DEST}/lib/pkgconfig
export GOBIN=${DEST}/bin

# If this script is passed an argument of "clean", then delete the
# locally compiled pieces
if [ "$1" = "clean" ]; then
  echo "Removing local SQLite and compiled DBHub.io executables"
  rm -rf ${DEST} other/cache
  exit
fi

# Builds a local SQLite
if [ ! -e "${DEST}/lib/libsqlite3.so" ]; then
  if [ ! -d "other/cache" ]; then
    mkdir -p other/cache
  fi
  cd other/cache  || exit 1
  if [ ! -f sqlite.tar.gz ]; then
    echo "Downloading SQLite source code"
    TARBALL=$(curl -s https://sqlite.org/download.html | awk '/<!--/,/-->/ {print}' | grep 'sqlite-autoconf' | cut -d ',' -f 3)
    SHA3=$(curl -s https://sqlite.org/download.html | awk '/<!--/,/-->/ {print}' | grep 'sqlite-autoconf' | cut -d ',' -f 5)
    curl -LsS -o sqlite.tar.gz https://sqlite.org/${TARBALL}
    VERIFY=$(openssl dgst -sha3-256 sqlite.tar.gz | cut -d ' ' -f 2)
    if [ "$SHA3" != "$VERIFY" ]; then exit 2 ; fi
  fi
  if [ ! -f sqlite.tar.gz ]; then
    echo "Downloading the SQLite source code did not work"
    exit 3
  fi
  echo "Compiling local SQLite"
  tar xfz sqlite.tar.gz
  cd sqlite-autoconf-* || exit 4
  CPPFLAGS="-DSQLITE_ENABLE_COLUMN_METADATA=1 -DSQLITE_MAX_VARIABLE_NUMBER=250000 -DSQLITE_ENABLE_RTREE=1 -DSQLITE_ENABLE_GEOPOLY=1 -DSQLITE_ENABLE_FTS3=1 -DSQLITE_ENABLE_FTS3_PARENTHESIS=1 -DSQLITE_ENABLE_FTS5=1 -DSQLITE_ENABLE_STAT4=1 -DSQLITE_ENABLE_JSON1=1 -DSQLITE_SOUNDEX=1 -DSQLITE_ENABLE_MATH_FUNCTIONS=1 -DSQLITE_MAX_ATTACHED=125 -DSQLITE_ENABLE_MEMORY_MANAGEMENT=1 -DSQLITE_ENABLE_SNAPSHOT=1" ./configure --prefix=${DEST} --enable-dynamic-extensions=no
  make -j "$(nproc)"
  make install
  cd ..
  rm -rf sqlite-autoconf-*
  cd ../..
fi

# Builds the Go binaries
echo "Compiling DB4S downloader daemon"
go install .

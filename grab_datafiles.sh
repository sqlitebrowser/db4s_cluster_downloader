#!/usr/bin/env sh

# This is a simple script to download the DB4S release files from GitHub, so they're present for the GitHub Actions
# based Go test workflow

# Immediately error out if any of the commands doesn't succeed
set -e

# Download the release files
mkdir -p data
cd data

OPTIONS="-sSOL"

# 3.10.1
for file in \
    DB.Browser.for.SQLite-3.10.1-win32.exe \
    DB.Browser.for.SQLite-3.10.1-win64.exe \
    DB.Browser.for.SQLite-3.10.1.dmg \
    SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.10.1/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.11.0
for file in \
    DB.Browser.for.SQLite-3.11.0-win32.msi \
    DB.Browser.for.SQLite-3.11.0-win32.zip \
    DB.Browser.for.SQLite-3.11.0-win64.msi \
    DB.Browser.for.SQLite-3.11.0-win64.zip \
    DB.Browser.for.SQLite-3.11.0.dmg; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.11.1
for file in \
    DB.Browser.for.SQLite-3.11.1-win32.msi \
    DB.Browser.for.SQLite-3.11.1-win32.zip \
    DB.Browser.for.SQLite-3.11.1-win64.msi \
    DB.Browser.for.SQLite-3.11.1-win64.zip \
    DB.Browser.for.SQLite-3.11.1v2.dmg; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.11.2
for file in \
    DB.Browser.for.SQLite-3.11.2-win32.msi \
    DB.Browser.for.SQLite-3.11.2-win32.zip \
    DB.Browser.for.SQLite-3.11.2-win64.msi \
    DB.Browser.for.SQLite-3.11.2-win64.zip \
    DB.Browser.for.SQLite-3.11.2.dmg \
    SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.12.0
for file in \
    DB.Browser.for.SQLite-3.12.0-win32.msi \
    DB.Browser.for.SQLite-3.12.0-win32.zip \
    DB.Browser.for.SQLite-3.12.0-win64.msi \
    DB.Browser.for.SQLite-3.12.0-win64.zip \
    DB.Browser.for.SQLite-3.12.0.dmg \
    SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.12.2
for file in \
    DB.Browser.for.SQLite-3.12.2-win32.msi \
    DB.Browser.for.SQLite-3.12.2-win32.zip \
    DB.Browser.for.SQLite-3.12.2-win64.msi \
    DB.Browser.for.SQLite-3.12.2-win64.zip \
    DB.Browser.for.SQLite-3.12.2.dmg \
    DB.Browser.for.SQLite-arm64-3.12.2.dmg \
    DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage \
    SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.13.0
for file in \
    DB.Browser.for.SQLite-v3.13.0.dmg \
    DB.Browser.for.SQLite-v3.13.0-win32.msi \
    DB.Browser.for.SQLite-v3.13.0-win32.zip \
    DB.Browser.for.SQLite-v3.13.0-win64.msi \
    DB.Browser.for.SQLite-v3.13.0-win64.zip \
    DB.Browser.for.SQLite-v3.13.0-x86.64.AppImage; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done

# 3.13.1
for file in \
    DB.Browser.for.SQLite-v3.13.1.dmg \
    DB.Browser.for.SQLite-v3.13.1-win32.msi \
    DB.Browser.for.SQLite-v3.13.1-win32.zip \
    DB.Browser.for.SQLite-v3.13.1-win64.msi \
    DB.Browser.for.SQLite-v3.13.1-win64.zip \
    DB.Browser.for.SQLite-v3.13.1-x86.64-v2.AppImage; do
  if [ ! -s "${file}" ]; then
    echo
    echo "Downloading ${file}"
    curl "${OPTIONS}" "https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.1/${file}"
  else
    echo " * ${file} already downloaded"
  fi
done
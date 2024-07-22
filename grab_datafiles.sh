#!/usr/bin/env sh

# This is a simple script to download the DB4S release files from GitHub, so they're present for the GitHub Actions
# based Go test workflow

# Immediately error out if any of the commands doesn't succeed
set -e

# Download the release files
mkdir data
cd data

# 3.10.1
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.10.1/DB.Browser.for.SQLite-3.10.1-win32.exe
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.10.1/DB.Browser.for.SQLite-3.10.1-win64.exe
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.10.1/DB.Browser.for.SQLite-3.10.1.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.10.1/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe

# 3.11.0
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/DB.Browser.for.SQLite-3.11.0-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/DB.Browser.for.SQLite-3.11.0-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/DB.Browser.for.SQLite-3.11.0-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/DB.Browser.for.SQLite-3.11.0-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.0/DB.Browser.for.SQLite-3.11.0.dmg

# 3.11.1
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/DB.Browser.for.SQLite-3.11.1-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/DB.Browser.for.SQLite-3.11.1-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/DB.Browser.for.SQLite-3.11.1-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/DB.Browser.for.SQLite-3.11.1-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.1/DB.Browser.for.SQLite-3.11.1v2.dmg

# 3.11.2
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/DB.Browser.for.SQLite-3.11.2-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/DB.Browser.for.SQLite-3.11.2-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/DB.Browser.for.SQLite-3.11.2-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/DB.Browser.for.SQLite-3.11.2-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/DB.Browser.for.SQLite-3.11.2.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.11.2/SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe

# 3.12.0
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/DB.Browser.for.SQLite-3.12.0-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/DB.Browser.for.SQLite-3.12.0-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/DB.Browser.for.SQLite-3.12.0-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/DB.Browser.for.SQLite-3.12.0-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/DB.Browser.for.SQLite-3.12.0.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.0/SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe

# 3.12.2
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-3.12.2-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-3.12.2-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-3.12.2-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-3.12.2-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-3.12.2.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB.Browser.for.SQLite-arm64-3.12.2.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.12.2/SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe

# 3.13.0
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0.dmg
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0-win32.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0-win32.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0-win64.msi
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0-win64.zip
curl -sSOL https://github.com/sqlitebrowser/sqlitebrowser/releases/download/v3.13.0/DB.Browser.for.SQLite-v3.13.0-x86.64.AppImage

# Database schema for the database stats

This creates an empty database ready to log DB4S downloads

To load this into a fresh PostgreSQL database run these commands
from the postgres superuser:

    $ createuser -U postgres -d db4s
    $ createdb -U postgres -O db4s db4s_stats
    $ psql -U db4s db4s_stats < schema.sql

It should finish with no errors.

Note - This schema is created using:

    $ pg_dump -Os -U postgres db4s_stats > schema.sql

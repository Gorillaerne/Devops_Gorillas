from dotenv import load_dotenv
load_dotenv()

import mysql.connector
import psycopg2
from psycopg2.extras import execute_values
import os

# --- Connection config (use env vars, not hardcoded credentials) ---
mysql_conn = mysql.connector.connect(
    host=os.getenv("MYSQL_HOST", "127.0.0.1"),
    port=int(os.getenv("MYSQL_PORT", 3306)),
    user=os.getenv("MYSQL_USER", "root"),
    password=os.getenv("MYSQL_PASSWORD", ""),
    database=os.getenv("MYSQL_DB", "whoknows")
)
pg_conn = psycopg2.connect(
    host=os.getenv("PG_HOST", "127.0.0.1"),
    port=int(os.getenv("PG_PORT", 5432)),
    user=os.getenv("PG_USER", "root"),
    password=os.getenv("PG_PASSWORD", ""),
    dbname=os.getenv("PG_DB", "whoknows")
)

BATCH_SIZE = 1000  # Rows per batch insert

mysql_cursor = mysql_conn.cursor(dictionary=True)
pg_cursor = pg_conn.cursor()

# --- Get all tables ---
mysql_cursor.execute("SHOW TABLES")
tables = [list(row.values())[0] for row in mysql_cursor.fetchall()]
print(f"Found tables: {tables}\n")

for table in tables:
    print(f"Migrating table: {table}")

    # Get total row count for progress tracking
    mysql_cursor.execute(f"SELECT COUNT(*) as cnt FROM {table}")
    total = mysql_cursor.fetchone()["cnt"]
    print(f"  Total rows: {total}")

    if total == 0:
        print(f"  Empty table, skipping.\n")
        continue

    # Stream rows in batches instead of fetchall()
    mysql_cursor.execute(f"SELECT * FROM {table}")

    inserted = 0
    batch = []

    first_row = mysql_cursor.fetchone()
    if not first_row:
        continue

    columns = list(first_row.keys())
    # Quote every column name to handle reserved words
    cols_quoted = ", ".join(f'"{col}"' for col in columns)

    batch.append(tuple(first_row.values()))

    for row in mysql_cursor:
        batch.append(tuple(row.values()))

        if len(batch) >= BATCH_SIZE:
            try:
                execute_values(
                    pg_cursor,
                    f'INSERT INTO "{table}" ({cols_quoted}) VALUES %s ON CONFLICT DO NOTHING',
                    batch
                )
                pg_conn.commit()
                inserted += len(batch)
                print(f"  Inserted {inserted}/{total} rows...")
            except Exception as e:
                pg_conn.rollback()
                print(f"  ERROR on batch insert into {table}: {e}")
                raise
            batch = []

    # Insert remaining rows
    if batch:
        try:
            execute_values(
                pg_cursor,
                f'INSERT INTO "{table}" ({cols_quoted}) VALUES %s ON CONFLICT DO NOTHING',
                batch
            )
            pg_conn.commit()
            inserted += len(batch)
        except Exception as e:
            pg_conn.rollback()
            print(f"  ERROR on final batch insert into {table}: {e}")
            raise

    print(f"  Done. Inserted {inserted} rows.\n")

    # --- Reset sequences for auto-increment columns ---
    pg_cursor.execute(f"""
        SELECT column_name
        FROM information_schema.columns
        WHERE table_name = %s
          AND column_default LIKE 'nextval%%'
    """, (table,))
    seq_columns = pg_cursor.fetchall()

    for (col,) in seq_columns:
        pg_cursor.execute(f"""
            SELECT setval(
                pg_get_serial_sequence('"{table}"', %s),
                COALESCE((SELECT MAX("{col}") FROM "{table}"), 1)
            )
        """, (col,))
        pg_conn.commit()
        print(f"  Reset sequence for {table}.{col}")

print("Migration complete!")

mysql_cursor.close()
pg_cursor.close()
mysql_conn.close()
pg_conn.close()
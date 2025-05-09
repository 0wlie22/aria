import glob
import hashlib
import os
import re
from datetime import datetime

import numpy as np
import pandas as pd
import psycopg2

from categories import expense_categories, income_categories, skip_categories

LOG_LEVELS = {"DEBUG": 10, "INFO": 20, "WARNING": 30, "ERROR": 40}
DEFAULT_LOG_LEVEL = "INFO"


def get_log_level():
    return os.getenv("LOG_LEVEL", DEFAULT_LOG_LEVEL).upper()


def log(message, level="INFO"):
    current_level = get_log_level()
    if LOG_LEVELS[level] >= LOG_LEVELS.get(current_level, 20):
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"{timestamp} [{level}] - {message}")


def extract_date_from_narrative(narrative):
    match = re.search(r"\b(\d{2}/\d{2}/\d{4})\b", str(narrative))
    return match.group(1).replace("/", ".") if match else None


def classify_category(narrative, amount):
    narrative_lower = str(narrative).lower()

    for skip in skip_categories:
        if skip in narrative_lower:
            return "skip"

    categories_to_check = income_categories if amount > 0 else expense_categories

    for category, keywords in categories_to_check.items():
        if any(kw in narrative_lower for kw in keywords):
            return category

    return None


def determine_type(amount):
    if amount > 0:
        return "income"
    elif amount < 0:
        return "expense"
    else:
        return None


def create_fingerprint(row):
    raw_str = (
        f"{row['type']}_{row['date']}_{row['amount']}_{str(row['Narrative'])[:50]}"
    )
    return hashlib.sha256(raw_str.lower().strip().encode("utf-8")).hexdigest()


def transaction_exists(cursor, fingerprint):
    cursor.execute("SELECT 1 FROM transactions WHERE fingerprint = %s", (fingerprint,))
    exists = cursor.fetchone() is not None
    log(
        f"Checked if transaction exists for fingerprint {fingerprint[:6]}...: {exists}",
        "DEBUG",
    )
    return exists


def load_and_preprocess_csv(file_path):
    df = pd.read_csv(file_path, skiprows=3, sep="|", encoding="cp1257")
    df = df.iloc[:-5]

    df["narrative_date"] = df["Narrative"].apply(extract_date_from_narrative)
    df["amount"] = (
        df["Amount DR/Amount CR"].replace(",", ".", regex=True) * 100
    ).astype(int)
    df["category"] = df.apply(
        lambda row: classify_category(row["Narrative"], row["amount"]), axis=1
    )
    df["date"] = np.where(
        df["narrative_date"].notnull(), df["narrative_date"], df["Date"]
    )
    df["date"] = pd.to_datetime(df["date"], format="%d.%m.%Y").dt.date
    df["type"] = df.apply(lambda row: determine_type(row["amount"]), axis=1)

    log(f"Preprocessed CSV {os.path.basename(file_path)}: {len(df)} rows", "DEBUG")
    return df


def ask_user_for_category(row):
    print("Input needed -----------------------------------")
    print(f"  Narrative: {row['Narrative']}")
    print(f"  Amount: {row['amount']/100}")
    print(f"  Date: {row['date']}")
    print(f"  Fingerprint: {row['fingerprint']}")

    if row["amount"] > 0:
        valid_categories = list(income_categories.keys())
    else:
        valid_categories = list(expense_categories.keys())

    user_input = (
        input(f"  Enter category {valid_categories} or ' ' to skip: ").strip().lower()
    )
    if user_input not in valid_categories and user_input != "":
        print(f"Invalid category '{user_input}'. Please try again.")
        return ask_user_for_category(row)

    log(f"User input category: {user_input}", "DEBUG")

    return user_input


def insert_transaction(cursor, row):
    cursor.execute(
        """
        INSERT INTO transactions (fingerprint, type, date, amount, category)
        VALUES (%s, %s, %s, %s, %s)
    """,
        (row["fingerprint"], row["type"], row["date"], row["amount"], row["category"]),
    )


def main():
    conn = psycopg2.connect(
        dbname="aria",
        user="finance_user",
        password=os.getenv("PSQL_PASSWORD"),
        host="localhost",
        port=5432,
    )
    cursor = conn.cursor()

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS transactions (
        fingerprint TEXT PRIMARY KEY,
        type TEXT,
        date DATE,
        amount INTEGER,
        category TEXT
    );""")
    conn.commit()
    log("Database initialized (transactions table ensured)", "INFO")

    csv_files = glob.glob("./data/*.csv")
    log(f"Found {len(csv_files)} CSV files", "INFO")

    for file_path in csv_files:
        log(f"Processing file: {os.path.basename(file_path)}", "INFO")
        df = load_and_preprocess_csv(file_path)

        # Create fingerprints (initial, maybe incomplete)
        df["fingerprint"] = df.apply(create_fingerprint, axis=1)

        # Manual input for uncategorized
        for idx, row in df[df["category"].isnull()].iterrows():
            if transaction_exists(cursor, row["fingerprint"]):
                continue

            category = ask_user_for_category(row)
            df.at[idx, "category"] = category

        # Update type + fingerprint after category fixes
        df["type"] = df.apply(lambda row: determine_type(row["amount"]), axis=1)
        df["fingerprint"] = df.apply(create_fingerprint, axis=1)

        df_clean = df[(df["category"].notnull()) & (df["type"].notnull())].copy()
        log(f"Transactions ready for insertion: {len(df_clean)}", "INFO")

        for _, row in df_clean.iterrows():
            try:
                if row["category"] in ["", "skip"]:
                    continue
                if transaction_exists(cursor, row["fingerprint"]):
                    continue
                insert_transaction(cursor, row)
            except psycopg2.IntegrityError:
                log(
                    f"IntegrityError: duplicate fingerprint {row['fingerprint'][:6]}...",
                    "WARNING",
                )
                continue

        conn.commit()

    df_sqlite = pd.read_sql_query(
        "SELECT type, category, SUM(amount)/100 as total_amount FROM transactions GROUP BY type, category",
        conn,
    )
    log("Final data from SQLite:", "INFO")
    print(df_sqlite)

    conn.close()
    log("Database connection closed", "INFO")


if __name__ == "__main__":
    main()

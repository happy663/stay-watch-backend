from __future__ import annotations
import pandas as pd
from datetime import date
from sqlalchemy.orm import Session
from . import struct as st

# userIdと期間を指定してlogを取得
def get_log_by_userId_and_period(db: Session, userId: int, start: date, end: date) -> pd.DataFrame:
    # logsを取得
    q = db.query(st.Log.start_at, st.Log.end_at, st.Log.user_id).filter(st.Log.user_id == userId, st.Log.start_at >= start, st.Log.end_at < end)
    logs = pd.read_sql(q.statement, con=db.connection()) # type: ignore
    return logs

# userIdから全てのlogを取得
def get_all_logs_by_userId(db: Session, userId: int) -> pd.DataFrame:
    # logsを取得
    q = db.query(st.Log.start_at, st.Log.end_at, st.Log.user_id).filter(st.Log.user_id == userId)
    logs = pd.read_sql(q.statement, con=db.connection()) # type: ignore
    return logs
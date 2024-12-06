from __future__ import annotations
from sqlalchemy.orm import Session
import pandas as pd
from . import struct as st

# EditedLogsからあるuserの最も古いデータを取得
def get_oldest_editedLog_by_userId(db: Session, userId: int) -> st.EditedLog | None:
    log: st.EditedLog | None = db.query(_entity=st.EditedLog).filter(st.EditedLog.user_id == userId).order_by(st.EditedLog.date).first()
    return log

# userIdから特定日のデータを取得し、dateframeに変換
def get_editedLog_by_userId_and_day(db: Session, userId: int, day: int) -> pd.DataFrame:
    # logsを取得
    q = db.query(_entity=st.EditedLog).filter(st.EditedLog.user_id == userId)
    edited_logs = pd.read_sql(q.statement, con=db.connection()) # type: ignore
    date = pd.to_datetime(edited_logs['date']) # type: ignore
    edited_logs['day'] = date.dt.weekday
    edited_logs = edited_logs[edited_logs['day'] == day]
    return edited_logs

# EditedLogを作成する
def create_editedLog(db: Session, uid: int, date: str, reporting_time: str, leavr_time: str) -> st.EditedLog:
    # editedLogを作成
    edited_log: st.EditedLog = st.EditedLog(user_id=uid, date=date, reporting=reporting_time, leave=leavr_time)
    db.add(edited_log)
    db.commit()
    db.refresh(edited_log)
    return edited_log
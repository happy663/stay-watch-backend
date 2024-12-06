from __future__ import annotations
from sqlalchemy import select
from sqlalchemy.orm import Session
from . import struct as st

# 全てのユーザを取得する
def get_all_users_by_community(community:int, db: Session) -> list[st.User]:
    users: list[st.User] = db.query(st.User).where(st.User.community_id == community).all()
    return users

# 特定のユーザを取得する
def get_user_by_id(db: Session, user_id: int) -> st.User | None:
    user: st.User | None = db.scalar(select(st.User).where(st.User.id == user_id))
    return user
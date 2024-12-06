from __future__ import annotations
from sqlalchemy import select
from sqlalchemy.orm import Session
import datetime
from . import struct as st

# userIdから最新のクラスタリングを取得する
def get_latest_cluster_by_userId(db: Session, userId: int, reporting:bool) -> st.Cluster | None:
    # clustersを取得
    cluster: st.Cluster | None = db.scalar(select(st.Cluster).where(st.Cluster.user_id == userId, st.Cluster.reporting == reporting).order_by(st.Cluster.date.desc()))
    return cluster

# userIdから最古のクラスタリングを取得する
def get_oldest_cluster_by_userId(db: Session, userId: int, reporting:bool) -> st.Cluster | None:
    # clustersを取得
    cluster: st.Cluster | None = db.scalar(select(st.Cluster).where(st.Cluster.user_id == userId, st.Cluster.reporting == reporting).order_by(st.Cluster.date))
    return cluster

# 取得したclusterと同じuser_id、dateのclusterを全て取得する
def get_all_cluster_by_userId_and_date(db: Session, userId: int, date: datetime.date, reporting:bool) -> list[st.Cluster]:
    # clustersを取得
    clusters: list[st.Cluster] = db.query(st.Cluster).where(st.Cluster.user_id == userId, st.Cluster.reporting == reporting, st.Cluster.date == date).all()
    return clusters

# clusterを作成する
def create_cluster(db: Session, uid: int, day: datetime.date, reporting: bool, mean: float, std: float, count: int) -> st.Cluster:
    # clusterを作成
    cluster = st.Cluster(user_id=uid, date=day, reporting=reporting, average=mean, sd=std, count=count)
    db.add(cluster)
    db.commit()
    db.refresh(cluster)
    return cluster
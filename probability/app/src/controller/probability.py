from __future__ import annotations
import math
from typing import List
from pydantic import BaseModel
from fastapi import APIRouter, Depends, HTTPException
from fastapi.encoders import jsonable_encoder
from fastapi.responses import ORJSONResponse
from sqlalchemy.orm import Session
from datetime import datetime, timedelta
from models import cluster as cl, user as us
from models import editedlog as el
from service import normal_distribution as nd
from lib.mysql import get_db

# レスポンス用のクラス
class ProbabilityResponse(BaseModel):
    userId: int
    userName: str
    probability: float

router = APIRouter()

# フロント側から現在の時刻を受け取ることとする（後で要変更)
# 今日のある時間までに特定のユーザーが入退室する確率、もしくはある時間以降に入退室する確率を返す
# 変数：user_id, true or false
@router.get("/app/probability/{reporting}/{before}" , response_class=ORJSONResponse, response_model=ProbabilityResponse)
async def get_probability_reporting_before(reporting:str, before:str, user_id:int = 0, date:str = '2024-1-1', time:str = "24:00:00", db: Session = Depends(get_db)):
    r = True if reporting == "reporting" else False
    b = True if before == "before" else False
    date_object= datetime.strptime(date, '%Y-%m-%d') # 今日の日付
    seven_days_ago= date_object - timedelta(days=7) # 1週間前の日付
    # userのAPIが叩かれた曜日の最新のクラスタを取得(1週間前までのデータで作成されたクラスタ)
    clusters = cl.get_all_cluster_by_userId_and_date(db, user_id, seven_days_ago, r)

    # 確率を算出してreturnを返す
        # 確率算出のためにデータ収集期間を求める
        # データが存在している場合は確率計算に進む
        # データが存在していない = userが追加されたばかりなので確率=0%
    user = us.get_user_by_id(db, user_id)
    if user is not None:
        # userの最も古いデータを取得
        oldest_users_editedLog = el.get_oldest_editedLog_by_userId(db,user_id)
        if oldest_users_editedLog is not None:
            # 今日の日付 - 最も古いデータの日付 = データ収集期間
            delta: timedelta = abs(clusters[0].date - oldest_users_editedLog.date + timedelta(days=1))
            # 差分を週単位に変換
            days_difference = math.floor(delta.days/7)
            # ここでクラスタリングの結果を元に確率を計算する(bがTrueなら以前, Falseなら以降)
            pr: float = nd.probability_from_normal_distribution(clusters, time, days_difference, b)
            result = ProbabilityResponse(userId=user_id, userName=user.name, probability=pr)
            result_json = jsonable_encoder(result)
            return ORJSONResponse(result_json)
        elif oldest_users_editedLog is None:
            none_result = ProbabilityResponse(userId=user_id, userName=user.name, probability=0)
            none_result_json = jsonable_encoder(none_result)
            return ORJSONResponse(none_result_json)
    elif user is None:
        raise HTTPException(status_code=404, detail="User not found")

# 全てのユーザがその日に入室する確率を返す
@router.get("/app/probability/{community}/all", response_class=ORJSONResponse, response_model=List[ProbabilityResponse])
async def get_probability_all(community:int, date:str = "2024-1-1", db: Session = Depends(get_db)):
    date_object= datetime.strptime(date, '%Y-%m-%d')
    seven_days_ago= date_object - timedelta(days=7)
    users = us.get_all_users_by_community(community,db)
    # 結果格納用のリスト
    result: list[ProbabilityResponse] = []
    # ユーザーごとに繰り返す
    for user in users:
        clusters = cl.get_all_cluster_by_userId_and_date(db, user.id, seven_days_ago, True)
        oldest_cluster = cl.get_oldest_cluster_by_userId(db, user.id, True)
        if oldest_cluster is not None:
            delta = abs(clusters[0].date - oldest_cluster.date + timedelta(days=1))
            # 差分を日単位に変換
            days_difference = math.floor(delta.days/7)
            # ここでクラスタリングの結果を元に確率を計算する(bがTrueなら以前, Falseなら以降)
            pr = nd.probability_from_normal_distribution(clusters, "24:00:00", days_difference, True)
            result.append(ProbabilityResponse(userId=user.id, userName=user.name, probability=pr))
        elif oldest_cluster is None:
            result.append(ProbabilityResponse(userId=user.id, userName=user.name, probability=0))
    # resultをjsonに変換
    result_json = jsonable_encoder(result)
    return ORJSONResponse(result_json)
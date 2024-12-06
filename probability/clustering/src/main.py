from datetime import date, timedelta
from lib.mysql import get_db
from models import cluster as cl, user as us, log, editedlog as el,

# 前提
# 毎日朝に実行される
# 手順
# 1. 前日の日付を取得、日付から曜日を取得
# 2. 前日と同じ曜日の入退室ログデータを取得(実行日が火曜日なら、これまでの月曜日のデータを取得)
# 3. 入退室ログを来訪・帰宅ログに変換(日付ごとの最初の入室と最後の退室を取得)
# 4. 入室と退室のそれぞれでクラスタリングを実行
# 5. クラスタごとの平均・分散を取得、保存
# 6. ユーザーごとに処理を繰り返す
# 補足
# 1. 入退室ログデータがない場合はスキップ
# 2. クラスタリングの結果がない場合は要素数1のクラスタが一つ(クラスタリングの結果がない = 来訪した日が1日だけの場合)
# 3. クラスタリングの結果がない場合は平均 = 値、分散 = 0


def main():
    db = get_db().__next__()
    # コミュニティに所属する全てのユーザを取得
    users = us.get_all_users_by_community(1, db)

    # 1. 前日の日付を取得、日付から曜日を取得
    # 今日と前日の日付を取得
    today = date.today()
    yesterday = today - timedelta(days=1)
    # 前日の曜日を取得
    day_of_week = yesterday.weekday()

    # 6. ユーザーごとに処理を繰り返す
    for user in users:
        user_id = user.id
        # 2. 前日と同じ曜日の入退室ログデータを取得
        # 入退室ログデータを取得
        df = log.get_log_by_userId_and_period
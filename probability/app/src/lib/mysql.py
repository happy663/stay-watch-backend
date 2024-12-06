from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import os

# 接続したいDBへの接続情報
user_name = os.environ['USER_NAME']
password = os.environ['PASSWORD']
host = os.environ['HOST']
port = os.environ['PORT']
database = os.environ['DATABASE']

SQLALCHEMY_DATABASE_URL = "mysql://" + user_name + ":" + password + "@" + host + ":" +port + "/" + database + "?charset=utf8&unix_socket=/var/run/mysqld/mysqld.sock"

engine = create_engine(SQLALCHEMY_DATABASE_URL, echo=True)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def get_db() :
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
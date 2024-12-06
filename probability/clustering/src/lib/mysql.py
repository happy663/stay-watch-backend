from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import os

# 接続したいDBへの接続情報
user_name = os.getenv('USER_NAME')
password = os.getenv('PASSWORD')
host = os.getenv('HOST')
port = os.getenv('PORT')
database = os.getenv('DATABASE')

SQLALCHEMY_DATABASE_URL = "mysql://" + user_name + ":" + password + "@" + host + ":" +port + "/" + database + "?charset=utf8&unix_socket=/var/run/mysqld/mysqld.sock"
print(SQLALCHEMY_DATABASE_URL)

engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def get_db() :
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
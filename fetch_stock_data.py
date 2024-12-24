# fetch_stock_data.py

import yfinance as yf
import pandas as pd


# pandas.DataFrame を戻す
def fetch_stock_data(symbol, start_date, end_date):
    # 日足データの取得
    daily_data = yf.download(symbol, start=start_date, end=end_date, interval="1d")
    return daily_data  # pandas.DataFrame を戻す


def main():
    symbol = "7203.T"  # トヨタ自動車の例
    start_date = "2023-01-01"  # 開始日
    end_date = "2023-12-31"  # 終了日

    # データの取得
    data = fetch_stock_data(symbol, start_date, end_date)

    # CSVファイルに出力
    output_file = f"{symbol}_stock_data.csv"
    data.to_csv(output_file)

    print(f"データをCSVファイルに保存しました: {output_file}")


if __name__ == "__main__":
    main()

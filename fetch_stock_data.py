# fetch_stock_data.py

import yfinance as yf
import pandas as pd
import os


# pandas.DataFrame を戻す
def fetch_stock_data(symbol, start_date, end_date):
    # 日足データの取得
    daily_data = yf.download(symbol, start=start_date, end=end_date, interval="1d")
    return daily_data  # pandas.DataFrame を戻す


#     """
#                   Open    High     Low   Close    Adj Close    Volume
#     Date
#     2023-12-01  2819.0  2842.0  2803.0  2833.0  2758.835693  26774000
#     2023-12-04  2802.0  2802.5  2744.5  2767.5  2695.050293  30495700
#     2023-12-05  2770.0  2784.5  2743.5  2753.5  2681.416748  24512600
#     """


def main():
    symbols = [
        "7203.T",  # Toyota Motor Corporation
        "7201.T",  # Nissan Motor Co., Ltd.
        "7267.T",  # Honda Motor Co., Ltd.
        "7261.T",  # Mazda Motor Corporation
        "7269.T",  # Suzuki Motor Corporation
        # "7262.T",  # Mitsubishi Motors Corporation 上場廃止
        "7270.T",  # Subaru Corporation
        "7202.T",  # Isuzu Motors Limited
        "7205.T",  # Hino Motors, Ltd.
        "7211.T",  # Mitsubishi Fuso Truck and Bus Corporation
        "7224.T",  # Shizuoka Daihatsu Motor Co., Ltd.
        "7266.T",  # Showa Corporation
        ]  # シンボルリストの例
    start_date = "2003-01-01"  # 開始日
    end_date = "2023-12-31"  # 終了日

    # CSVファイルを保存するフォルダの作成
    output_dir = "csv"
    os.makedirs(output_dir, exist_ok=True)

    for symbol in symbols:
        # データの取得
        data = fetch_stock_data(symbol, start_date, end_date)

        # CSVファイルに出力
        output_file = os.path.join(output_dir, f"{symbol}_stock_data.csv")
        data.to_csv(output_file, index_label="Date")
        print(f"データをCSVファイルに保存しました: {output_file}")


if __name__ == "__main__":
    main()

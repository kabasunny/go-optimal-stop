# run_commands.ps1

# コマンドを20回繰り返して実行するスクリプト
for ($i = 1; $i -le 20; $i++) {
    Write-Host "Run number: $i"
    go run ./cmd --random
    Start-Sleep -Seconds 1  # 1秒の待機時間を追加（必要に応じて調整）
}

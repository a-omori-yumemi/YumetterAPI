---
marp: true
---

# 6000 requests per second くらいのサービスを作る
大森　章裕 

---
# 自己紹介
## 大森　章裕  
東京工業大学 学部3年生
GitHub: https://github.com/Mojashi
Twitter: @oreha_senpai

- :heart: CTF(web)（最近激アツ）:heart:
- 競技プログラミング
- ボードゲーム
- ボルダリング（最近始めた）
- その他パソコンカタカタ等
---
# もともとできたこと
## webっぽいやつ
- Golang
- React.js
- MySQL
- データベース技術、トランザクションあたりに興味
## それ以外
- アルゴリズム, キカイガクシュー, 言語処理系 etc
---
# いまできること
## webっぽいやつ
- Golang
- React.js
- MySQL
- データベース技術、トランザクションあたりに興味
- <span style="color:red;font-weight:bold;">AWSチョット触れる</span>
- <span style="color:red;font-weight:bold;">負荷試験チョットできる<span>
- <span style="color:red;font-weight:bold;">その他設計や高速化に関する知見など・・・</span>

---
#  Agenda
1. YumetterAPI v1 の仕様を紹介
2. 設計を眺める
3. 負荷試験
4. 改善
5. 所感

    ---
大体作業時系列を追いかける形
随所に得た知見を散りばめていく

---
# Yumetterとは！！

---
![](https://i.gyazo.com/209a3a2ae58fd2ba238e8d022c1170d4.png)

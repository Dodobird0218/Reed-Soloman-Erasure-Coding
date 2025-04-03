# Reed-Solomon Erasure Coding

此專案為Erasure Coding的小實作，有兩組程式碼，分別以Lagrange和Vandermonde的方法來實現。

## 一、專案簡介

### 1. 有限域 GF(2^8) 實現 (@field.go)

有限域 GF(2^8) 是 Reed-Solomon 編碼的數學基礎，在Lagrange和Vandermonde方法中皆使用相同之規則

主要實現：

- 使用指數表和對數表加速運算
- 基本運算：加法、減法（同為異或運算）、乘法、除法
- 冪運算和乘法逆元計算
- 多項式 x^4 + x^3 + x^2 + 1 (0x1D) 作為本原多項式

核心方法包括：
- `Add/Sub`：有限域加/減法（XOR運算）
- `Mul`：利用對數表實現快速乘法
- `Div`：利用對數表實現快速除法
- `Inv`：計算乘法逆元
- `generateTables`：生成指數表和對數表

### 2. 編碼器和解碼器實現

#### Lagrange Encoder (@encoder.go)
- `RSEncoder2` 使用連續整數作為評估點
- 使用 Lagrange 插值多項式計算冗餘數據
- 提供高效編碼方法 `EncodeEfficient` 使用霍納法則
- 實現數據重建功能 `ReconstructData`

核心算法：
- 評估點生成：使用連續整數（1, 2, 3...）
- Lagrange 插值：計算多項式在特定評估點的值

#### Lagrange decoder (@decoder.go)
- `RSDecoder` 實現從任意足夠數量的分片中恢復原始數據
- 使用與編碼器相同的評估點
- 利用 Lagrange 插值公式重建原始數據

核心方法：
- `Decode`：從任意位置的分片恢復數據
- `DecodeLastShards`：從最後幾個分片恢復數據

#### Vandermonde Encoder (@encoder.go)
- `RSEncoder` 使用 Vandermonde 矩陣進行編碼
- 生成評估點和 Vandermonde 矩陣以計算冗餘數據
- 提供編碼方法 `Encode`，使用 Vandermonde 矩陣計算冗餘分片

核心算法：
- 評估點生成：使用連續整數（1, 2, 3...）
- Vandermonde 矩陣生成：計算每個冗餘分片的多項式值
- 使用 Vandermonde 矩陣進行編碼，計算冗餘數據

#### Vandermonde Decoder (@decoder.go)
- `VandermondeDecoder` 實現從任意足夠數量的分片中恢復原始數據
- 使用與編碼器相同的評估點
- 利用 Lagrange 插值公式重建原始數據

核心方法：
- `Decode`：從可用的分片中恢復原始數據
- `DecodeLastShards`：從最後幾個分片恢復數據

### 3. 主程式邏輯

#### 編碼主程式 (@encode_main.go)
- 讀取輸入 JSON 檔案獲取原始訊息
- 初始化有限域和 Reed-Solomon 編碼器
- 對原始訊息進行編碼
- 將編碼結果輸出到 JSON 檔案

#### 解碼主程式 (@decode_main.go)
- 讀取包含編碼分片的 JSON 檔案
- 初始化 Reed-Solomon 解碼器
- 從分片恢復原始訊息
- 將解碼結果保存到 JSON 檔案

### 4. CMD 執行範例

編碼範例：
```
./encode message.json encoded.json
```

解碼範例：
```
./decode encoded.json decoded.json
```

執行結果會顯示：
- 原始訊息和對應的十六進制表示
- 編碼/解碼結果及其十六進制表示
- 確認操作成功的資訊

### 5. 編碼結果/解碼結果介紹 (@data)

- `data/encode_01/`: 第一組編碼測試資料集
- `data/encode_02/`: 第二組編碼測試資料集
- `data/decode_01/`: 第一組解碼測試資料集
- `data/decode_02/`: 第二組解碼測試資料集

  - `data/encode/message.json`: 原始輸入訊信息
  - `data/encode/encoded.json`: 編碼結果
  - `data/decode/middle_message.json`: 編碼結果的中間六個shards
  - `data/decode/middle_decoded.json`: 中間六個shards的最終解碼結果
  - `data/decode/tail_message.json`: 編碼結果的最後六個shards
  - `data/decode/tail_decoded.json`: 最後六個shards的最終解碼結果

## 二、遇到困難

### 1. 兩種方法編碼結果不相同

第一次將兩個編碼器製作出來時編出來的碼是不一樣的，後來發現問題在於評估點的選擇，將評估點修改為一致即可解決問題。

### 2. 解碼不順利

我遇到的另外一個狀況是解碼的時候再解最後六個shards時能夠順利還原，但在解中間六個shards時無法順利還原。
後來發現問題在於需要知道是第幾個index開始的encode message才能解碼，原本的程式預設都是最後六個，所以中間的shards無法還原。
因此在解碼器的input json檔案中增加了`start_index`來判斷，最終解決了此問題。

## 三、心得

這是我第一個使用go語言撰寫的專案，讓我熟悉了go語言的資料結構與程式架構。
Erasure Coding 是一個很神奇的加密方法，竟然可以使用片段資料就能還原原始資料！
在一開始的時候我因為太急著想做出來，所以為了求快完成得很草率，帶來的結果就是功能不符預期。
之後我重新想了架構，重新想了EC背後的數學邏輯（尤其是評估點和有限域的設計），就順利做出來了！
當然做完的時候滿有成就感的，也藉由大型語言模型的幫助，讓我順利了許多。
我認為AI是很好用的工具，可以大幅提升我們的工作效率，不過整個架構設計、底層邏輯還是要靠自己，所以基本功還是很重要的。
使用不熟悉的程式語言製作，花費幾個小時完成了這個專案，我認為自己這次的表現算是及格了，對自己滿滿意的！

結果這次專案花滿多時間在打報告的部分XD

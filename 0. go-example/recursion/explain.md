# Giải thích mã nguồn `recursion.go`

## Tổng quan

File `recursion.go` là một chương trình Go được thiết kế để so sánh hai struct bất kỳ, phát hiện sự khác biệt về giá trị, kiểu dữ liệu, hoặc độ dài của slice/mảng. Nó hỗ trợ các tính năng nâng cao như:

- **So sánh đồng thời (concurrent)** : Sử dụng goroutine và semaphore để xử lý so sánh nhiều trường đồng thời, tối ưu hiệu suất.
- **Quy tắc so sánh tùy chỉnh** : Cho phép định nghĩa hàm so sánh riêng cho các trường hoặc kiểu cụ thể.
- **Giới hạn độ sâu đệ quy** : Ngăn chặn tràn stack khi so sánh các struct lồng nhau sâu.
- **Hỗ trợ struct tag** : Cho phép đánh dấu các trường để bỏ qua hoặc áp dụng so sánh tùy chỉnh thông qua tag.
- **Ghi log** : Hỗ trợ ghi lại các sự kiện trong quá trình so sánh để debug.
- **Tái sử dụng** : Có thể reset trạng thái để sử dụng lại mà không cần khởi tạo mới.

Mã nguồn này được viết theo phong cách của một lập trình viên cấp cao (senior), tập trung vào tính sạch (clean), rõ ràng (clarify), dễ bảo trì (maintainable), và khả năng mở rộng (scalable).

---

## Mục đích và cách sử dụng

### Mục đích

`StructComparer` được thiết kế để so sánh hai giá trị struct trong Go, trả về danh sách các khác biệt (`ComparisonResult`) về giá trị, kiểu, hoặc độ dài của slice/mảng. Nó hữu ích trong các trường hợp như:

- **Kiểm thử (testing)** : So sánh dữ liệu đầu ra với dữ liệu mong đợi.
- **Debugging** : Xác định sự khác biệt giữa hai phiên bản của một struct.
- **Đồng bộ hóa dữ liệu** : Phát hiện thay đổi giữa hai phiên bản của một đối tượng.

### Cách sử dụng

1. **Khởi tạo StructComparer** :

- Tạo một instance của `StructComparer` với cấu hình (`StructComparerConfig`), bao gồm số lượng goroutine tối đa, trường bỏ qua, thời gian timeout, v.v.
- Ví dụ:
  ```go
  comparer := NewStructComparer(StructComparerConfig{
      MaxGoroutines:    5,
      IgnoreFields:     []string{"Age"},
      CompareTimeout:   5 * time.Second,
      MaxDepth:         10,
      CustomComparers:  map[string]CustomCompareFunc{"case_insensitive": caseInsensitiveCompare},
      EnableLogging:    true,
  })
  ```

1. **So sánh struct** :

- Gọi phương thức `Compare` với hai giá trị struct và một context.
- Nhận danh sách `ComparisonResult` chứa các khác biệt.
- Ví dụ:
  ```go
  results, err := comparer.Compare(context.Background(), p1, p2)
  ```

1. **In kết quả** :

- Lặp qua `results` để in các khác biệt (đường dẫn, loại khác biệt, giá trị, kiểu trường).
- Ví dụ:
  ```go
  for _, result := range results {
      fmt.Printf("Difference at %s (%s, Type: %v): %v != %v\n",
          result.Path, result.DiffType, result.FieldType, result.Value1, result.Value2)
  }
  ```

1. **Tái sử dụng** :

- Gọi `Reset` để xóa cache và kết quả, cho phép sử dụng lại `StructComparer`.
- Ví dụ:
  ```go
  comparer.Reset()
  ```

---

## Các thành phần chính

### 1. Structs và Types

- **`Address` và `Person`** :
- `Address`: Biểu diễn một địa chỉ với các trường `City` và `Zip`. Có tag `comparer:"ignore"` để bỏ qua trường `City`.
- `Person`: Biểu diễn một người với các trường `Name`, `Age`, `Address`, và `Hobbies`. Có tag `comparer:"custom=case_insensitive"` cho `Name` và `comparer:"ignore"` cho `Age`.
- **`ComparisonResult`** :
- Lưu trữ thông tin về một khác biệt:
  - `Path`: Đường dẫn đến trường (e.g., `root.Address.Zip`).
  - `DiffType`: Loại khác biệt (`TypeMismatch`, `LengthMismatch`, `ValueMismatch`, `CustomMismatch`).
  - `Value1`, `Value2`: Giá trị của hai struct tại trường đó.
  - `FieldType`: Kiểu của trường (thêm để cung cấp ngữ cảnh).
  - `Timestamp`: Thời điểm phát hiện khác biệt.
- **`CustomCompareFunc`** :
- Hàm tùy chỉnh để so sánh hai giá trị (trả về `true` nếu bằng nhau, `false` nếu khác).
- **`StructComparerConfig`** :
- Cấu hình cho `StructComparer`, bao gồm:
  - `MaxGoroutines`: Số lượng goroutine tối đa.
  - `IgnoreFields`: Danh sách trường bỏ qua.
  - `CompareTimeout`: Thời gian timeout cho so sánh.
  - `MaxDepth`: Giới hạn độ sâu đệ quy.
  - `CustomComparers`: Bản đồ các hàm so sánh tùy chỉnh.
  - `EnableLogging`, `Logger`: Hỗ trợ ghi log.
  - `IgnoreTag`, `CustomCompareTag`: Tên tag để bỏ qua hoặc so sánh tùy chỉnh.
- **`StructComparer`** :
- Quản lý quá trình so sánh, với các trường:
  - `config`: Lưu cấu hình.
  - `ignoreFields`: Bản đồ các trường bỏ qua (O(1) lookup).
  - `customFields`: Bản đồ các hàm so sánh tùy chỉnh.
  - `results`: Danh sách các khác biệt.
  - `cache`: Cache kết quả so sánh để tránh lặp lại.
  - `semaphore`: Kênh để giới hạn số goroutine.
  - `mu`: Khóa để đảm bảo an toàn đồng thời.
  - `logger`: Logger tùy chọn.

### 2. Hàm và phương thức chính

- **`NewStructComparer`** :
- Khởi tạo `StructComparer` với cấu hình, thiết lập giá trị mặc định và chuyển đổi `IgnoreFields` thành map.
- **`Compare`** :
- Phương thức công khai để so sánh hai struct, trả về danh sách `ComparisonResult`.
- **`compareStructs`** :
- So sánh hai giá trị đệ quy, kiểm tra context, độ sâu, cache, và xử lý các kiểu khác nhau (struct, slice/array, hoặc nguyên thủy).
- **`compareStructFields`** :
- So sánh các trường của struct, kiểm tra tag và cấu hình để bỏ qua hoặc áp dụng so sánh tùy chỉnh, sử dụng goroutine cho mỗi trường.
- **`compareArrayElements`** :
- So sánh các phần tử của slice/mảng, sử dụng goroutine cho mỗi phần tử.
- **`comparePrimitiveValues`** :
- So sánh các giá trị nguyên thủy (int, string, v.v.) bằng `reflect.DeepEqual`.
- **`addResult`** :
- Ghi lại khác biệt một cách an toàn đồng thời.
- **`GetResults`** :
- Trả về bản sao của danh sách khác biệt.
- **`Reset`** :
- Xóa cache và kết quả để tái sử dụng.
- **`log`** :
- Ghi log nếu được bật, hỗ trợ debug.

### 3. Cơ chế đệ quy

Có, mã nguồn sử dụng **đệ quy** để so sánh các struct lồng nhau. Cách hoạt động như sau:

- **Điểm bắt đầu** : Phương thức `Compare` gọi `compareStructs` với hai giá trị `reflect.Value` và đường dẫn `root`.
- **Xử lý struct** :
- Nếu giá trị là struct, `compareStructFields` lặp qua các trường và gọi `compareStructs` đệ quy cho mỗi trường không bị bỏ qua.
- Mỗi lời gọi đệ quy tăng `depth` để theo dõi độ sâu.
- **Xử lý slice/mảng** :
- Nếu giá trị là slice/mảng, `compareArrayElements` lặp qua các phần tử và gọi `compareStructs` đệ quy cho mỗi phần tử.
- **Điều kiện dừng** :
- Khi gặp giá trị nguyên thủy, gọi `comparePrimitiveValues`.
- Khi vượt quá `MaxDepth`, trả về lỗi.
- Khi context hết hạn hoặc giá trị được cache, dừng đệ quy.

**Giới hạn độ sâu (`MaxDepth`)** đảm bảo không bị tràn stack khi xử lý các struct lồng nhau sâu (ví dụ: struct chứa struct chứa struct, v.v.).

### 4. Đồng thời (Concurrency)

- **Goroutine** : Mỗi trường của struct hoặc phần tử của slice/mảng được so sánh trong một goroutine riêng, cho phép xử lý song song.
- **Semaphore** : Giới hạn số lượng goroutine đồng thời bằng `MaxGoroutines`, sử dụng kênh `semaphore`.
- **An toàn đồng thời** :
- `sync.RWMutex` bảo vệ việc ghi/đọc `results`.
- `sync.Map` được sử dụng cho cache để tránh xung đột.

### 5. Các tính năng nâng cao

- **Quy tắc so sánh tùy chỉnh** :
- Cho phép định nghĩa hàm so sánh (e.g., so sánh không phân biệt hoa thường cho chuỗi).
- Ví dụ: `caseInsensitiveCompare` so sánh `Name` của `Person` không phân biệt hoa/thường.
- **Struct Tag** :
- Hỗ trợ tag `comparer:"ignore"` để bỏ qua trường (e.g., `City` trong `Address`).
- Hỗ trợ tag `comparer:"custom=case_insensitive"` để áp dụng so sánh tùy chỉnh (e.g., `Name` trong `Person`).
- **Ghi log** :
- Ghi lại các sự kiện như cache hit, trường bị bỏ qua, hoặc khác biệt được phát hiện.
- **Reset** :
- Xóa cache và kết quả để tái sử dụng `StructComparer`.

---

## Ý nghĩa của mã nguồn

Mã nguồn này là một công cụ mạnh mẽ và linh hoạt để so sánh struct trong Go, phù hợp với các ứng dụng yêu cầu kiểm tra sự khác biệt giữa các đối tượng phức tạp. Nó:

- **Hiệu quả** : Sử dụng đồng thời và cache để tối ưu hiệu suất.
- **Linh hoạt** : Hỗ trợ so sánh tùy chỉnh và bỏ qua trường thông qua cấu hình hoặc tag.
- **An toàn** : Đảm bảo an toàn đồng thời và giới hạn độ sâu đệ quy.
- **Dễ bảo trì** : Mã sạch, có comment rõ ràng, và cấu trúc modular.
- **Mở rộng** : Dễ dàng thêm các tính năng mới như quy tắc so sánh hoặc hỗ trợ kiểu dữ liệu khác.

---

## Ví dụ thực thi

Trong hàm `main`, mã nguồn so sánh hai struct `Person`:

```go
p1 := Person{
    Name: "Alice",
    Age:  30,
    Address: Address{
        City: "New York",
        Zip:  "10001",
    },
    Hobbies: []string{"reading", "hiking"},
}

p2 := Person{
    Name: "ALICE",
    Age:  31,
    Address: Address{
        City: "New York",
        Zip:  "10002",
    },
    Hobbies: []string{"reading", "swimming"},
}
```

**Kết quả** :

- `Age` bị bỏ qua (do cấu hình).
- `City` bị bỏ qua (do tag `comparer:"ignore"`).
- `Name` được so sánh không phân biệt hoa/thường, nên "Alice" và "ALICE" được coi là bằng nhau.
- Phát hiện khác biệt ở `Zip` (`10001` vs `10002`) và `Hobbies` (`hiking` vs `swimming`).
- Sau khi gọi `Reset`, số lượng kết quả về 0.

  **Output mẫu** :

```
Difference at root.Address.Zip (ValueMismatch, Type: string): 10001 != 10002
Difference at root.Hobbies[1] (ValueMismatch, Type: string): hiking != swimming
After reset, results count: 0
```

---

## Kết luận

`recursion.go` là một công cụ mạnh mẽ, được thiết kế tốt để so sánh struct trong Go. Nó sử dụng đệ quy để xử lý các struct lồng nhau, đồng thời để tối ưu hiệu suất, và cung cấp các tính năng nâng cao như tag, so sánh tùy chỉnh, và ghi log. Mã nguồn này phù hợp cho các ứng dụng phức tạp, dễ dàng mở rộng, và tuân thủ các tiêu chuẩn của lập trình viên cấp cao.

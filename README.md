# Layer Edge Light Node

## Giới thiệu

Layer Edge Light Node là một client kết nối với mạng Layer Edge để xác minh cây Merkle bằng cách thu thập các mẫu ngẫu nhiên từ các cây có sẵn và xác minh tính toàn vẹn của chúng. Light node thực hiện các hoạt động xác minh bằng chứng Zero-Knowledge thông qua dịch vụ ZK prover cục bộ và gửi các bằng chứng đã xác minh lên mạng.

Light node hoạt động bằng cách:
1. Khám phá các cây Merkle có sẵn từ mạng Layer Edge
2. Thu thập các mẫu ngẫu nhiên từ các cây để xác minh
3. Tạo và xác minh các bằng chứng Zero-Knowledge
4. Gửi các bằng chứng đã xác minh để nhận phần thưởng

## Tính năng chính

- Tự động khám phá các cây Merkle có sẵn từ mạng
- Thu thập các mẫu ngẫu nhiên từ cây để xác minh
- Tạo và xác minh các bằng chứng Zero-Knowledge
- Triển khai cơ chế ngủ thông minh để tránh công việc dư thừa trên các cây không thay đổi
- Gửi các bằng chứng đã xác minh để nhận phần thưởng
- Hỗ trợ proxy để tăng cường khả năng kết nối và bảo mật
- Quản lý ví tích hợp để xác thực và nhận phần thưởng

## Cấu trúc dự án

```
light-node-dev/
├── clients/             # Các client giao tiếp với mạng Layer Edge
│   ├── cosmos.go        # Client Cosmos cho giao tiếp blockchain
│   └── request.go       # Xử lý các yêu cầu HTTP/gRPC
├── node/                # Lõi của light node
│   └── verifier.go      # Logic xác minh cây Merkle
├── risc0-merkle-service/# Dịch vụ ZK prover dựa trên RISC Zero
│   ├── cli/             # Giao diện dòng lệnh cho dịch vụ
│   ├── host/            # Máy chủ dịch vụ
│   └── methods/         # Các phương thức ZK
├── scripts/             # Scripts hỗ trợ xây dựng và chạy
├── utils/               # Các tiện ích
│   ├── get_env.go       # Xử lý biến môi trường
│   ├── hash_string.go   # Hàm băm chuỗi
│   ├── proxy_ops.go     # Hoạt động proxy
│   ├── random_sampler.go# Lấy mẫu ngẫu nhiên
│   └── wallet_ops.go    # Hoạt động ví
├── main.go              # Điểm vào chính của ứng dụng
├── go.mod               # Quản lý phụ thuộc Go
└── Makefile             # Tệp Makefile để xây dựng và chạy
```

## Yêu cầu hệ thống

- Go 1.18 trở lên
- Rust 1.81.0 trở lên
- Truy cập vào điểm cuối gRPC của Layer Edge
- Bộ công cụ RISC Zero

### Cài đặt bộ công cụ RISC Zero

```bash
curl -L https://risczero.com/install | bash && rzup install
```

## Hướng dẫn cài đặt

### Cấu hình biến môi trường

Tạo tệp `.env` hoặc cấu hình biến môi trường:

```env
GRPC_URL=grpc.testnet.layeredge.io:9090
CONTRACT_ADDR=cosmos1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqt56709
ZK_PROVER_URL=http://127.0.0.1:3001
# Hoặc sử dụng ZK Prover từ xa:
# ZK_PROVER_URL=https://layeredge.mintair.xyz/
API_REQUEST_TIMEOUT=100
POINTS_API=https://light-node.layeredge.io
PRIVATE_KEY='cli-node-private-key'
```

Đảm bảo URL ZK Prover giống với URL của máy chủ nơi dịch vụ merkle đang chạy, hoặc sử dụng URL từ xa nếu bạn không muốn chạy dịch vụ cục bộ.

### Cấu hình proxy (tùy chọn)

Nếu bạn muốn sử dụng proxy, hãy tạo tệp `proxy.txt` với danh sách các proxy theo định dạng:

```
http://username:password@host:port
http://username:password@host:port
...
```

### Cấu hình ví

Tạo tệp `wallet.txt` với khóa riêng tư của ví:

```
your-private-key-here
```

## Xây dựng và chạy

### Sử dụng scripts

Chúng tôi cung cấp các scripts để đơn giản hóa quá trình xây dựng và chạy:

```bash
# Xây dựng dịch vụ RISC Zero
./scripts/build-risczero.sh

# Xây dựng light node
./scripts/build-light-node.sh

# Chạy cả hai dịch vụ
./scripts/runner.sh
```

### Chạy thủ công

#### 1. Khởi động dịch vụ RISC Zero Merkle

```bash
cd risc0-merkle-service
cargo build && cargo run
```

#### 2. Khởi động Light Node

Trong một terminal khác:

```bash
go build
./light-node
```

Đảm bảo cả hai dịch vụ đều chạy độc lập.

## Hướng dẫn cài đặt và chạy trên Windows

### Cài đặt các công cụ cần thiết

1. **Cài đặt Go**:
   - Tải Go từ [trang chủ Go](https://golang.org/dl/) cho Windows
   - Chạy trình cài đặt và làm theo hướng dẫn
   - Kiểm tra cài đặt bằng cách mở Command Prompt và gõ: `go version`

2. **Cài đặt Rust**:
   - Tải và chạy [rustup-init.exe](https://www.rust-lang.org/tools/install)
   - Chọn cài đặt mặc định
   - Kiểm tra cài đặt bằng cách mở Command Prompt mới và gõ: `rustc --version`

3. **Cài đặt Git**:
   - Tải Git từ [trang chủ Git](https://git-scm.com/download/win)
   - Chạy trình cài đặt và làm theo hướng dẫn
   - Kiểm tra cài đặt bằng cách mở Command Prompt và gõ: `git --version`

4. **Cài đặt bộ công cụ RISC Zero**:
   - Mở PowerShell với quyền Administrator
   - Chạy lệnh sau để cài đặt:
     ```powershell
     iwr -useb https://risczero.com/install.ps1 | iex
     rzup install
     ```

### Cấu hình biến môi trường

Tạo tệp `.env` trong thư mục gốc của dự án với nội dung tương tự như đã mô tả ở trên.

Hoặc cấu hình biến môi trường Windows thông qua PowerShell:

```powershell
$env:GRPC_URL="grpc.testnet.layeredge.io:9090"
$env:CONTRACT_ADDR="cosmos1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqt56709"
$env:ZK_PROVER_URL="http://127.0.0.1:3001"
$env:API_REQUEST_TIMEOUT="100"
$env:POINTS_API="https://light-node.layeredge.io"
$env:PRIVATE_KEY="cli-node-private-key"
```

### Xây dựng và chạy trên Windows

#### Sử dụng PowerShell

1. **Clone repository**:
   ```powershell
   git clone https://github.com/hieusats/light-node-dev.git
   cd light-node-dev
   ```

2. **Xây dựng và chạy dịch vụ RISC Zero Merkle**:
   ```powershell
   cd risc0-merkle-service
   cargo build
   cargo run
   ```

3. **Mở một PowerShell mới và xây dựng Light Node**:
   ```powershell
   cd path\to\light-node-dev
   go build
   .\light-node.exe
   ```

#### Sử dụng batch scripts (tùy chọn)

Bạn có thể tạo các batch scripts để đơn giản hóa quá trình:

1. Tạo file `build-risczero.bat`:
   ```batch
   @echo off
   cd risc0-merkle-service
   cargo build
   echo RISC Zero Merkle service built successfully
   ```

2. Tạo file `build-light-node.bat`:
   ```batch
   @echo off
   go build
   echo Light Node built successfully
   ```

3. Tạo file `run-all.bat`:
   ```batch
   @echo off
   start cmd /k "cd risc0-merkle-service && cargo run"
   timeout /t 5
   start cmd /k ".\light-node.exe"
   echo Services started
   ```

### Lưu ý khi chạy trên Windows

1. **Đường dẫn tệp**: Windows sử dụng dấu gạch ngược (`\`) thay vì dấu gạch chéo (`/`) cho đường dẫn tệp. Đảm bảo điều chỉnh đường dẫn phù hợp.

2. **Tường lửa Windows**: Bạn có thể cần cho phép ứng dụng qua tường lửa Windows khi được nhắc.

3. **Quyền Administrator**: Một số thao tác có thể yêu cầu quyền Administrator. Chạy Command Prompt hoặc PowerShell với quyền Administrator nếu cần.

4. **Proxy trên Windows**: Nếu sử dụng proxy, đảm bảo định dạng trong tệp `proxy.txt` là chính xác và không có ký tự đặc biệt không mong muốn.

5. **Lỗi kết nối**: Nếu gặp lỗi kết nối, hãy kiểm tra cấu hình mạng và tường lửa Windows.

## Ghi nhật ký và giám sát

Light node cung cấp ghi nhật ký chi tiết về các hoạt động của nó. Bạn có thể theo dõi đầu ra nhật ký để theo dõi:

- Khám phá cây
- Tạo và xác minh bằng chứng
- Gửi các bằng chứng đã xác minh
- Trạng thái ngủ của cây
- Hoạt động proxy
- Giao dịch ví

## Khắc phục sự cố

Nếu bạn gặp vấn đề:

1. Kiểm tra kết nối gRPC của bạn với mạng Layer Edge
2. Đảm bảo dịch vụ ZK prover đang chạy và có thể truy cập được
3. Xác minh địa chỉ ví và định dạng chữ ký của bạn
4. Kiểm tra nhật ký để biết thông báo lỗi cụ thể
5. Xác minh cấu hình proxy nếu bạn đang sử dụng
6. Đảm bảo khóa riêng tư của ví hợp lệ

## Đóng góp

Chúng tôi hoan nghênh đóng góp! Vui lòng làm theo các bước sau:

1. Fork repository
2. Tạo nhánh tính năng (`git checkout -b feature/amazing-feature`)
3. Commit các thay đổi của bạn (`git commit -m 'Add some amazing feature'`)
4. Push lên nhánh (`git push origin feature/amazing-feature`)
5. Mở Pull Request

## Giấy phép

Dự án này được cấp phép theo Giấy phép MIT - xem tệp LICENSE để biết chi tiết.

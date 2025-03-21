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

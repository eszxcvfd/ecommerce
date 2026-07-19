# Research: Provider Boundaries for Payment, Payout, and File Storage (Issue #2)

**Status:** Research note — not a final ADR or production recommendation
**Date:** 2026-07-19 (refreshed)
**Context:** Issue #2 / MVP Issue #1
**Legal scope:** Intentionally excluded (see §5)

---

## 1. Payment / Top-up Provider Capabilities (First-Party Docs)

### 1.1 VNPay

**Developer portal:** `sandbox.vnpayment.vn/apis/`
**Product:** VNPAY-QR — redirect-based payment gateway with QR code, domestic ATM/card, and international card support.

**Documented capabilities:**

- **Create payment request/URL:** The merchant creates an order with the `pay` command and redirects the user to VNPay's payment page (`https://sandbox.vnpayment.vn/paymentv2/vpcpay.html` in sandbox). Parameters include `vnp_Amount` (in VND × 100, to eliminate decimals), `vnp_TxnRef` (unique per day), `vnp_ReturnUrl`, and `vnp_OrderInfo`. [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html]
- **IPN callback:** VNPay supports server-to-server IPN (Instant Payment Notification) at the merchant's registered `IPN URL` to notify the backend of payment results. The IPN is sent after the user completes payment on VNPay's interface. [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html — IPN URL section]
- **Return URL redirect:** After payment completion, VNPay redirects the user back to the merchant's `vnp_ReturnUrl`. The merchant should verify the signature on this return as well (not just the IPN). [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html]
- **Signature verification:** Uses HMAC-SHA512 with the merchant's `vnp_HashSecret` secret key. The signing algorithm sorts parameters alphabetically, builds a query string, and computes `hash_hmac('sha512', hashdata, vnp_HashSecret)`. [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html — SecureHash section]
- **Payment expiration:** Supports `vnp_ExpireDate` parameter to set the payment URL expiration time (format `yyyyMMddHHmmss`, GMT+7). [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html — vnp_ExpireDate]
- **Query transaction status (QueryDr):** The merchant can query transaction results via POST to `https://sandbox.vnpayment.vn/merchant_webapi/api/transaction` with `vnp_Command=querydr`. Returns detailed status including `vnp_TransactionStatus`, `vnp_BankCode`, `vnp_PayDate`, `vnp_TransactionNo`. [sandbox.vnpayment.vn/apis/docs/truy-van-hoan-tien/querydr&refund.html]
- **Refund:** Supports partial (`vnp_TransactionType=03`) and full (`vnp_TransactionType=02`) refunds via the same API with `vnp_Command=refund`. [sandbox.vnpayment.vn/apis/docs/truy-van-hoan-tien/querydr&refund.html]
- **Payment method selection:** The merchant can optionally specify `vnp_BankCode` to force a specific payment method (`VNPAYQR`, `VNBANK`, `INTCARD`). If omitted, the user chooses at the VNPay interface. [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html]
- **Sandbox:** Available at `sandbox.vnpayment.vn/`. Test credentials can be registered at `sandbox.vnpayment.vn/devreg/`. [sandbox.vnpayment.vn/apis/]
- **Currency:** VND only (`vnp_CurrCode=VND`). No multi-currency support documented. [sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html]
- **Prerequisites:** Requires `vnp_TmnCode` (merchant code) and `vnp_HashSecret` (secret key), obtained by registering at `sandbox.vnpayment.vn/devreg/`. Documentation says "If you do not have configuration information yet, you can register here" — implying self-service test credentials. Public documentation does not disclose production fees, transaction limits, or contract terms.

**Unverified from first-party docs:** Fee structure per transaction, daily/monthly production limits, exact processing times, production contract requirements — these require merchant registration.

### 1.2 MoMo (Payment)

**Developer portal:** `developers.momo.vn` (redirects to `developer.momo.vn`)
**GitHub:** `github.com/momo-wallet` — public SDK repositories (Java, PHP, JavaScript, iOS, Android).

**Documented capabilities from first-party sources:**

- **All-In-One (AIO) Payment Gateway:** Supports redirect-based online payment. The merchant sends a payment creation request and receives a `payUrl` for user redirect. [github.com/momo-wallet/payment — README; github.com/momo-wallet/java — README]
- **Payment methods:** Online Payment (Desktop, Mobile website), Offline (POS, Static QR, Dynamic QR), Mobile (App-to-App, In-MoMo-Application). [github.com/momo-wallet/java — README]
- **Callback / IPN:** MoMo sends server-to-server IPN callbacks to notify merchants of transaction results. The README references PayPal's IPN model. `ipnUrl` parameter is available in the disbursement API and presumably in the payment API. [github.com/momo-wallet/payment — README]
- **Signature verification:** Uses HMAC-SHA256 for request signing, RSA for encrypting sensitive fields, and AES for additional data protection. [github.com/momo-wallet/payment — README]
- **Query/status:** The AIO gateway supports transaction status queries via the Java SDK. [github.com/momo-wallet/java — README; github.com/momo-wallet/java — QueryStatusTransactionRequest/Response models]
- **Sandbox:** Two environments: `dev` (development/sandbox) and `prod` (production). Environment configuration is available in the Java SDK. [github.com/momo-wallet/java — README]
- **Two environments:** `dev` (sandbox) and `prod` (production) with different endpoint configurations. [github.com/momo-wallet/java — Environment.java]
- **Prerequisites:** Requires `partnerCode` and `accessKey` and a secret key for signing. Exact contract terms are not public.
- **Security algorithms:** HMAC 256, RSA, AES. [github.com/momo-wallet/payment — README]

**Unverified from first-party docs:** Fee details, transaction limits, settlement timelines, production service-level commitments — these require an active merchant account with MoMo.

### 1.3 ZaloPay (Payment)

**Developer portal:** `docs.zalopay.vn`
**GitHub:** `github.com/zalopay-samples` — quickstart repositories (Next.js, Node.js, mobile SDKs).

**Documented capabilities from first-party sources:**

- **Create payment order:** The merchant calls the `createOrder` API with `app_id`, `app_trans_id`, `app_user`, `app_time`, `amount`, `item`, `description`, `embed_data`, `callback_url`, and `mac` (HMAC signature). The API returns an `order_url` for redirect or `zp_trans_token` for app-to-app. [github.com/zalopay-samples/quickstart-payment-gateway]
- **Callback / Webhook:** ZaloPay sends a POST callback to the merchant's registered `callback_url` when a payment is successful. The callback contains `data` (JSON string), `mac` (HMAC-SHA256), and `type` (`1` = order, `2` = agreement). The merchant must respond with `return_code: 1` for success or `2` for failure. ZaloPay retries callbacks up to 3 times if no 1/2 response is received. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Signature verification for callback:** Computed as `HMAC-SHA256(key2, callback_data.data)`, compared with `callback_data.mac`. `key2` is provided by ZaloPay at merchant registration. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Callback data fields:** Include `app_trans_id`, `app_time`, `amount`, `zp_trans_id` (ZaloPay transaction code), `server_time`, `channel` (payment channel ID), `user_fee_amount`, `discount_amount`. The callback data JSON contains all fields as a string under the `data` key. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Query/status:** If no callback is received within 15 minutes of order creation, the merchant should proactively call the `QueryOrder` API to get the final result. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Idempotency requirement:** ZaloPay's callback documentation explicitly warns that callback endpoints "may receive duplicate events" and recommends "tracking the events that have been processed and avoiding reprocessing events that have already been tracked." [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Payment channels documented:** Includes channel ID numbers in the callback data. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Multi-product support:** Callback mechanism supports Order payments, Tokenization/Agreement, and ZOD (ZaloPay Open Domain) products with different callback data schemas. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Sandbox:** Available with test account credentials in `zalopay-samples/test-apps` and test wallet apps in `zalopay-samples/test-wallets`. Callback logs can be inspected via the Sandbox Merchant Portal at `sbmc.zalopay.vn/devtool`. [github.com/zalopay-samples/test-apps; github.com/zalopay-samples/test-wallets]
- **Code examples:** Callback handling code is provided in Node.js, Python, and Go in the official documentation. [docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/]
- **Prerequisites:** Requires `app_id`, `key1`, `key2` — provided by ZaloPay upon merchant registration. Documentation references `ZLP_MERCHANT_CALLBACK_URL` in environment config.

**Unverified from first-party docs:** Fee details, transaction limits, settlement schedules, production SLA.

---

## 2. Payout Capabilities (First-Party Docs)

### 2.1 MoMo (Disbursement API v2)

**Developer portal:** `developer.momo.vn/v3/docs/payment/api/disbursement-v2/`

**Documented capabilities (verified from docs page):**

- **Check Wallet Status:** POST `/v2/gateway/api/disbursement/verify` — checks whether a recipient wallet can receive money before disbursement. Validates MoMo account age (18+), NFC C06 verification status, and bank account linking (Circular No. 40 of the State Bank of Vietnam). Request fields: `partnerCode`, `orderId`, `requestId`, `requestType=checkWallet`, `disbursementMethod` (RSA-encrypted JSON with `walletId`, `walletName`, `personalId`), `signature` (HMAC-SHA256). [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Get Current Merchant Balance:** POST `/v2/gateway/api/disbursement/balance` — checks the merchant's remaining disbursement fund balance. Supports VND and global currencies (USD, EUR, AUD, etc.). Optional `orderGroupId` for distinguishing merchant balance groups. [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Single Disbursement (wallet):** POST `/v2/gateway/api/disbursement/pay` with `requestType=disburseToWallet`. Uses RSA-encrypted `disbursementMethod` JSON with `walletId` (phone number), `walletName`, `personalId`. Amount range: min 1,000 VND, max 200,000,000 VND per transaction. Requires minimum 30-second API timeout. [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Single Disbursement (bank):** POST `/v2/gateway/api/disbursement/pay` with `requestType=disburseToBank`. Supports both bank account number and bank card number. Fields: `bankAccountNo`/`bankCardNo`, `bankAccountHolderName`, `bankCode` (e.g., VCB, ACB, BIDV). Amount range: min 10,000 VND, max 20,000,000 VND per transaction. [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **IPN callback:** The `ipnUrl` parameter is available in the disbursement pay request. MoMo sends POST callbacks to this URL with payment results including `partnerCode`, `orderId`, `requestId`, `amount`, `transId`, `resultCode`, `message`, `responseTime`, `extraData`, `signature`. The callback uses HMAC-SHA256 for signature verification. [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Idempotency:** The `requestId` field is used for idempotency control — duplicate `requestId`s return result code 40 ("Duplicated requestId"). [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Minimum API timeout:** Minimum 30 seconds when calling the disbursement pay API to ensure response is received from MoMo. [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Result codes:** Extensive result code table for disbursement: 0 (success), 10 (server error), 20 (bad format), 22 (amount out of range), 40 (duplicate requestId), 42 (invalid orderId), 99 (unknown), 1007 (inactive account), 1008 (exceeds receiving limit), 1100 (insufficient merchant balance), 1507 (bank not found), 4001 (account restricted), 4002 (not verified by C06/NFC), 4003 (invalid receiver info), 7000/7002 (pending/processing). [developer.momo.vn/v3/docs/payment/api/disbursement-v2/]
- **Batch Disbursement:** The existing note cited batch endpoint `/v2/gateway/api/disbursement/batch/pay` with ≤1,000 items. **Not re-verified from the current docs page** — the single-disbursement page does not include batch endpoint documentation. A separate batch page may exist but was not found at the same URL. Verify at merchant onboarding.

**Corrections from previous version of this note:**
- **Bank disbursement IS supported** — the docs clearly define `disburseToBank` with `bankAccountNo`, `bankCardNo`, and `bankCode` fields. The earlier version incorrectly stated that only wallet disbursement was documented.
- **IPN callbacks ARE documented for disbursement** — the `ipnUrl` parameter and callback parameter table are present in the docs.
- **Amount limits ARE documented** — wallet: 1,000–200,000,000 VND; bank: 10,000–20,000,000 VND.

**Unverified:** Batch disbursement endpoint not re-confirmed on the current docs page. Exact fees per payout (not documented in public API docs). Production SLA and contract terms.

### 2.2 ZaloPay (Disbursement)

**Status as of July 2026:** The previously cited URL `docs.zalopay.vn/en/v2/payments/disbursement/overview.html` now returns **HTTP 404**. The disbursement documentation page has been removed or moved to a different URL. The quickstart repository at `github.com/zalopay-samples/quickstart-disbursement` still exists and confirms a disbursement feature.

**What is still confirmed:**
- ZaloPay offers disbursement/payout APIs ("merchant transfer money to user" including "payroll" scenarios). [github.com/zalopay-samples/quickstart-disbursement/README.md]
- Sandbox/test credentials are provided in a sample `.env` including `REACT_APP_APP_ID`, `REACT_APP_PAYMENT_ID`, `REACT_APP_KEY1`, `REACT_APP_PRIVATE_KEY`. [github.com/zalopay-samples/quickstart-disbursement/.env]
- The quickstart is a Node.js + React payroll application using Redux Toolkit. [github.com/zalopay-samples/quickstart-disbursement]

**Unverified:** Request/response structure, batch support, callback mechanism, fee structure, limits, wallet vs. bank account support — all unknown since the official docs page is no longer accessible at the previously known URL. The quickstart repo does not document API details.

### 2.3 VNPay

**Public payout documentation status:** No public payout/disbursement API documentation was found on VNPay's developer portal. The available documentation covers payment gateway (accepting payments via `pay` command), transaction query (`querydr`), and refund (`refund`) — but not sending payouts to users. VNPay may offer disbursement through dedicated merchant channels, but this is not verified from public first-party documentation.

---

## 3. File Storage Capabilities (First-Party Docs)

### 3.1 AWS S3

**First-party docs:** `docs.aws.amazon.com/AmazonS3/latest/userguide/`

**Documented capabilities:**

- **Object upload:** Single PUT upload (up to 5 GB), Multipart Upload (up to 50 TB, parts from 5 MB to 5 GB). Multipart uploads can be done independently and in parallel. [docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html]
- **Resumable upload:** The S3 Transfer Manager (Java, Python, AWS CLI with CRT) supports automatic resumable uploads. [docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html]
- **Presigned URLs:** Time-limited URLs for download (GET), upload (PUT), and metadata read (HEAD). Valid for up to **7 days** when using IAM user credentials with AWS Signature Version 4 (SigV4). Supports checksum verification: CRC-64/NVME, CRC32, CRC32C, SHA-1, SHA-256, MD5, XXHash64, XXHash3, XXHash128, SHA-512 (with SigV4). [docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html]
- **Credentials for presigned URLs:** IAM user (up to 7 days), IAM role session (up to session expiration), AWS STS temporary credentials (up to credential expiration). [docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html]
- **Console presigned URLs:** Up to 12 hours. [docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html]
- **Private/public visibility:** By default, buckets and objects are private. Access is controlled via IAM policies, bucket policies, ACLs, Block Public Access settings, and presigned URLs. Network path restrictions can be applied via IAM `aws:SourceIp` conditions and VPC endpoint policies. [docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html]
- **Object metadata:** Custom metadata as name-value pairs. Standard metadata: Content-Type, Content-Disposition, Content-Encoding. [docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html]
- **Delete:** Single object delete and S3 Batch Operations for bulk delete. [docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html]
- **Server-side encryption:** SSE-S3 (AES-256) for all new objects by default. SSE-KMS and SSE-C also supported. [docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html]
- **Event notifications:** SNS, SQS, or Lambda on object create/delete events. [docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html]
- **Location info:** As of 2026-07-19, the AWS Global Infrastructure page lists 39 geographic regions with 123 Availability Zones. **No Vietnam region or Local Zone is listed** on the official regions page (page updated 2026-07-15). The earlier mention of a "Hanoi Local Zone (announced June 2026)" could not be confirmed from the current region list. Check `aws.amazon.com/about-aws/global-infrastructure/regions_az/` at deployment time for the latest status. The closest confirmed regions are in Asia Pacific (Singapore, Tokyo, Seoul, Mumbai, Sydney, etc.). [aws.amazon.com/about-aws/global-infrastructure/regions_az/]

### 3.2 Google Cloud Storage

**First-party docs:** `cloud.google.com/storage/docs/overview`

**Documented capabilities:**

- **Object upload:** Via JSON API, XML API, gcloud CLI, or client libraries. [cloud.google.com/storage/docs/overview]
- **Resumable upload:** Dedicated resumable upload mechanism. After initiating a resumable upload, a session URI is returned that the client uses to upload data without additional signed URLs. The session URI can be used by anyone in possession of it. [cloud.google.com/storage/docs/access-control/signed-urls]
- **Signed URLs:** Time-limited URLs for reading or writing objects. Created using the V4 signing process with service account credentials. Maximum expiration is 604,800 seconds (7 days). Supports read, write, and delete operations. The signed URL only works via XML API endpoints. [cloud.google.com/storage/docs/access-control/signed-urls]
- **Credentials for signed URLs:** Typically a service account with an HMAC key, or a user account with an associated HMAC key. [cloud.google.com/storage/docs/access-control/signed-urls]
- **Private/public visibility:** IAM-based and ACL-based control. Signed URLs allow fine-grained access without making objects public.
- **Object metadata:** System-defined metadata (name, generation, content-type) and custom metadata (user-defined key-value pairs). [cloud.google.com/storage/docs/objects]
- **Delete:** Object deletion via API. Also supports Object Lifecycle Management for automatic deletion. [cloud.google.com/storage/docs/overview]
- **Location info:** Closest region to Vietnam is `asia-southeast1` (Singapore). No Vietnam-specific region or local zone is documented.

### 3.3 Cloudinary

**First-party docs:** `cloudinary.com/documentation/upload_images`

**Documented capabilities:**

- **Object upload:** REST Upload API at `https://api.cloudinary.com/v1_1/<cloud_name>/<resource_type>/upload`. Supports image, video, raw, and auto-detect resource types. Upload methods: server-side authenticated (signed), direct browser upload, fetch/remote URL, auto-upload from S3/GCS, multipart/resumable via SDK. SDKs for Node.js, Python, PHP, Java, Go, Ruby, .NET, iOS, Android. [cloudinary.com/documentation/upload_images]
- **Authenticated vs. unsigned upload:** Authenticated (signed) uploads require `api_key` and `signature`. Unsigned uploads use an upload preset (with fewer available parameters as a security precaution). [cloudinary.com/documentation/upload_images]
- **Private access:** Resources can be managed via Admin API for listing, updating, and deleting. Assets can be delivered via signed URLs.
- **Metadata:** Custom metadata fields can be defined and attached to assets. Cloudinary automatically analyzes uploaded assets (format, size, resolution, colors) and indexes this data for search. [cloudinary.com/documentation/upload_images]
- **Delete:** Assets can be deleted via Admin API (`/admin/assets/<public_id>/destroy`) or via the MCP server (Asset Management MCP). [cloudinary.com/documentation/upload_images]
- **Location info:** Cloudinary operates on a global CDN with multi-region storage. No specific Vietnam data residency information is documented. Data residency for Vietnam is not guaranteed.

---

## 4. Provider-Neutral Flows and Boundaries (Architecture Recommendations)

*These are architecture recommendations derived from patterns common across providers, not provider-specific facts from first-party docs.*

### 4.1 Top-Up Flow

```
User → [Payment Provider] → IPN/Webhook → Identity mapping → Wallet crediting
```

1. User selects a deposit amount → system creates `TopUpRequest` with `status=pending` and a domain-generated internal transaction ID (UUID).
2. System calls `PaymentProvider.CreatePaymentURL()` → user is redirected to the provider's payment page.
3. User completes payment on the provider's interface.
4. Provider sends IPN/webhook callback to the system's `POST /webhooks/payment` endpoint (publicly reachable URL).
5. System verifies the webhook signature using the provider's signature algorithm (typically HMAC-SHA256 or HMAC-SHA512 with a per-provider secret).
6. System checks idempotency — if this `provider_reference` has already been processed, return a 200/OK without crediting again.
7. System updates `TopUpRequest` → `completed`, appends an immutable `XuTransaction` (type=`topup`, credit), and updates the user's `XuWallet` balance.
8. If the webhook has not arrived within a timeout window, the system polls `PaymentProvider.QueryTransaction()` using the internal transaction ID or provider reference.

### 4.2 Withdrawal Flow

```
User → WithdrawalRequest → Admin approval → [Payout Provider] → Status → Wallet debiting
```

1. User creates `WithdrawalRequest` (≥ 100,000 Xu, holding period elapsed).
2. System validates balance, holds the Xu (`xu_balance` minus held amount), sets `status=pending`.
3. Admin reviews and approves → system updates `status=approved`.
4. System calls `PayoutProvider.Payout()` (or `BatchPayout()` for batch mode) with the internal transaction ID as idempotency key.
5. Provider processes the transfer (wallet-to-wallet, bank transfer, or bank card transfer, depending on provider capabilities).
6. Provider returns a synchronous response with `resultCode` or sends an async IPN callback.
7. System updates `WithdrawalRequest` status to `completed` or `failed`, debits the wallet with an immutable `XuTransaction` (type=`withdrawal`).

### 4.3 File Access Flow (Private Assets)

1. When a user purchases a product, the system records their entitlement in the `download_permissions` table.
2. When the user requests a download, the system verifies the permission exists and is not expired.
3. The system generates a presigned/signed URL for the file via `FileStorageProvider.GetDownloadURL()` with a short TTL (e.g., 15–60 minutes).
4. The user downloads directly from the provider's endpoint — the system never proxies the file bytes.
5. For public preview images, the system returns a public URL or a long-lived signed URL.

### 4.4 Webhook Handling

- **Exposed endpoint:** One or more `POST /webhooks/{provider}` endpoints reachable from the internet.
- **Signature verification:** Every provider signs webhook payloads with a secret shared at merchant registration. The system must implement per-provider signature verification before accepting any webhook as valid.
- **Idempotency:** Webhook handlers must be idempotent — the same `provider_reference` may be delivered multiple times. Always check whether the corresponding request has already transitioned to a terminal state before processing.
- **Response:** Return HTTP 200 (or provider-specific success code) within a timeout. Providers may retry on non-200 (ZaloPay: up to 3 retries).
- **Fallback polling:** If no webhook arrives within a reasonable window (ZaloPay recommends 15 minutes), proactively call `QueryTransaction()` to resolve the status.

### 4.5 Reconciliation Data (Append-Only Ledger)

Each financial transaction must be recorded in an append-only table (`xu_transactions`) that serves as the system of record for reconciliation:

| Field | Purpose | Source |
|---|---|---|
| `internal_transaction_id` | Domain-generated primary key, UUID | System |
| `idempotency_key` | Ensures safe retry of provider calls | System (same as internal_tx_id) |
| `provider_reference` | Provider's transaction ID for reconciliation | Provider |
| `provider_name` | Which provider handled this leg | System |
| `amount_vnd` | Gross amount in VND | Provider callback or system |
| `amount_xu` | Converted Xu amount (1:1 fixed rate) | System |
| `provider_fee_vnd` | Fee deducted by provider | Provider callback; nullable |
| `gross_amount` | Amount before fees | System |
| `net_amount` | Amount after provider fees | Provider callback or derived |
| `currency` | Always `VND` for MVP | System |
| `status` | `pending` → `completed` / `failed` | State machine |
| `requested_at` | When the user initiated the action | System |
| `completed_at` | When confirmed by provider | Provider callback |
| `webhook_raw` | Raw provider callback JSON | Provider |
| `webhook_signature` | Signature for re-verification | Provider |
| `settlement_batch_id` | Payout batch grouping ID | System (for withdrawals) |

### 4.6 Status Transitions

```
TopUpRequest:       pending → processing → completed | failed
WithdrawalRequest:  pending → approved → processing → completed | failed | rejected
Purchase:           pending → completed (Xu transfer)
Settlement (7 day): pending → completed (Xu released to seller)
```

### 4.7 Provider-Agnostic Domain Interfaces

The payment/payout/file-storage responsibilities must be behind seam interfaces in the domain layer. Provider implementations are adapters that fulfill these interfaces.

```
internal/domain/payment/      — PaymentProvider interface
internal/domain/payout/       — PayoutProvider interface
internal/domain/filestorage/  — FileStorageProvider interface
internal/adapter/payment/*    — per-provider implementations
internal/adapter/payout/*     — per-provider implementations
internal/adapter/filestorage/* — per-provider implementations
```

---

## 5. Legal Scope Exclusion

Legal and compliance research (Vietnamese regulations on e-money, payment intermediaries, AML/KYC, tax, consumer protection, IP, e-commerce platform obligations, etc.) is **excluded from this document** per the Issue #2 scope narrowing directive. All legal questions remain unresolved and must be researched and concluded separately before the financial system uses real VND.

---

## 6. Provider Summary Table

### 6.1 Payment / Top-Up

| Provider | Create Payment | IPN/Callback | Signature | Query Status | Refund | Sandbox | First-Party Source |
|---|---|---|---|---|---|---|---|
| **VNPay** | Redirect-based payment URL (`vnp_Command=pay`) | IPN server-to-server + Return URL redirect | HMAC-SHA512 (SecureHash with secret key) | QueryDr API (`vnp_Command=querydr`) | Full/partial refund (`vnp_Command=refund`) | Yes (self-registration at `sandbox.vnpayment.vn/devreg/`) | sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/pay.html; sandbox.vnpayment.vn/apis/docs/truy-van-hoan-tien/querydr&refund.html |
| **MoMo** | AIO Gateway → `payUrl` | IPN callback (PayPal IPN-style) | HMAC-SHA256 + RSA + AES | Transaction status query | Refund (via SDK) | Yes (dev/prod env) | github.com/momo-wallet; developer.momo.vn |
| **ZaloPay** | `createOrder` → `order_url` / `zp_trans_token` | Webhook (POST callback, retry ×3, duplicate warnings) | HMAC-SHA256 (`key2`) | QueryOrder (poll after 15 min) | Unknown | Yes (test apps + test wallets + callback log portal) | docs.zalopay.vn/docs/developer-tools/knowledge-base/callback/; github.com/zalopay-samples |

### 6.2 Payout

| Provider | Payout API | Batch | Wallet / Bank | Pre-Verify | Balance Check | IPN | First-Party Source |
|---|---|---|---|---|---|---|---|
| **MoMo** | Single Disbursement `POST /disbursement/pay` (`disburseToWallet` / `disburseToBank`) | Previously cited at ≤1,000 items (not re-verified from current docs) | Both: wallet (1k-200M VND) and bank (10k-20M VND) | Wallet status check (`/disbursement/verify`) | Balance (`/disbursement/balance`) | Yes (`ipnUrl` field + callback params documented) | developer.momo.vn/v3/docs/payment/api/disbursement-v2/ |
| **ZaloPay** | Disbursement API (confirmed exists via quickstart repo) | Unknown | Unknown | Unknown | Unknown | Unknown | github.com/zalopay-samples/quickstart-disbursement |
| **VNPay** | No public payout documentation found | N/A | N/A | N/A | N/A | N/A | — |

### 6.3 File Storage

| Provider | Upload | Signed URL | Private/Public | Metadata | Delete | Location Closest to Vietnam |
|---|---|---|---|---|---|---|
| **AWS S3** | Single PUT (≤5 GB) + Multipart (≤50 TB) + Transfer Manager | Presigned URL (up to 7 days with IAM user, SigV4; up to 12h from console) | IAM + Bucket Policy + ACL + Block Public Access + network path restriction | System + custom | Single + Batch | Asia Pacific regions (Singapore, Tokyo, Seoul, etc.). **No Vietnam Local Zone confirmed** on aws.amazon.com/about-aws/global-infrastructure/regions_az/ as of 2026-07-19 |
| **GCS** | JSON/XML API + gcloud CLI + SDK + Resumable upload | Signed URL (up to 7 days, V4 signing, service account) | IAM + ACL | System + custom | Single + Object Lifecycle Management | `asia-southeast1` (Singapore) |
| **Cloudinary** | REST Upload API (`/v1_1/<cloud>/<type>/upload`) + SDKs (Node, Python, PHP, Java, Go, .NET, etc.) + unsigned (preset) | Signed delivery URLs; Admin API | Upload presets + Admin API + signed delivery | Custom metadata fields + auto-generated analysis data | Admin API | Global CDN (no Vietnam data residency guarantee) |

---

## 7. Decisions, Risks, and Architecture Notes

### 7.1 Decisions Carried Forward (from Issue #2 scope)

| # | Decision | Source | Status |
|---|---|---|---|
| D1 | Xu 1:1 VND, fixed rate | CONTEXT.md:63 | Confirmed |
| D2 | Platform fee 25%, seller 75% | CONTEXT.md:66-69 | Confirmed |
| D3 | Settlement after 7 days | CONTEXT.md:71-76 | Confirmed |
| D4 | Top-up hold period 1 day before withdrawal | CONTEXT.md:34-35 | Confirmed (outside this research scope) |
| D5 | Minimum withdrawal 100,000 Xu | CONTEXT.md:78-80 | Confirmed |
| D6 | No refund after purchase | CONTEXT.md:82-84 | Confirmed (outside this research scope) |
| D7 | Provider behind interface seam | CONTEXT.md — implied by domain boundaries | Confirmed |
| D8 | Manual withdrawal via admin approval (MVP) | CONTEXT.md — implied by withdrawal flow | Confirmed |
| D9 | MVP currency: VND only | CONTEXT.md:62-63 | Confirmed |

### 7.2 Risks and Technical Known Unknowns

| # | Risk | Severity | Mitigation / Note |
|---|---|---|---|
| R1 | **Provider fees unknown until merchant registration** | MEDIUM | All three payment providers and two payout providers do not disclose fees publicly. Budget model must be built after contract negotiation. |
| R2 | **Payout capabilities vary significantly** | MEDIUM | MoMo has the most complete public payout documentation (wallet + bank, wallet verify, balance check, IPN). ZaloPay's payout docs URL now returns 404 — capability is confirmed to exist via repo but API details are unknown. VNPay has no documented payout API. Do not rely on one provider for both payment and payout. |
| R3 | **File storage provider Vietnam latency / data residency** | MEDIUM | AWS does not list a Vietnam Local Zone on the official regions page as of 2026-07-19. GCS Singapore adds latency. Cloudinary offers best CDN performance but no documented Vietnam data residency. |
| R4 | **Webhook delivery is not guaranteed** | LOW | ZaloPay explicitly warns of duplicate callbacks and retries up to 3 times. VNPay relies on IPN + Return URL dual channel. The system must implement both webhook handling AND status polling. |
| R5 | **Idempotency is the merchant's responsibility** | LOW | All providers expect the merchant to generate unique transaction IDs and handle duplicate webhooks (ZaloPay: explicit warning; MoMo: `requestId` idempotency with `Duplicated requestId` result code). The system must enforce idempotency at the application layer. |
| R6 | **MoMo payout amount limits per transaction** | LOW | Wallet payout max 200,000,000 VND. Bank payout max 20,000,000 VND. For withdrawals above these thresholds, the system must split into multiple payout requests or choose wallet payouts for larger amounts. |
| R7 | **ZaloPay Disbursement documentation degraded** | MEDIUM | The English disbursement overview page (`/en/v2/payments/disbursement/overview.html`) now returns 404. This documentation was previously accessible but appears to have been removed or relocated. If ZaloPay is a candidate for payout, the merchant must contact ZaloPay for current API documentation. |
| R8 | **VNPay `vnp_Amount` ×100 multiplier** | LOW | VNPay requires multiplying the amount by 100 (to eliminate decimals). This is a per-provider implementation detail that the adapter layer must handle — the domain layer should deal in whole VND amounts, with conversion only at the adapter boundary. |
| R9 | **MoMo disbursement API 30-second minimum timeout** | LOW | The disbursement pay API requires a minimum 30-second HTTP timeout. Withdrawal flow must account for this in request timeout configuration. |

### 7.3 Open Questions Requiring Merchant Registration

These questions cannot be answered from public first-party documentation and require merchant account registration and contract negotiation with each provider:

1. **Fees per transaction** for both payment (top-up) and payout (disbursement) — none of the three payment providers or two payout candidates publish fee schedules publicly.
2. **Transaction limits** (daily, monthly) — only per-transaction limits are documented for MoMo payout.
3. **Settlement timelines** — how quickly funds from payment top-ups settle to the merchant's provider account.
4. **Minimum initial deposit** or reserve requirement for the payout fund (MoMo disbursement balance).
5. **Production SLA** and uptime commitments.
6. **Onboarding time** and documentation required for merchant account approval.
7. **ZaloPay disbursement API details** — the official docs page is no longer accessible at the previously known URL.
8. **VNPay payout API** — confirm with VNPay merchant support whether payout/disbursement is available through dedicated channels.

### 7.4 Architecture Design Notes

| # | Note |
|---|---|
| N1 | **Two separate webhook endpoints by provider.** Each payment provider requires its own `POST /webhooks/{provider}` endpoint because signature verification algorithms and callback data schemas differ. |
| N2 | **Fallback polling pattern.** Every provider recommends or requires status polling as a fallback when callbacks are missed. Implement a background job that queries unresolved transactions after a configurable timeout (e.g., 15 minutes for ZaloPay). |
| N3 | **Provider adapter isolation for amount scaling.** VNPay requires `amount × 100` in the payment URL (VND without decimals), while other providers may expect raw VND amounts. This translation belongs in the adapter layer only. |
| N4 | **Withdrawal amount splitting.** If using MoMo bank payout (max 20M VND/transaction), the domain layer must support splitting a single withdrawal into multiple provider calls when the amount exceeds the per-transaction limit. |
| N5 | **Append-only ledger for reconciliation.** Every financial event (top-up, purchase, settlement, withdrawal) must produce an immutable XuTransaction record, keyed by the domain-generated internal transaction ID. The ledger serves as the single source of truth for reconciling against provider reports. |

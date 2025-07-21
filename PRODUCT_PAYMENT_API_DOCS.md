# Product Payment API Documentation

This document provides comprehensive information about the Product Payment APIs for the referral system.

## Base URL
```
/product-payment
```

## APIs Overview

1. **Get Referral Code Statistics** - Get payment statistics for a specific referral code
2. **Get All Referral Code Statistics (Paginated)** - Get payment statistics for all referral codes
3. **Get Successful Payments (Paginated)** - Get all successful payments with pagination
4. **Get Total Payment Statistics** - Get overall payment statistics

---

## 1. Get Referral Code Statistics

### Endpoint
```
GET /product-payment/referral-stats/{referral_code}
```

### Description
Retrieves payment statistics for a specific referral code including success, pending, failed, and total counts.

### Request Parameters

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `referral_code` | string | Yes | The referral code to get statistics for |

### Response Format
```json
{
  "referral_code": "string",
  "success_count": 0,
  "pending_count": 0,
  "failed_count": 0,
  "total_count": 0
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `referral_code` | string | The referral code |
| `success_count` | number | Number of successful payments |
| `pending_count` | number | Number of pending payments |
| `failed_count` | number | Number of failed payments |
| `total_count` | number | Total number of payments |

### Example Request
```bash
GET /product-payment/referral-stats/REF123
```

### Example Response
```json
{
  "referral_code": "REF123",
  "success_count": 25,
  "pending_count": 3,
  "failed_count": 2,
  "total_count": 30
}
```

### Error Responses
- **400 Bad Request**: When referral code is missing
- **500 Internal Server Error**: When database error occurs

---

## 2. Get All Referral Code Statistics

### Endpoint
```
GET /product-payment/referral-stats
```

### Description
Retrieves payment statistics for all referral codes with pagination, ordered by success count and total count in descending order.

### Request Parameters

#### Query Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | number | No | 1 | Page number (starts from 1) |
| `page_size` | number | No | 10 | Number of items per page |

### Response Format
```json
{
  "data": [
    {
      "referral_code": "string",
      "success_count": 0,
      "pending_count": 0,
      "failed_count": 0,
      "total_count": 0
    }
  ],
  "total_count": 0,
  "page": 0,
  "page_size": 0,
  "total_pages": 0
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `data` | array | Array of referral code statistics |
| `data[].referral_code` | string | The referral code |
| `data[].success_count` | number | Number of successful payments |
| `data[].pending_count` | number | Number of pending payments |
| `data[].failed_count` | number | Number of failed payments |
| `data[].total_count` | number | Total number of payments for this referral |
| `total_count` | number | Total number of referral codes |
| `page` | number | Current page number |
| `page_size` | number | Number of items per page |
| `total_pages` | number | Total number of pages |

### Example Request
```bash
GET /product-payment/referral-stats?page=1&page_size=5
```

### Example Response
```json
{
  "data": [
    {
      "referral_code": "REF123",
      "success_count": 25,
      "pending_count": 3,
      "failed_count": 2,
      "total_count": 30
    },
    {
      "referral_code": "REF456",
      "success_count": 18,
      "pending_count": 1,
      "failed_count": 1,
      "total_count": 20
    }
  ],
  "total_count": 200,
  "page": 1,
  "page_size": 5,
  "total_pages": 40
}
```

### Error Responses
- **500 Internal Server Error**: When database error occurs

---

## 3. Get Successful Payments (Paginated)

### Endpoint
```
GET /product-payment/successful-payments
```

### Description
Retrieves all successful payments with pagination support. Returns payment details including amount, invoice, user info, and referral code.

### Request Parameters

#### Query Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | number | No | 1 | Page number (starts from 1) |
| `page_size` | number | No | 10 | Number of items per page |

### Response Format
```json
{
  "data": [
    {
      "amount_student_paid": "string",
      "invoice": "string",
      "user_id": 0,
      "referral_code": "string",
      "transaction_id": "string",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total_count": 0,
  "page": 0,
  "page_size": 0,
  "total_pages": 0
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `data` | array | Array of successful payment objects |
| `data[].amount_student_paid` | string\|null | Amount paid by student |
| `data[].invoice` | string\|null | Invoice number |
| `data[].user_id` | number | User ID who made the payment |
| `data[].referral_code` | string\|null | Referral code used |
| `data[].transaction_id` | string\|null | Transaction ID |
| `data[].created_at` | string\|null | Payment creation timestamp |
| `total_count` | number | Total number of successful payments |
| `page` | number | Current page number |
| `page_size` | number | Number of items per page |
| `total_pages` | number | Total number of pages |

### Example Request
```bash
GET /product-payment/successful-payments?page=1&page_size=5
```

### Example Response
```json
{
  "data": [
    {
      "amount_student_paid": "500.00",
      "invoice": "INV-2024-001",
      "user_id": 12345,
      "referral_code": "REF123",
      "transaction_id": "TXN_ABC123",
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "amount_student_paid": "750.00",
      "invoice": "INV-2024-002",
      "user_id": 12346,
      "referral_code": "REF456",
      "transaction_id": "TXN_DEF456",
      "created_at": "2024-01-14T15:45:00Z"
    }
  ],
  "total_count": 150,
  "page": 1,
  "page_size": 5,
  "total_pages": 30
}
```

### Pagination Usage
- To get the first page: `?page=1&page_size=10`
- To get the second page: `?page=2&page_size=10`
- To get 20 items per page: `?page=1&page_size=20`

### Error Responses
- **500 Internal Server Error**: When database error occurs

---

## 4. Get Total Payment Statistics

### Endpoint
```
GET /product-payment/total-stats
```

### Description
Retrieves overall payment statistics across all payments with referral codes.

### Request Parameters
None

### Response Format
```json
{
  "success_count": 0,
  "pending_count": 0,
  "failed_count": 0,
  "total_count": 0
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `success_count` | number | Total number of successful payments |
| `pending_count` | number | Total number of pending payments |
| `failed_count` | number | Total number of failed payments |
| `total_count` | number | Total number of payments |

### Example Request
```bash
GET /product-payment/total-stats
```

### Example Response
```json
{
  "success_count": 1250,
  "pending_count": 45,
  "failed_count": 23,
  "total_count": 1318
}
```

### Error Responses
- **500 Internal Server Error**: When database error occurs

---

## Common HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success - Request completed successfully |
| 400 | Bad Request - Invalid request parameters |
| 500 | Internal Server Error - Server-side error occurred |

## Common Headers

### Request Headers
```
Content-Type: application/json
```

### Response Headers
```
Content-Type: application/json
```

## Notes for Frontend Development

1. **Pagination**: The successful payments API supports pagination. Always handle the pagination metadata to implement proper UI controls.

2. **Null Values**: Some fields can be null (marked as `string|null`). Handle these appropriately in your UI.

3. **Error Handling**: Always implement proper error handling for 400 and 500 status codes.

4. **Date Format**: Timestamps are returned in ISO 8601 format (`YYYY-MM-DDTHH:mm:ssZ`).

5. **Referral Codes**: All successful payment responses include a non-null `referral_code`.

6. **Currency**: Amount fields are returned as strings to preserve decimal precision.

## Example Frontend Usage (JavaScript)

```javascript
// Get referral code stats
async function getReferralStats(referralCode) {
  try {
    const response = await fetch(`/product-payment/referral-stats/${referralCode}`);
    if (!response.ok) throw new Error('Failed to fetch referral stats');
    return await response.json();
  } catch (error) {
    console.error('Error:', error);
  }
}

// Get paginated successful payments
async function getSuccessfulPayments(page = 1, pageSize = 10) {
  try {
    const response = await fetch(`/product-payment/successful-payments?page=${page}&page_size=${pageSize}`);
    if (!response.ok) throw new Error('Failed to fetch payments');
    return await response.json();
  } catch (error) {
    console.error('Error:', error);
  }
}

// Get all referral stats with pagination
async function getAllReferralStats(page = 1, pageSize = 10) {
  try {
    const response = await fetch(`/product-payment/referral-stats?page=${page}&page_size=${pageSize}`);
    if (!response.ok) throw new Error('Failed to fetch all referral stats');
    return await response.json();
  } catch (error) {
    console.error('Error:', error);
  }
}

// Get total payment stats
async function getTotalStats() {
  try {
    const response = await fetch('/product-payment/total-stats');
    if (!response.ok) throw new Error('Failed to fetch total stats');
    return await response.json();
  } catch (error) {
    console.error('Error:', error);
  }
}
``` 
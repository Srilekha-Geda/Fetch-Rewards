# Fetch-Rewards Receipt Processor API Service
<!-- Receipt Processor Service
This README provides an overview of the Receipt Processor API, its setup, usage, and endpoints. -->


This README provides an overview of the Receipt Processor API, its setup, usage, and endpoints.

---

## **Overview**

The Receipt Processor API allows users to submit receipt details and calculate reward points based on specific rules. The API generates a unique identifier for each receipt and allows users to retrieve the points for a specific receipt using the unique ID.

---

## **Features**
1. Submit receipt details via a `POST` request to process and calculate reward points.
2. Retrieve calculated points for a receipt using its unique ID.
3. Unique ID generation using SHA-256 based on receipt attributes.

---

## **Endpoints**

### **1. Root Endpoint**
- **URL**: `/`
- **Method**: `GET`
- **Description**: Basic health check for the server.
- **Response**:
  - **200 OK**: `Server is running!`

### **2. Process Receipt**
- **URL**: `/receipts/process`
- **Method**: `POST`
- **Description**: Submits a receipt for processing and calculates reward points.
- **Request Body** (JSON):
  ```json
  {
    "retailer": "Target",
    "purchaseDate": "2023-09-15",
    "purchaseTime": "14:33",
    "items": [
      {
        "shortDescription": "Pepsi - 12 pack",
        "price": "5.99"
      },
      {
        "shortDescription": "Bread",
        "price": "2.49"
      }
    ],
    "total": "8.48"
  }
- **Response** (JSON):
  ```json
  {
  "id": "f87c2f66b780bb86e4dfbeefc3d2b4dc"
  }

### **3. Get Receipt Points**
- **URL**: `/receipts/{id}`
- **Method**: `GET`
- **Description**: Retrieves the points for a receipt using its unique ID.
- **Path Parameter**: 
  - **{id}**: The unique identifier of the receipt.
- **Response** (JSON):
  - **Success**: 
  ```json
  {
  "points": "64"
  }
  - **Error**: 
  ```json
  {
  "error": "Receipt not found"
  }

---

## **Reward Calculation Rules**
The API calculates points for a receipt based on the following rules:
1. **Retailer Name**: 1 point for each alphanumeric character.
2. **Total Amount**: 
- 50 points if the total is a whole number (e.g., `10.00`).
- 25 points if the total is a multiple of `0.25`.
3. **Items**: 
- 5 points for every 2 items.
- If the item's description length is a multiple of 3, add `ceil(price * 0.2)` points.
4. **Purchase Date**: 6 points if the day of the date is odd.
5. **Purchase Time**: 10 points if the time is between `2:00 PM` and `3:59 PM`.

---

## **Reward Calculation Rules**
The API calculates points for a receipt based on the following rules:
1. **Retailer Name**: 1 point for each alphanumeric character.
2. **Total Amount**: 
- 50 points if the total is a whole number (e.g., `10.00`).
- 25 points if the total is a multiple of `0.25`.
3. **Items**: 
- 5 points for every 2 items.
- If the item's description length is a multiple of 3, add `ceil(price * 0.2)` points.
4. **Purchase Date**: 6 points if the day of the date is odd.
5. **Purchase Time**: 10 points if the time is between `2:00 PM` and `3:59 PM`.


---

## **Setup and Installation**
1. **Install Go**: Ensure you have Go installed. You can download it from golang.org.
2. **Clone the Repository**: 
    ```bash
    git clone <repository-url>
    cd <repository-directory>
 
3. **Run the Server**:  The server will start on `http://localhost:8080`.
    ```bash
    go run main.go



---

## **Testing**
### **Using Postman**
1. Import the following endpoints into Postman:
- `POST http://localhost:8080/receipts/process`
- `GET http://localhost:8080/receipts/{id}`
2. Send the requests with appropriate body and parameters.

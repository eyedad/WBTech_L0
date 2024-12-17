import requests
import json
import os

BASE_URL = "http://localhost:8080"  # Убедитесь, что ваш сервис запущен на этом адресе

def test_get_all_orders():
    response = requests.get(f"{BASE_URL}/orders")
    if response.status_code == 200:
        print("GET /orders passed")
        print("Status Code:", response.status_code)
        print("Response:", response.json())
    else:
        print("GET /orders failed")
        print("Status Code:", response.status_code)


def test_get_order_by_id(order_id):
    response = requests.get(f"{BASE_URL}/orders/{order_id}")
    if response.status_code == 200:
        print(f"GET /orders/{order_id} passed")
        print("Status Code:", response.status_code)
        print("Response:", response.json())
    else:
        print(f"GET /orders/{order_id} failed")
        print("Status Code:", response.status_code)


def test_add_order(order_data):
    response = requests.post(f"{BASE_URL}/orders", json=order_data)
    if response.status_code == 201:
        print("POST /orders passed")
        print("Status Code:", response.status_code)
        print("Response:", response.json())
    else:
        print("POST /orders failed")
        print("Status Code:", response.status_code)


def test_all_orders_from_files(test_files_directory):
    for filename in os.listdir(test_files_directory):
        if filename.endswith('.json'):
            file_path = os.path.join(test_files_directory, filename)
            with open(file_path) as f:
                try:
                    new_order_data = json.load(f)
                    print(f"Testing adding order from {filename}...")
                    test_add_order(new_order_data)
                except json.JSONDecodeError as e:
                    print(f"Failed to decode JSON from {filename}: {e}")


if __name__ == "__main__":
    test_all_orders_from_files('test_files')
    test_get_all_orders()
    test_get_order_by_id("h890feb7b2b84b6test")

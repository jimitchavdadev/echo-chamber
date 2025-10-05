# test_websockets.py

import requests
import websocket
import threading
import time
import json

BASE_URL = "http://localhost:8080/api"
WS_URL = "ws://localhost:8080/api/ws"

# --- ANSI Color Codes for pretty printing ---
GREEN = "\033[92m"
YELLOW = "\033[93m"
RED = "\033[91m"
RESET = "\033[0m"

# Global list to store messages received by user_one
received_messages = []
ws_connection_open = True

def print_step(step, message):
    print(f"\n{YELLOW}{step}: {message}{RESET}")

def print_success(message):
    print(f"  - {GREEN}SUCCESS:{RESET} {message}")

def print_failure(message):
    print(f"  - {RED}FAILURE:{RESET} {message}")

def run_user_one_ws_client(cookie_header):
    """This function runs in a separate thread, listening for messages."""
    global ws_connection_open
    try:
        ws = websocket.create_connection(WS_URL, header=cookie_header)
        print("  - [Thread 1] User One WebSocket connected and listening...")
        while ws_connection_open:
            try:
                # Set a timeout so the loop doesn't block forever
                ws.settimeout(1.0)
                message = ws.recv()
                print(f"  - [Thread 1] User One Received: {message[:100]}...")
                received_messages.append(json.loads(message))
            except websocket.WebSocketTimeoutException:
                continue # No message received, continue listening
            except Exception as e:
                print_failure(f"[Thread 1] Error receiving message: {e}")
                break
        ws.close()
        print("  - [Thread 1] User One WebSocket connection closed.")
    except Exception as e:
        print_failure(f"Could not connect User One WebSocket: {e}")

def main():
    global ws_connection_open
    
    print_step("STEP 1", "Setting up test users and content...")
    
    # Use sessions to automatically handle cookies
    session_one = requests.Session()
    session_two = requests.Session()

    # Register users (ignore errors if they already exist)
    session_one.post(f"{BASE_URL}/register", json={"username": "py_user_one", "email": "py_one@test.com", "password": "password123"})
    session_two.post(f"{BASE_URL}/register", json={"username": "py_user_two", "email": "py_two@test.com", "password": "password123"})

    # Login and get user data
    res_one = session_one.post(f"{BASE_URL}/login", json={"email": "py_one@test.com", "password": "password123"})
    user_one_id = res_one.json()["user"]["id"]
    
    res_two = session_two.post(f"{BASE_URL}/login", json={"email": "py_two@test.com", "password": "password123"})
    
    print(f"  - Logged in User One (ID: {user_one_id}) and User Two.")

    # User One creates a post
    res_post = session_one.post(f"{BASE_URL}/posts", json={"content": "A post for the Python test suite!"})
    post_id = res_post.json()["ID"]
    print(f"  - User One created Post ID: {post_id}")

    print_step("STEP 2", "Starting WebSocket listener for User One...")
    
    # Prepare cookie header for websocket connection
    user_one_cookie_str = f"jwt={session_one.cookies.get('jwt')}"
    user_one_ws_header = {"Cookie": user_one_cookie_str}

    # Start the listening client in a background thread
    listener_thread = threading.Thread(target=run_user_one_ws_client, args=(user_one_ws_header,))
    listener_thread.start()
    time.sleep(2) # Give it a moment to connect

    print_step("STEP 3", "User Two performs actions to trigger real-time events...")
    
    # Action 1: User Two likes User One's post
    print("  - User Two is liking the post...")
    session_two.post(f"{BASE_URL}/posts/{post_id}/like")
    time.sleep(1) # Wait for the notification to be received

    # Action 2: User Two sends a chat message to User One
    print("  - User Two is sending a chat message...")
    try:
        user_two_cookie_str = f"jwt={session_two.cookies.get('jwt')}"
        ws_two = websocket.create_connection(WS_URL, header={"Cookie": user_two_cookie_str})
        chat_payload = {
            "type": "chat_message",
            "payload": {"receiverId": user_one_id, "content": "Hello from Python!"}
        }
        ws_two.send(json.dumps(chat_payload))
        ws_two.close()
        print("  - Chat message sent.")
    except Exception as e:
        print_failure(f"Could not send chat message: {e}")
    
    time.sleep(2) # Wait for the chat message to be received

    print_step("STEP 4", "Stopping listener and verifying results...")
    ws_connection_open = False # Signal the listener thread to stop
    listener_thread.join() # Wait for the thread to finish

    # --- VERIFICATION ---
    if len(received_messages) != 2:
        print_failure(f"Expected 2 messages, but received {len(received_messages)}")
    else:
        print_success(f"Received 2 messages as expected.")

        # Check for notification
        notification = next((msg for msg in received_messages if msg.get("type") == "new_notification"), None)
        if notification:
            print_success("Found 'new_notification' message.")
            if notification["payload"]["EntityID"] != post_id:
                print_failure(f"Notification post ID was {notification['payload']['EntityID']}, expected {post_id}")
        else:
            print_failure("Did not find 'new_notification' message.")

        # Check for chat message
        chat_message = next((msg for msg in received_messages if msg.get("type") == "new_chat_message"), None)
        if chat_message:
            print_success("Found 'new_chat_message' message.")
            if chat_message["payload"]["Content"] != "Hello from Python!":
                print_failure("Chat message content was incorrect.")
        else:
            print_failure("Did not find 'new_chat_message' message.")

if __name__ == "__main__":
    main()
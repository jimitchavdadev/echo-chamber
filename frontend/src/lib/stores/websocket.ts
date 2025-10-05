import { writable } from 'svelte/store';
import type { Notification, ChatMessage } from '$lib/types';

interface WsMessage {
  type: string;
  payload: any;
}

// Store for raw notifications
export const notifications = writable<Notification[]>([]);
export const unreadNotificationCount = writable<number>(0);

// Store for chat messages, organized by user ID
export const chatMessages = writable<Record<number, ChatMessage[]>>({});

let socket: WebSocket | null = null;

export function connectWebSocket() {
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return;
  }

  socket = new WebSocket('ws://localhost:8080/api/ws');

  socket.onopen = () => {
    console.log('WebSocket connection established');
  };

  socket.onmessage = (event) => {
    const message: WsMessage = JSON.parse(event.data);

    if (message.type === 'new_notification') {
      const newNotification: Notification = message.payload;
      notifications.update(items => [newNotification, ...items]);
      unreadNotificationCount.update(n => n + 1);
    }

    if (message.type === 'new_chat_message') {
      const newChatMessage: ChatMessage = message.payload;
      const partnerId = newChatMessage.SenderID; // We are the receiver
      chatMessages.update(chats => {
        if (!chats[partnerId]) {
          chats[partnerId] = [];
        }
        chats[partnerId].push(newChatMessage);
        return chats;
      });
    }
  };

  socket.onclose = () => {
    console.log('WebSocket connection closed');
    socket = null;
    // Optional: attempt to reconnect
  };

  socket.onerror = (error) => {
    console.error('WebSocket error:', error);
  };
}

export function sendChatMessage(receiverId: number, content: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const message = {
      type: 'chat_message',
      payload: { receiverId, content }
    };
    socket.send(JSON.stringify(message));
  }
}

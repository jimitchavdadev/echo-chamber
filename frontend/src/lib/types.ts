export interface User {
  id: number;
  username: string;
  email: string;
  bio: string;
}

export interface Post {
  ID: number;
  CreatedAt: string;
  content: string;
  author: User;
  likeCount: number;
  isLiked: boolean;
}

export interface Comment {
  ID: number;
  CreatedAt: string;
  content: string;
  postId: number;
  author: User;
}

export interface Notification {
  ID: number;
  CreatedAt: string;
  ActorID: number;
  Type: 'like' | 'comment';
  EntityID: number; // e.g., Post ID
  IsRead: boolean;
  Actor: User;
}

export interface ChatMessage {
  ID: number;
  CreatedAt: string;
  SenderID: number;
  ReceiverID: number;
  Content: string;
  Sender: User;
}

export interface ProfileData {
  id: number;
  username:string;
  bio: string;
  isFollowing: boolean;
}

export interface User {
  id: number;
  username: string;
  email: string;
  bio: string;
}

// Updated Post type
export interface Post {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  content: string;
  author: User;
  likeCount: number;
  isLiked: boolean;
}

// New Comment type
export interface Comment {
  ID: number;
  CreatedAt: string;
  content: string;
  postId: number;
  author: User;
}

export interface ProfileData {
  id: number;
  username:string;
  bio: string;
  isFollowing: boolean;
}

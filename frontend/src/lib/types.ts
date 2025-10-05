export interface User {
  id: number;
  username: string;
  email: string;
  bio: string;
}

export interface Post {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  content: string;
  author: User;
}

export interface ProfileData {
  id: number;
  username: string;
  bio: string;
  isFollowing: boolean;
}

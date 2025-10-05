import { writable, type Writable } from 'svelte/store';

interface User {
  id: number;
  username: string;
  email: string;
  bio?: string;
}

interface AuthState {
  user: User | null;
}

export const authStore: Writable<AuthState> = writable({
  user: null
});

import { createSlice } from '@reduxjs/toolkit';

export interface User {
  id: number;
  username: string;
  full_name: string;
  is_active: boolean;
}

interface AuthState {
  accessToken: string | null;
  expiresAt: string | null;
  isAuth: boolean;
  loading: boolean;
  user: User | null;
}

const initialState: AuthState = {
  accessToken: null,
  expiresAt: null,
  isAuth: false,
  loading: true,
  user: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setAuth(state, action) {
      state.accessToken = action.payload.access_token;
      state.expiresAt = action.payload.expires_at;
      state.isAuth = true;
      state.loading = false;
      if (action.payload.user) {
        state.user = action.payload.user;
      }
      // state.user = action.payload.user;
    },
    clearAuth(state) {
      state.accessToken = null;
      state.expiresAt = null;
      state.isAuth = false;
      state.loading = false;
      state.user = null;
    },
    setLoading(state, action) {
      state.loading = action.payload;
    },
  },
});

export const { setAuth, clearAuth, setLoading } = authSlice.actions;
export default authSlice.reducer;

// 📌 Slice = логика + состояние + экшены в одном файле
// import { createSlice } from '@reduxjs/toolkit';
// import type { PayloadAction } from '@reduxjs/toolkit';

// interface User {
//   username: string;
// }

// interface AuthState {
//   user: User | null;
//   isAuth: boolean;
// }

// const initialState: AuthState = {
//   user: null,
//   isAuth: false,
// };

// const authSlice = createSlice({
//   name: 'auth',
//   initialState,
//   reducers: {
//     loginSuccess(state, action: PayloadAction<User>) {
//       state.user = action.payload;
//       state.isAuth = true;
//     },
//     logout(state) {
//       state.user = null;
//       state.isAuth = false;
//     },
//   },
// });

// export const { loginSuccess, logout } = authSlice.actions;
// export default authSlice.reducer;

import { createSlice } from '@reduxjs/toolkit';

interface AuthState {
  accessToken: string | null;
  expiresAt: string | null;
  isAuth: boolean;
}

const initialState: AuthState = {
  accessToken: null,
  expiresAt: null,
  isAuth: false,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setAuth(state, action) {
      state.accessToken = action.payload.access_token;
      state.expiresAt = action.payload.expires_at;
      state.isAuth = true;
    },
    clearAuth(state) {
      state.accessToken = null;
      state.expiresAt = null;
      state.isAuth = false;
    },
  },
});

export const { setAuth, clearAuth } = authSlice.actions;
export default authSlice.reducer;
